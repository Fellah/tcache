package jobs

import (
	"sync"
	"time"

	"github.com/fellah/tcache/data"
	"github.com/fellah/tcache/db"
	"github.com/fellah/tcache/log"
	"github.com/fellah/tcache/sletat"
	"github.com/fellah/tcache/stat"
)

func init() {
	mx = sync.Mutex{}
}

var mx sync.Mutex

const (
	workersNum = 16
	bulkSize   = 516
)

func fetchTours(packets <-chan data.PacketInfo, stat *stat.Tours, end <-chan bool) {
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

				collect := make(chan data.Tour)
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
}

func preProcessTour(packet data.PacketInfo, tour *data.Tour) {
	tour.DptCityId = packet.DptCityId
}

func processTour(packet data.PacketInfo, tour *data.Tour) {
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

func collectTours(tours <-chan data.Tour, stat *stat.Tours) {
	toursBulk := make([]data.Tour, 0, bulkSize)
	for tour := range tours {
		toursBulk = append(toursBulk, tour)

		if len(toursBulk) == bulkSize {
			mx.Lock()
			start := time.Now()
			println("SaveTours START")
			db.SaveTours(toursBulk)
			println("SaveTours END ", time.Since(start).String())
			mx.Unlock()

			toursBulk = make([]data.Tour, 0, bulkSize)
		}
	}
	mx.Lock()
	start := time.Now()
	println("SaveTours START")
	db.SaveTours(toursBulk)
	println("SaveTours END ", time.Since(start).String())
	mx.Unlock()
}

func isSkipped(tour *data.Tour) bool {
	if !isCityActive(tour.TownId) {
		return true
	}

	return false
}

func finalize(ends []chan bool, stat *stat.Tours, channels []chan data.PacketInfo) {
	go func() {
		// wait end signal from all end channels
		for _,end := range ends {
			<-end
			close(end)
		}

		for _,channel := range channels {
			close(channel)
		}

		stat.Output()

		log.Info.Println("END")
	}()
}
