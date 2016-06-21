package jobs

import (
	"github.com/fellah/tcache/db"
	"github.com/fellah/tcache/log"
)

var citiesIds []int
var departCitiesIds []int

func queryCities() {
	var err error

	citiesIds, err = db.QueryCities()
	if err != nil {
		log.Error.Fatal(err)
	}

	departCitiesIds, err = db.QueryDepartCities()
	if err != nil {
		log.Error.Fatal(err)
	}
}

func isCityActive(cityId int) bool {
	for _, activeCityId := range citiesIds {
		if activeCityId == cityId {
			return true
		}
	}

	return false
}

func isDepartCityActive(cityId int) bool {
	for _, activeCityId := range departCitiesIds {
		if activeCityId == cityId {
			return true
		}
	}

	return false
}