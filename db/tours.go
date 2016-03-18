package db

import (
	"log"
	"time"

	"github.com/lib/pq"

	"github.com/fellah/tcache/sletat"
)

func RemoveExistTours(t time.Time) {
	_, err := db.Query(`DELETE FROM cached_sletat_tours WHERE created_at >= $1`, t)
	if err != nil {
		log.Println(err)
	}
}

func SaveTours(tours []sletat.Tour) {
	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
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
		"hotel_id",
		"town_id",
		"star_id",
		"meal_id",
		"htplace_id",
		"hotel_is_in_stop",
		"tickets_included",
		"has_econom_tickets_dpt",
		"has_econom_tickets_rtn",
		"has_business_tickets_dpt",
		"has_business_tickets_rtn",
		"tour_name",
		"hotel_name",
		"town_name",
		"star_name",
		"meal_name",
		"room_name",
		"htplace_name",
		"original_hotel_name",
		"original_town_name",
		"original_star_name",
		"original_meal_name",
		"original_room_name",
		"original_htplace_name",
		"original_country_name",
		"original_dptcity_name",
		"original_currency_id",
		"original_currency_name",
		"original_price",
		"tour_url",
		"price_type",
		"flags",
		"created_at",
		"updated_at",

		"dpt_city_id",

		"price_byr",
		"price_eur",
		"price_usd",

		//"uploader_id",
	))
	if err != nil {
		log.Fatal(err)
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
			tour.HotelId,
			tour.TownId,
			tour.StarId,
			tour.MealId,
			tour.HtplaceId,
			tour.HotelIsInStop,
			tour.TicketsIncluded,
			tour.HasEconomTicketsDpt,
			tour.HasEconomTicketsRtn,
			tour.HasBusinessTicketsDpt,
			tour.HasBusinessTicketsRtn,
			tour.TourName,
			tour.HotelName,
			tour.TownName,
			tour.StarName,
			tour.MealName,
			tour.RoomName,
			tour.HtplaceName,
			tour.OriginalHotelName,
			tour.OriginalTownName,
			tour.OriginalStarName,
			tour.OriginalMealName,
			tour.OriginalRoomName,
			tour.OriginalHtplaceName,
			tour.OriginalCountryName,
			tour.OriginalDptCityName,
			tour.OriginalCurrencyId,
			tour.OriginalCurrencyName,
			tour.OriginalPrice,
			tour.TourUrl,
			tour.PriceType,
			tour.Flags,
			tour.CreateDate,
			tour.UpdateDate,

			tour.DptCityId,

			tour.PriceByr,
			tour.PriceEur,
			tour.PriceUsd,

			// "uploader_id"
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	/*err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}*/

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
