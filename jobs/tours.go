package jobs

import (
	"sync"
	"time"

	"github.com/fellah/tcache/data"
	"github.com/fellah/tcache/db"
	"github.com/fellah/tcache/log"
	"github.com/fellah/tcache/sletat"
	"github.com/fellah/tcache/stat"
	"github.com/fellah/tcache/cache"
	"github.com/fellah/tcache/prefilter"
)

func init() {
	mx = sync.Mutex{}
}

var mx sync.Mutex

const (
	workersNum = 16
	bulkSize   = 516
)

func fetchTours(packets <-chan data.PacketInfo, stat *stat.Tours, end chan bool) {
	wg := new(sync.WaitGroup)
	wg.Add(workersNum)

	// Run multiply workers to read concurrently from one channel.
	for i := 0; i < workersNum; i++ {
		go func() {
			for packet := range packets {
				stat.Output()
				var count uint64 = 0
				var skipped uint64 = 0

				log.Info.Println("fetchTours Run ...")
				tours_channels, err := sletat.FetchTours(packet.Id, tours_channels_used)
				tours := tours_channels[0]

				if err != nil {
					log.Error.Println(err)
					continue
				}

				collect := make(chan data.Tour)
				log.Info.Println("fetchTours collect tours Run ...")
				go collectTours(collect, stat)

				// Process tours before send the to the database.
				log.Info.Println("fetchTours tours loop Run ...")
				for tour := range tours {
					count++

					preProcessTour(packet, &tour)

					processTour(packet, &tour)

					// Partners
					if prefilter.ForPartnersTours(&tour) {
						cache.RegisterTourGroup(tour)
					}

					if isSkipped(&tour) {
						skipped++
						continue
					}

					if !isKidsValid(&tour) {
						stat.KidsIssue <- 1
					}

					collect <- tour
				}
				log.Info.Println("fetchTours tours loop FINISH ...")

				close(collect)

				stat.Total <- count
				stat.Skipped <- skipped
			}

			log.Info.Println("fetchTours gorotine FINISH")
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		log.Info.Println("fetchTours FINISH ...")

		cache.SaveTourGroupsToDB()

		end <- true
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

func finalize(end chan bool, stat *stat.Tours) {
	go func() {
		// wait end signal from all channels
		<-end
		close(end)

		stat.Output()

		log.Info.Println("END")
	}()
}
