package jobs

import (
	"log"

	"github.com/fellah/tcache/db"
	"github.com/fellah/tcache/sletat"
)

func fetchTours(chPacket <-chan sletat.PacketInfo, chTour chan<- sletat.Tour) {
	for packet := range chPacket {
		chRawTour := make(chan sletat.Tour)

		// Process raw tour and send it to the channel for save.
		go func(packet sletat.PacketInfo, chRawTour <-chan sletat.Tour) {
			for tour := range chRawTour {
				processTour(packet, &tour)
				chTour <- tour
			}
		}(packet, chRawTour)

		// Read tours from XML stream and send them to them channel
		err := sletat.FetchTours(packet.Id, chRawTour)
		if err != nil {
			log.Println(err)
		}

		close(chRawTour)
	}
}

func processTour(packet sletat.PacketInfo, tour *sletat.Tour) {
	tour.CreateDate = packet.CreateDate

	tour.DptCityId = packet.DptCityId

	if operator, ok := operators[tour.SourceId]; ok {
		tour.PriceByr = currencyPrice(tour.Price, operator.ExchangeRateRur)
		tour.PriceEur = currencyPrice(tour.Price, operator.ExchangeRateRur)
		tour.PriceEur = currencyPrice(tour.Price, operator.ExchangeRateRur)
	}
}

func currencyPrice(price int, exchange float64) int {
	return price
	//return price * int(exchange)
}

func saveTours(chTour <-chan sletat.Tour) {
	tours := make([]sletat.Tour, 0, 64)

	for tour := range chTour {
		tours = append(tours, tour)
		if len(tours) == 64 {
			db.SaveTours(tours)
			tours = make([]sletat.Tour, 0, 256)
		}
	}
}
