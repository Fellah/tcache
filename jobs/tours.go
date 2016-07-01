package jobs

import (
	"sync"

	"github.com/fellah/tcache/db"
	"github.com/fellah/tcache/log"
	"github.com/fellah/tcache/sletat"
	"github.com/fellah/tcache/stat"
)

const (
	workersNum = 32
	bulkSize   = 516
)

func fetchTours(packets <-chan sletat.PacketInfo, stat *stat.Tours) chan sletat.Tour {
	out := make(chan sletat.Tour)

	wg := new(sync.WaitGroup)
	wg.Add(workersNum)

	// Run multiply workers to read concurrently from one channel.
	for i := 0; i < workersNum; i++ {
		go func() {
			for packet := range packets {
				stat.Output()
				var count uint64 = 0
				var skipped uint64 = 0

				tours, err := tryGetToursChannel(packet.Id)
				if err != nil {
					log.Error.Println(err)
					continue
				}

				for tour := range tours {
					count++

					preProcessTour(packet, &tour)

					if isSkipped(&tour) {
						skipped++
						continue
					}

					processTour(packet, &tour)
					out <- tour
				}

				stat.Total <- count
				stat.Skipped <- skipped
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

func saveTours(tours <-chan sletat.Tour, stat *stat.Tours) <-chan bool {
	end := make(chan bool)

	go func() {
		toursBulk := make([]sletat.Tour, 0, bulkSize)
		for tour := range tours {
			checkKidsIssue(&tour, stat)

			i := findDuplicate(tour, toursBulk)
			if i >= 0 && i < len(toursBulk) {
				if tour.Price < toursBulk[i].Price {
					toursBulk[i] = tour
				}
			} else {
				toursBulk = append(toursBulk, tour)
			}

			if len(toursBulk) == bulkSize {
				db.SaveTours(toursBulk)

				toursBulk = make([]sletat.Tour, 0, bulkSize)
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

	return false
}

func finalize(end <-chan bool) {
	go func() {
		<-end
		log.Info.Println("END")
	}()
}

func findDuplicate(tour sletat.Tour, toursBulk []sletat.Tour) int {
	for i := range toursBulk {
		if tour.HotelId == toursBulk[i].HotelId &&
			tour.Checkin == toursBulk[i].Checkin &&
			tour.DptCityId == toursBulk[i].DptCityId &&
			tour.Nights == toursBulk[i].Nights &&
			tour.Adults == toursBulk[i].Adults &&
			tour.Kids == toursBulk[i].Kids &&
			tour.SourceId == toursBulk[i].SourceId &&
			tour.MealId == toursBulk[i].MealId &&
			compareKidsValues(tour.Kid1Age, toursBulk[i].Kid1Age) &&
			compareKidsValues(tour.Kid2Age, toursBulk[i].Kid2Age) &&
			compareKidsValues(tour.Kid3Age, toursBulk[i].Kid3Age) {
			return i
		}
	}

	return -1
}

func checkKidsIssue(tour *sletat.Tour, stat *stat.Tours) {
	kids := 0

	if tour.Kid1Age != nil {
		kids++
	}

	if tour.Kid2Age != nil {
		kids++
	}

	if tour.Kid3Age != nil {
		kids++
	}

	if tour.Kids != kids {
		stat.KidsIssue <- 1
	}
}

func compareKidsValues(vA *int, vB *int) bool {
	if vA == nil && vB != nil {
		return false
	} else if vA != nil && vB == nil {
		return false
	} else if vA == nil && vB == nil {
		return true
	} else {
		return *vA == *vB
	}
}
