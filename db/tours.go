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

	stmt, err := txn.Prepare(`INSERT INTO cached_sletat_tours (
		tour_hash,
		offer_id,
		request_id,
		source_id,
		update_date,
		price,
		currency_id,
		checkin,
		nights,
		adults,
		kids,
		kid1age,
		kid2age,
		kid3age,
		hotel_id,
		town_id,
		star_id,
		meal_id,
		room_id,
		room_name,
		htplace_id,
		htplace_name,
		hotel_is_in_stop,
		tickets_included,
		has_econom_tickets_dpt,
		has_econom_tickets_rtn,
		has_business_tickets_dpt,
		has_business_tickets_rtn,
		tour_name,
		original_price,
		tour_url,
		price_type,
		flags,
		created_at,
		updated_at,

		active,

		dpt_city_id,
		country_id,

		price_byr,
		price_eur,
		price_usd
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		$11,
		$12,
		$13,
		$14,
		$15,
		$16,
		$17,
		$18,
		$19,
		$20,
		$21,
		$22,
		$23,
		$24,
		$25,
		$26,
		$27,
		$28,
		$29,
		$30,
		$31,
		$32,
		$33,
		$34,
		$35,

		$36,

		$37,
		$38,

		$39,
		$40,
		$41
	)ON CONFLICT (offer_id, request_id) DO UPDATE SET
		tour_hash = $1,
		source_id = $4,
		update_date = $5,
		price = $6,
		currency_id = $7,
		checkin = $8,
		nights = $9,
		adults = $10,
		kids = $11,
		kid1age = $12,
		kid2age = $13,
		kid3age = $14,
		hotel_id = $15,
		town_id = $16,
		star_id = $17,
		meal_id = $18,
		room_id = $19,
		room_name = $20,
		htplace_id = $21,
		htplace_name = $22,
		hotel_is_in_stop = $23,
		tickets_included = $24,
		has_econom_tickets_dpt = $25,
		has_econom_tickets_rtn = $26,
		has_business_tickets_dpt = $27,
		has_business_tickets_rtn = $28,
		tour_name = $29,
		original_price = $30,
		tour_url = $31,
		price_type = $32,
		flags = $33,
		created_at = $34,
		updated_at = $35,

		active = $36,

		dpt_city_id = $37,
		country_id = $38,

		price_byr = $39,
		price_eur = $40,
		price_usd = $41`)
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

	err = txn.Commit()
	if err != nil {
		log.Error.Println(err)
		return
	}

	err = stmt.Close()
	if err != nil {
		log.Error.Println(err)
	}
}

func SaveToursV2(tours []sletat.Tour) {
	txn, err := db.Begin()
	if err != nil {
		log.Error.Println(err)
		return
	}

	stmt, err := txn.Prepare(`INSERT INTO cached_sletat_tours (
		tour_hash,
		offer_id,
		request_id,
		source_id,
		update_date,
		price,
		currency_id,
		checkin,
		nights,
		adults,
		kids,
		kid1age,
		kid2age,
		kid3age,
		hotel_id,
		town_id,
		star_id,
		meal_id,
		room_id,
		room_name,
		htplace_id,
		htplace_name,
		hotel_is_in_stop,
		tickets_included,
		has_econom_tickets_dpt,
		has_econom_tickets_rtn,
		has_business_tickets_dpt,
		has_business_tickets_rtn,
		tour_name,
		original_price,
		tour_url,
		price_type,
		flags,
		created_at,
		updated_at,

		active,

		dpt_city_id,
		country_id,

		price_byr,
		price_eur,
		price_usd
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		$11,
		$12,
		$13,
		$14,
		$15,
		$16,
		$17,
		$18,
		$19,
		$20,
		$21,
		$22,
		$23,
		$24,
		$25,
		$26,
		$27,
		$28,
		$29,
		$30,
		$31,
		$32,
		$33,
		$34,
		$35,

		$36,

		$37,
		$38,

		$39,
		$40,
		$41
	)`)
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

	err = txn.Commit()
	if err != nil {
		log.Error.Println(err)
		return
	}

	err = stmt.Close()
	if err != nil {
		log.Error.Println(err)
	}
}

func SaveToursV3(tours []sletat.Tour) {
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
	t = t.Add(-36 * time.Hour)

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
