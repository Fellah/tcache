package db

import (
	"fmt"
	"strings"

	_ "github.com/lib/pq"

	"github.com/fellah/tcache/data"
	"github.com/fellah/tcache/log"
	"database/sql"
)

func SaveTours(tours []data.Tour) {
	if len(tours) == 0 {
		return
	}

	filteredTours := removeDuplicates(tours, isEqual)

	{
		values := makeToursValues(filteredTours)
		query := fmt.Sprintf(`
		INSERT INTO cached_sletat_tours as cst (%s)
		VALUES (%s)
		ON CONFLICT (%s)
		DO UPDATE SET %s
		`, toursFields, values, toursUnique, toursUpdate)

		if err := sendQuery(query); err != nil {
			log.Error.Println(err)
		}
	}
}


func SaveMapTour(tour map[string]string, transaction *sql.Tx) {
	if len(tour) == 0 {
		return
	}

	values := makeMapToursValues(tour)
	query := fmt.Sprintf(`
		INSERT INTO cached_sletat_tours as cst (%s)
		VALUES (%s)
		ON CONFLICT (%s)
		DO UPDATE SET %s
		`, toursFields, values, toursUnique, toursUpdate)
	err := SendQueryParamsRaw(transaction, query)

	if err != nil {
		log.Error.Println(err)
	}
}


func CleanMapTours() {
	tx, err := StartTransaction()
	if err != nil {
		log.Error.Println(err)
		return
	}

	SendQueryParamsRaw(tx, "DELETE FROM cached_sletat_tours WHERE checkin < (NOW() - '1 day'::interval)")
	SendQueryParamsRaw(tx, "DELETE FROM cached_sletat_tours WHERE updated_at < (NOW() - '1 hours'::interval)")

	CommitTransaction(tx)
}

func sendQuery(query string) error {
	txn, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := txn.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	if err = stmt.Close(); err != nil {
		return err
	}

	if err = txn.Commit(); err != nil {
		return err
	}

	return nil
}

func makeToursValues(tours []data.Tour) string {
	values := make([]string, len(tours))
	for i, tour := range tours {
		values[i] = fmt.Sprintf(toursValues,
			tour.SourceId, tour.Price, tour.CurrencyId, tour.Checkin,
			tour.Nights, tour.Adults, tour.Kids, tour.HotelId,
			tour.TownId, tour.MealId, tour.CreateDate,
			tour.UpdateDate, tour.DptCityId, tour.CountryId, tour.PriceByr,
			tour.PriceEur, tour.PriceUsd, true, *tour.Kid1Age,
			*tour.Kid2Age, *tour.Kid3Age,
			tour.TicketsIncluded, tour.HasEconomTicketsDpt,
			tour.HasEconomTicketsRtn, tour.HotelIsInStop,
		)
	}

	return strings.Join(values, "), (")
}

func makeMapToursValues(tour map[string]string) string {
	values := fmt.Sprintf(toursValues,
		a2i(tour["source_id"]), a2i(tour["price"]),
		a2i(tour["currency_id"]), tour["checkin"],
		a2i(tour["nights"]), a2i(tour["adults"]), a2i(tour["kids"]),
		a2i(tour["hotel_id"]), a2i(tour["town_id"]), a2i(tour["meal_id"]),
		a2i(tour["dpt_city_id"]), a2i(tour["country_id"]),
		a2i(tour["price_byr"]), a2i(tour["price_eur"]),
		a2i(tour["price_usd"]), true,
		a2i(tour["kid1age"]), a2i(tour["kid2age"]), a2i(tour["kid3age"]),
		a2i(tour["tickets_included"]), a2i(tour["has_econom_tickets_dpt"]),
		a2i(tour["has_econom_tickets_rtn"]), a2i(tour["hotel_is_in_stop"]),
	)

	return values
}

func makeToursValuesPartition(tours []data.Tour) string {
	values := make([]string, len(tours))
	for i, tour := range tours {
		values[i] = fmt.Sprintf(toursValuesPartition,
			tour.SourceId, tour.Price, tour.CurrencyId, tour.Checkin,
			tour.Nights, tour.Adults, tour.Kids, tour.HotelId,
			tour.TownId, tour.MealId, tour.CreateDate,
			tour.UpdateDate, tour.DptCityId, tour.CountryId, tour.PriceByr,
			tour.PriceEur, tour.PriceUsd, true, *tour.Kid1Age,
			*tour.Kid2Age, *tour.Kid3Age, " now()")
	}

	return strings.Join(values, "), (")
}

func makeToursValuesEHI(tours []data.Tour) string {
	values := make([]string, len(tours))
	for i, tour := range tours {
		values[i] = fmt.Sprintf(toursValuesEHI,
			tour.SourceId, tour.Price, tour.CurrencyId, tour.Checkin,
			tour.Nights, tour.Adults, tour.Kids,
			tour.TownId, tour.MealId, tour.CreateDate,
			tour.UpdateDate, tour.DptCityId, tour.CountryId, tour.PriceByr,
			tour.PriceEur, tour.PriceUsd, true, *tour.Kid1Age,
			*tour.Kid2Age, *tour.Kid3Age)
	}

	return strings.Join(values, "), (")
}

func removeDuplicates(tours []data.Tour, isEqual func(data.Tour, data.Tour) bool) []data.Tour {
	remove := make([]bool, len(tours))

	for i, _ := range tours {
		for j, _ := range tours[i+1:] {
			if isEqual(tours[i], tours[i+j+1]) {
				// TODO: if tour.Price < toursBulk[i].Price
				remove[i+j+1] = true
			}
		}
	}

	filteredTours := make([]data.Tour, 0)
	for i, v := range tours {
		if !remove[i] {
			filteredTours = append(filteredTours, v)
		}
	}

	return filteredTours
}

func isEqual(aTour, bTour data.Tour) bool {
	return aTour.HotelId == bTour.HotelId &&
		aTour.Checkin == bTour.Checkin &&
		aTour.DptCityId == bTour.DptCityId &&
		aTour.Nights == bTour.Nights &&
		aTour.Adults == bTour.Adults &&
		aTour.Kids == bTour.Kids &&
		aTour.MealId == bTour.MealId &&
		*aTour.Kid1Age == *bTour.Kid1Age &&
		*aTour.Kid2Age == *bTour.Kid2Age &&
		*aTour.Kid3Age == *bTour.Kid3Age
}

func isEqualEHI(aTour, bTour data.Tour) bool {
	return aTour.CountryId == bTour.CountryId &&
		aTour.TownId == bTour.TownId &&
		aTour.Checkin == bTour.Checkin &&
		aTour.DptCityId == bTour.DptCityId &&
		aTour.Nights == bTour.Nights &&
		aTour.Adults == bTour.Adults &&
		aTour.MealId == bTour.MealId &&
		aTour.Kids == bTour.Kids &&
		*aTour.Kid1Age == *bTour.Kid1Age &&
		*aTour.Kid2Age == *bTour.Kid2Age &&
		*aTour.Kid3Age == *bTour.Kid3Age
}
