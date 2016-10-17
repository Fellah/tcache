package prefilter

import (
	"github.com/fellah/tcache/data"
	"github.com/fellah/tcache/db"
	"github.com/fellah/tcache/log"
)

var hotelsNAPIds []int
var hotelsNAIds []int
var townsIds []int

// Public

func PrepareData() {
	var err error

	hotelsNAPIds, err = db.QueryHotelsNameActivePictures()
	if err != nil {
		log.Error.Fatal(err)
	}

	hotelsNAIds, err = db.QueryHotelsNameActiveNoImages()
	if err != nil {
		log.Error.Fatal(err)
	}

	townsIds, err = db.QueryCities()
	if err != nil {
		log.Error.Fatal(err)
	}
}

func ForPartnersTours(tour *data.Tour) bool {
	return (tour.TicketsIncluded != 0 &&
		(tour.HasEconomTicketsDpt == 1 || tour.HasEconomTicketsDpt == 2) &&
		(tour.HasEconomTicketsRtn == 1 || tour.HasEconomTicketsRtn == 2) &&
		(tour.HotelIsInStop == 0 || tour.HotelIsInStop == 2) &&
		tour.HotelId != 0 &&
		isTownGood(tour.TownId))
}

func IsHotelNameActivePictures(hotelId int) bool {
	return isInListInt(hotelsNAPIds, hotelId)
}

func IsHotelNameActive(hotelId int) bool {
	return (isInListInt(hotelsNAPIds, hotelId) || isInListInt(hotelsNAIds, hotelId))
}

func isTownGood(townId int) bool {
	return isInListInt(townsIds, townId)
}

func isInListInt(list []int, id int) bool {
	for _, goodId := range list {
		if goodId == id {
			return true
		}
	}

	return false
}
