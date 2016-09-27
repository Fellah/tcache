package jobs

import (
	"github.com/fellah/tcache/db"
	"github.com/fellah/tcache/log"
)

var hotelsIds []int

func queryHotels() {
	var err error

	hotelsIds, err = db.QueryHotels()
	if err != nil {
		log.Error.Fatal(err)
	}
}

func IsHotelGood(hotelId int) bool {
	for _, goodHotelId := range hotelsIds {
		if goodHotelId == hotelId {
			return true
		}
	}

	return false
}
