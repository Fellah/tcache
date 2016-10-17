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

var (
	mx sync.Mutex
	once_save_data = sync.Once{}
)

const (
	workersNum = 16
	bulkSize   = 516
)

func fetchTours(packets <-chan data.PacketInfo, stat *stat.Tours, end chan bool) {
	wg := new(sync.WaitGroup)
	wg.Add(workersNum)

	end_save_process := cronSaveTourGroupsToDB();

	// Run multiply workers to read concurrently from one channel.
	for i := 0; i < workersNum; i++ {
		go func() {
			for packet := range packets {
				stat.Output()
				var count uint64 = 0
				var skipped uint64 = 0

				log.Info.Println("fetchTours Run ...")
				tours, err := sletat.FetchTours(packet.Id)

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
					// ------
					if prefilter.ForPartnersTours(&tour) {
						cache.RegisterTourGroup(tour)
					}
					// ------

					if !isOperatorActive(packet.SourceId) {
						continue
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

		cache.SaveTourGroupsToDB(&once_save_data)

		end_save_process <- true
		end <- true
	}()
}

func cronSaveTourGroupsToDB() (chan bool) {
	end_chan := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker_save_data.C:
				log.Info.Println("CRON: Save Partners Group Data By")
				once_save_data.Do(func() { go cache.SaveTourGroupsToDB(&once_save_data) })
			case <-end_chan:
				log.Info.Println("CRON: End function - finish")
				return
			}
		}
	}()

	return end_chan
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
	if tour.TicketsIncluded != 1 ||
		(tour.HasEconomTicketsDpt != 1 && tour.HasEconomTicketsDpt != 2) ||
		(tour.HasEconomTicketsRtn != 1 && tour.HasEconomTicketsRtn != 2) ||
		(tour.HotelIsInStop != 0 && tour.HotelIsInStop != 2) ||
		tour.HotelId == 0 {
		return true
	}

	if !isCityActive(tour.TownId) {
		return true
	}

	if prefilter.ForHotel(tour) {
		return true
	}

	return false
}
