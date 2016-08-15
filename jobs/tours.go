package jobs

import (
	"sync"

	"github.com/fellah/tcache/db"
	"github.com/fellah/tcache/log"
	"github.com/fellah/tcache/sletat"
	"github.com/fellah/tcache/stat"
)

const (
	workersNum = 16
	bulkSize   = 516
)

func fetchTours(packets <-chan sletat.PacketInfo, stat *stat.Tours) chan bool {
	end := make(chan bool)

	wg := new(sync.WaitGroup)
	wg.Add(workersNum)

	// Run multiply workers to read concurrently from one channel.
	for i := 0; i < workersNum; i++ {
		go func() {
			for packet := range packets {
				stat.Output()
				var count uint64 = 0
				var skipped uint64 = 0

				tours, err := sletat.FetchTours(packet.Id)
				if err != nil {
					log.Error.Println(err)
					continue
				}

				collect := make(chan sletat.Tour)
				go collectTours(collect, stat)

				// Process tours before send the to the database.
				for tour := range tours {
					count++

					preProcessTour(packet, &tour)

					if isSkipped(&tour) {
						skipped++
						continue
					}

					processTour(packet, &tour)

					if !isKidsValid(&tour) {
						stat.KidsIssue <- 1
					}

					collect <- tour
				}

				close(collect)

				stat.Total <- count
				stat.Skipped <- skipped
			}

			wg.Done()
		}()
	}

	go func() {
		wg.Wait()

		end <- true
		close(end)
	}()

	return end
}

func preProcessTour(packet sletat.PacketInfo, tour *sletat.Tour) {
	tour.DptCityId = packet.DptCityId
}

func processTour(packet sletat.PacketInfo, tour *sletat.Tour) {
	tour.CreateDate = packet.CreateDate

	tour.CountryId = packet.CountryId

	if operator, ok := operators[tour.SourceId]; ok {
		// BYN = RUB * exchange rate
		tour.PriceByr = int(float64(tour.Price) * operator.ExchangeRateRur)

		// EUR = BYN / exchange rate
		if tour.PriceEur > 0 && operator.ExchangeRateEur > 0 {
			tour.PriceEur = int(float64(tour.PriceByr) / operator.ExchangeRateEur)
		} else {
			tour.PriceEur = 0
		}

		// USD = BYN / exchange rate
		if tour.PriceByr > 0 && operator.ExchangeRateUsd > 0 {
			tour.PriceUsd = int(float64(tour.PriceByr) / operator.ExchangeRateUsd)
		} else {
			tour.PriceUsd = 0
		}
	}

	processKidsValue(tour)
}

func collectTours(tours <-chan sletat.Tour, stat *stat.Tours) {
	go func() {
		toursBulk := make([]sletat.Tour, 0, bulkSize)
		for tour := range tours {
			toursBulk = append(toursBulk, tour)

			if len(toursBulk) == bulkSize {
				db.SaveTours(toursBulk)

				toursBulk = make([]sletat.Tour, 0, bulkSize)
			}
		}
		db.SaveTours(toursBulk)
	}()
}

func isSkipped(tour *sletat.Tour) bool {
	if !isCityActive(tour.TownId) {
		return true
	}

	return false
}

func finalize(end <-chan bool, stat *stat.Tours) {
	go func() {
		<-end
		stat.Output()
		log.Info.Println("END")
	}()
}
