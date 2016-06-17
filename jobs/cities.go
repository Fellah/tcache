package jobs

import (
	"sort"

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

func isCityActive(city_id int) bool {
	idx := sort.SearchInts(citiesIds, city_id)
	return !(idx == len(citiesIds))
}

func isDepartCityActive(city_id int) bool {
	idx := sort.SearchInts(departCitiesIds, city_id)
	return !(idx == len(departCitiesIds))
}