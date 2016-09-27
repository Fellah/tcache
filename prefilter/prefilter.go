package prefilter

import (
	"github.com/fellah/tcache/data"
	"github.com/fellah/tcache/db"
	"github.com/fellah/tcache/log"
)

var hotelsIds []int

func PrepareData() {
	var err error

	hotelsIds, err = db.QueryHotels()
	if err != nil {
		log.Error.Fatal(err)
	}
}

func TourEnable(tour *data.Tour) bool {
	return isHotelGood(tour.HotelId)
}

func isHotelGood(hotelId int) bool {
	for _, goodHotelId := range hotelsIds {
		if goodHotelId == hotelId {
			return true
		}
	}

	return false
}

