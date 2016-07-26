package db

import (
	"fmt"
	"strings"

	_ "github.com/lib/pq"

	"github.com/fellah/tcache/log"
	"github.com/fellah/tcache/sletat"
)

func SaveTours(tours []sletat.Tour) {
	if len(tours) == 0 {
		return
	}

	txn, err := db.Begin()
	if err != nil {
		log.Error.Println(err)
		return
	}

	values := makeToursValues(tours)
	query := fmt.Sprintf(`
		INSERT INTO cached_sletat_tours as cst (%s)
		VALUES (%s)
		ON CONFLICT (%s)
		DO UPDATE SET %s
	`, toursFields, values, toursUnique, toursUpdate)

	stmt, err := txn.Prepare(query)
	if err != nil {
		log.Error.Println(err, query)
		return
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Error.Println(err, query)
		return
	}

	err = stmt.Close()
	if err != nil {
		log.Error.Println(err, query)
		return
	}

	err = txn.Commit()
	if err != nil {
		log.Error.Println(err, query)
		return
	}
}

func makeToursValues(tours []sletat.Tour) string {
	values := make([]string, len(tours))
	for i, tour := range tours {
		values[i] = fmt.Sprintf(toursValues,
			tour.SourceId, tour.Price, tour.CurrencyId, tour.Checkin,
			tour.Nights, tour.Adults, tour.Kids, tour.HotelId,
			tour.TownId, tour.MealId, tour.CreateDate,
			tour.UpdateDate, tour.DptCityId, tour.CountryId, tour.PriceByr,
			tour.PriceEur, tour.PriceUsd, true,  *tour.Kid1Age,
			*tour.Kid2Age, *tour.Kid3Age)
	}

	return strings.Join(values, "), (")
}
