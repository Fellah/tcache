package db

import (
	"fmt"
	"strings"
	//"strconv"

	_ "github.com/lib/pq"

	"github.com/fellah/tcache/data"
	"github.com/fellah/tcache/log"
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

	/*
		{
			partition := "p" + strconv.Itoa(tours[0].CountryId)

			values := makeToursValuesPartition(filteredTours)
			query := fmt.Sprintf(`
			INSERT INTO partitioned_cached_sletat_tours_partitions.%s as cst (%s)
			VALUES (%s)
			ON CONFLICT (%s)
			DO UPDATE SET price = EXCLUDED.price, updated_at = now(), updated_price = now()
			`, partition, toursFieldsPartition, values, toursUnique)

			if err := sendQuery(query); err != nil {
				log.Error.Println(err)
			}
		}*/

	/*
		{
			filteredTours := removeDuplicates(tours, isEqualEHI)

			values := makeToursValuesEHI(filteredTours)
			query := fmt.Sprintf(`
			INSERT INTO cached_sletat_tour_by_cities as cst (%s)
			VALUES (%s)
			ON CONFLICT (%s)
			DO UPDATE SET %s
			`, toursFieldsEHI, values, toursUniqueEHI, toursUpdate)

			if err := sendQuery(query); err != nil {
				log.Error.Println("db:", err)
			}
		}
	*/
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
			*tour.Kid2Age, *tour.Kid3Age)
	}

	return strings.Join(values, "), (")
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
