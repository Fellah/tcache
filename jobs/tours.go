package jobs

import (
	"log"
	"time"

	"github.com/fellah/tcache/db"
	"github.com/fellah/tcache/sletat"
)

func FetchTours(chPocket <-chan sletat.PacketInfo) {
	for {
		select {
		case pocket, ok := <-chPocket:
			if !ok {
				return
			}
			tours, err := sletat.FetchTours(pocket.Id)
			if err != nil {
				log.Println(err)
			}

			for i := range tours {
				tours[i].CreateDate = pocket.CreateDate

				tours[i].DptCityId = pocket.DptCityId

				if operator, ok := operators[tours[i].SourceId]; ok {
					tours[i].PriceByr = currencyPrice(tours[i].Price, operator.ExchangeRateRur)
					tours[i].PriceEur = currencyPrice(tours[i].Price, operator.ExchangeRateRur)
					tours[i].PriceEur = currencyPrice(tours[i].Price, operator.ExchangeRateRur)
				}

			}

			go db.SaveTours(tours)

		case <-time.After(900 * time.Second):
			log.Println("TIMEOUT")
			return
		}
	}
}

func currencyPrice(price int, exchange float64) int {
	return price * int(exchange)
}