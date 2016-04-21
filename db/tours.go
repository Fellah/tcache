package db

import (
	"time"

	"github.com/lib/pq"

	"github.com/fellah/tcache/log"
	"github.com/fellah/tcache/sletat"
)

func RemoveExistTours(t time.Time) {
	_, err := db.Query("DELETE FROM cached_sletat_tours WHERE created_at >= $1", t)
	if err != nil {
		log.Error.Println(err)
	}
}

func SaveTours(tours []sletat.Tour) {
	txn, err := db.Begin()
	if err != nil {
		log.Error.Println(err)
		return
	}

	stmt, err := txn.Prepare(pq.CopyIn(
		"cached_sletat_tours",
		"tour_hash",
		"offer_id",
		"request_id",
		"source_id",
		"update_date",
		"price",
		"currency_id",
		"checkin",
		"nights",
		"adults",
		"kids",
		"kid1age",
		"kid2age",
		"kid3age",
		"hotel_id",
		"town_id",
		"star_id",
		"meal_id",
		"room_id",
		"room_name",
		"htplace_id",
		"htplace_name",
		"hotel_is_in_stop",
		"tickets_included",
		"has_econom_tickets_dpt",
		"has_econom_tickets_rtn",
		"has_business_tickets_dpt",
		"has_business_tickets_rtn",
		"tour_name",
		"original_price",
		"tour_url",
		"price_type",
		"flags",
		"created_at",
		"updated_at",

		"active",

		"dpt_city_id",
		"country_id",

		"price_byr",
		"price_eur",
		"price_usd",
	))
	if err != nil {
		log.Error.Println(err)
		return
	}

	for _, tour := range tours {
		_, err = stmt.Exec(
			tour.Hash,
			tour.OfferId,
			tour.RequestId,
			tour.SourceId,
			tour.UpdateDate,
			tour.Price,
			tour.CurrencyId,
			tour.Checkin,
			tour.Nights,
			tour.Adults,
			tour.Kids,
			tour.Kid1Age,
			tour.Kid2Age,
			tour.Kid3Age,
			tour.HotelId,
			tour.TownId,
			tour.StarId,
			tour.MealId,
			tour.RoomId,
			tour.RoomName,
			tour.HtplaceId,
			tour.HtplaceName,
			tour.HotelIsInStop,
			tour.TicketsIncluded,
			tour.HasEconomTicketsDpt,
			tour.HasEconomTicketsRtn,
			tour.HasBusinessTicketsDpt,
			tour.HasBusinessTicketsRtn,
			tour.TourName,
			tour.OriginalPrice,
			tour.TourUrl,
			tour.PriceType,
			tour.Flags,
			tour.CreateDate,
			tour.UpdateDate,

			true,

			tour.DptCityId,
			tour.CountryId,

			tour.PriceByr,
			tour.PriceEur,
			tour.PriceUsd,
		)
		if err != nil {
			log.Error.Println(err)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Error.Println(err)
		return
	}

	/*err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}*/

	err = txn.Commit()
	if err != nil {
		log.Error.Println(err)
		return
	}
}

func RemoveExpiredTours() {
	t := time.Now().UTC()
	t = t.Add(-48 * time.Hour) // 2 days ago.

	_, err := db.Query("DELETE FROM cached_sletat_tours WHERE created_at <= $1", t)
	if err != nil {
		log.Error.Println(err)
	}
}

func VacuumTours() {
	_, err := db.Exec("VACUUM (FULL, FREEZE, VERBOSE, ANALYZE) cached_sletat_tours")
	if err != nil {
		log.Error.Println(err)
	}
}

func AgregateToursByCountry() {
	rows, err := db.Query(`
	SELECT
		country_id,
		MIN(price) as price,
		price_byr,
		price_eur,
		price_usd
	FROM cached_sletat_tours
	GROUP BY country_id, price_byr, price_eur, price_usd
	ORDER BY min(price) ASC`)
	if err != nil {
		log.Error.Println(err)
		return
	}

	defer rows.Close()

	var countryId int
	var price int
	var priceByr int64
	var priceEur int
	var priceUsd int

	for rows.Next() {
		err = rows.Scan(&countryId, &price, &priceByr, &priceEur, &priceUsd)
		if err != nil {
			log.Error.Println(err)
			continue
		}

		if isAgregateToursExist(countryId, "country") {
			err := db.QueryRow(`
			UPDATE agregate_data_for_cached_sletat_tours
			SET price = $3, price_byr = $4, price_eur = $5, price_usd = $5
			WHERE agregate_item_id = $1 AND agregate_for_type = $2`, countryId, "country", price, priceByr, priceEur, priceUsd)
			if err != nil {
				//log.Debug.Println(err)
				continue
			}
		} else {
			err := db.QueryRow(`
			INSERT INTO agregate_data_for_cached_sletat_tours(
				agregate_item_id,
				agregate_for_type,
				price,
				price_byr,
				price_eur,
				price_usd
			) VALUES($1, $2, $3, $4, $5, $6)`, countryId, "country", price, priceByr, priceEur, priceUsd)
			if err != nil {
				//log.Debug.Println(err)
				continue
			}
		}
	}
}

func AgregateToursByHotel() {
	rows, err := db.Query(`
	SELECT
		hotel_id,
		MIN(price) as price,
		price_byr,
		price_eur,
		price_usd
	FROM cached_sletat_tours
	GROUP BY country_id, price_byr, price_eur, price_usd
	ORDER BY min(price) ASC`)
	if err != nil {
		log.Error.Println(err)
		return
	}

	defer rows.Close()

	var hotelId int
	var price int
	var priceByr int64
	var priceEur int
	var priceUsd int

	for rows.Next() {
		err = rows.Scan(&hotelId, &price, &priceByr, &priceEur, &priceUsd)
		if err != nil {
			log.Error.Println(err)
			continue
		}

		if isAgregateToursExist(hotelId, "hotel") {
			err := db.QueryRow(`
			UPDATE agregate_data_for_cached_sletat_tours
			SET price = $3, price_byr = $4, price_eur = $5, price_usd = $5
			WHERE agregate_item_id = $1 AND agregate_for_type = $2`, hotelId, "hotel", price, priceByr, priceEur, priceUsd)
			if err != nil {
				//log.Debug.Println(err)
				continue
			}
		} else {
			err := db.QueryRow(`
			INSERT INTO agregate_data_for_cached_sletat_tours(
				agregate_item_id,
				agregate_for_type,
				price,
				price_byr,
				price_eur,
				price_usd
			) VALUES($1, $2, $3, $4, $5, $6)`, hotelId, "hotel", price, priceByr, priceEur, priceUsd)
			if err != nil {
				//log.Debug.Println(err)
				continue
			}
		}
	}
}

func isAgregateToursExist(id int, t string) bool {
	res := false

	rows, err := db.Query(`
		SELECT COUNT(agregate_item_id)
		FROM agregate_data_for_cached_sletat_tours
		WHERE agregate_item_id = $1 AND agregate_for_type = $2`, id, t)
	if err != nil {
		log.Error.Println(err)
		return res
	}

	defer rows.Close()

	var count int

	for rows.Next() {
		rows.Scan(&count)
	}

	if count > 0 {
		res = true
	}

	return res
}
