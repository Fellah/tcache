package prefilter

import (
	"github.com/fellah/tcache/data"
	"github.com/fellah/tcache/db"
	"github.com/fellah/tcache/log"
)

var hotelsNAPIds []int
var hotelsNAIds []int
var townsIds []int
var departCitiesIds []int

var townsPartnersIds []int
var departCitiesPartnersIds []int

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

	departCitiesIds, err = db.QueryDepartCities()
	if err != nil {
		log.Error.Fatal(err)
	}

	townsPartnersIds, err = db.QueryPartnersCities()
	if err != nil {
		log.Error.Fatal(err)
	}

	departCitiesPartnersIds, err = db.QueryPartnersDepartCities()
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
		IsTownGood(tour.TownId))
}

func IsHotelNameActivePictures(hotelId int) bool {
	return isInListInt(hotelsNAPIds, hotelId)
}

func IsHotelNameActive(hotelId int) bool {
	return (isInListInt(hotelsNAPIds, hotelId) || isInListInt(hotelsNAIds, hotelId))
}

func IsTownGood(townId int) bool {
	return isInListInt(townsIds, townId)
}

func IsDepartCityActive(dptCityId int) bool {
	return isInListInt(departCitiesIds, dptCityId)
}


func IsPartnersTownGood(townId int) bool {
	return isInListInt(townsPartnersIds, townId)
}

func IsPartnersDepartCityActive(dptCityId int) bool {
	return isInListInt(departCitiesPartnersIds, dptCityId)
}


func isInListInt(list []int, id int) bool {
	for _, goodId := range list {
		if goodId == id {
			return true
		}
	}

	return false
}
