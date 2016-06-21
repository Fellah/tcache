package jobs

import (
	"sync"

	"github.com/fellah/tcache/log"
	"github.com/fellah/tcache/sletat"
	"github.com/fellah/tcache/db"
)

const (
	WORKERS_NUM = 32
	BULK_SIZE   = 516
)

func fetchTours(packets <-chan sletat.PacketInfo) chan sletat.Tour {
	out := make(chan sletat.Tour)

	wg := new(sync.WaitGroup)
	wg.Add(WORKERS_NUM)

	// Run multiply workers to read concurrently from one channel.
	for i := 0; i < WORKERS_NUM; i++ {
		go func() {
			for packet := range packets {
				tours, err := tryGetToursChannel(packet.Id)
				if err != nil {
					log.Error.Println(err)
					continue
				}

				for tour := range tours {
					preProcessTour(packet, &tour)

					if isSkipped(&tour) {
						continue
					}

					processTour(packet, &tour)
					out <- tour
				}
			}

			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func tryGetToursChannel(packetId string) (chan sletat.Tour, error) {
	var err error

	for i := 0; i < 3; i++ {
		tours, err := sletat.FetchTours(packetId)
		if err != nil {
			log.Info.Println(packetId, i, err)
			continue
		}

		return tours, nil
	}

	return nil, err
}

func preProcessTour(packet sletat.PacketInfo, tour *sletat.Tour) {
	tour.DptCityId = packet.DptCityId
}

func processTour(packet sletat.PacketInfo, tour *sletat.Tour) {
	tour.CreateDate = packet.CreateDate

	tour.CountryId = packet.CountryId

	if operator, ok := operators[tour.SourceId]; ok {
		// BYR = RUB * exchange rate
		tour.PriceByr = int(float64(tour.Price) * operator.ExchangeRateRur)

		// EUR = BYR * exchange rate
		if tour.PriceEur > 0 && operator.ExchangeRateEur > 0 {
			tour.PriceEur = int(float64(tour.PriceByr) / operator.ExchangeRateEur)
		} else {
			tour.PriceEur = 0
		}

		// USD = BYR * exchange rate
		if tour.PriceByr > 0 && operator.ExchangeRateUsd > 0 {
			tour.PriceUsd = int(float64(tour.PriceByr) / operator.ExchangeRateUsd)
		} else {
			tour.PriceUsd = 0
		}
	}
}

func saveTours(tours <-chan sletat.Tour) <-chan bool {
	end := make(chan bool)

	go func() {
		toursBulk := make([]sletat.Tour, 0, BULK_SIZE)
		for tour := range tours {
			toursBulk = append(toursBulk, tour)
			if len(toursBulk) == BULK_SIZE {
				db.SaveTours(toursBulk)
				toursBulk = make([]sletat.Tour, 0, BULK_SIZE)
			}
		}
		db.SaveTours(toursBulk)

		end <- true
		close(end)
	}()

	return end
}

func isSkipped(tour *sletat.Tour) bool {
	if !isCityActive(tour.TownId) {
		return true
	}

	if tour.TicketsIncluded != 1 {
		return true
	}

	if tour.HotelIsInStop != 0 && tour.HotelIsInStop != 2 {
		return true
	}

	return false
}

func finalize(end <-chan bool) {
	go func() {
		<-end

		//db.VacuumTours()

		log.Info.Println("END")
	}()
}
