package jobs

import (
	"sync"

	"github.com/fellah/tcache/db"
	"github.com/fellah/tcache/log"
	"github.com/fellah/tcache/sletat"
)

const (
	WORKERS_NUM = 16
	BULK_SIZE   = 256
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
					if tour.TicketsIncluded != 1 {
						continue
					}

					if tour.HotelIsInStop != 0 && tour.HotelIsInStop != 2 {
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
		} else {
			return tours, nil
		}
	}

	return nil, err
}

func processTour(packet sletat.PacketInfo, tour *sletat.Tour) {
	tour.CreateDate = packet.CreateDate

	tour.DptCityId = packet.DptCityId
	tour.CountryId = packet.CountryId

	if operator, ok := operators[tour.SourceId]; ok {
		// BYR = RUB * exchange rate
		tour.PriceByr = int(float64(tour.Price) * operator.ExchangeRateRur)
		// EUR = BYR * exchange rate
		tour.PriceEur = int(float64(tour.PriceByr) * operator.ExchangeRateEur)
		// USD = BYR * exchange rate
		tour.PriceUsd = int(float64(tour.PriceByr) * operator.ExchangeRateUsd)
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

func finalize(end <-chan bool) {
	go func() {
		<-end

		//db.RemoveExpiredTours()
		db.VacuumTours()
		db.AgregateToursByCountry()
		db.AgregateToursByHotel()

		log.Info.Println("END")
	}()
}
