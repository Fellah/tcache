package db

import (
	"github.com/fellah/tcache/log"
)

func QueryCities() ([]int, error) {
	rows, err := db.Query("SELECT sletat_city_id FROM sletat_cities WHERE active = true")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var cityId int
	citiesIds := make([]int, 0)

	for rows.Next() {
		err = rows.Scan(&cityId)
		if err != nil {
			log.Error.Println(err)
		}

		citiesIds = append(citiesIds, cityId)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return citiesIds, nil
}

func QueryDepartCities() ([]int, error) {
	rows, err := db.Query("SELECT sletat_depart_city_id FROM sletat_depart_cities WHERE active = true")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var cityId int
	citiesIds := make([]int, 0)

	for rows.Next() {
		err = rows.Scan(&cityId)
		if err != nil {
			log.Error.Println(err)
		}

		citiesIds = append(citiesIds, cityId)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return citiesIds, nil
}
