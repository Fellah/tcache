package jobs

import (
	"log"
	"sync"

	//"github.com/fellah/tcache/db"
	"github.com/fellah/tcache/sletat"
)

func fetchTours(packets <-chan sletat.PacketInfo) chan sletat.Tour {
	var wg sync.WaitGroup

	out := make(chan sletat.Tour)

	// Run multiply workers to read concurrently from one channel.
	for i := 0; i < WORKERS_NUM; i++ {
		wg.Add(1)

		go func(i int) {
			for packet := range packets {
				tours, err := sletat.FetchTours(packet.Id)
				if err != nil {
					log.Println(err)
					continue
				}

				for tour := range tours {
					processTour(packet, &tour)
					out <- tour
				}
			}

			wg.Done()
			log.Println("END", i)
		}(i)
	}

	go func() {
		wg.Wait()
		log.Println("CLOSE out")
		close(out)
	}()

	return out
}

/*func fetchTours(packets <-chan sletat.PacketInfo, tours chan<- sletat.Tour, wg sync.WaitGroup) {
	for {
		packet, ok := <-packets
		if !ok {
			break
		}

		rawTours := make(chan sletat.Tour)

		fmt.Println("AAAA1")

		// Process raw tour and send it to the channel for save.
		go func() {
			for tour := range rawTours {
				processTour(packet, &tour)
				tours <- tour
			}
		}()

		fmt.Println("AAAA2")



		fmt.Println("AAAA")

		// Read tours from XML stream and send them to them channel
		err := sletat.FetchTours(packet.Id, rawTours)
		if err != nil {
			log.Println(err)
		}

		fmt.Println("AAAA3")

		close(rawTours)

		wg.Done()
	}
}*/

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
		for _ = range tours {

		}

		end <- true

		//close(end)
	}()

	/*
	go func() {
		tours := make([]sletat.Tour, 0, 64)

		for {
			tour, ok := <-chTour
			if !ok {
				end <- true
				break
			}

			tours = append(tours, tour)
			if len(tours) == 16 {
				db.SaveTours(tours)
				tours = make([]sletat.Tour, 0, 256)
			}
		}
	}()*/

	return end
}

func finalize(end <-chan bool) {
	go func() {
		<-end

		log.Println("END")
	}()
}
