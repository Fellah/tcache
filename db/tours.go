package db

import (
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/fellah/tcache/log"
	"github.com/fellah/tcache/sletat"
)

func SaveTours(tours []sletat.Tour) {
	txn, err := db.Begin()
	if err != nil {
		log.Error.Println(err)
		return
	}

	values := makeToursValues(tours)
	query := fmt.Sprintf(`
		INSERT INTO cached_sletat_tours(%s)
		VALUES (%s)
		ON CONFLICT (offer_id, request_id) DO UPDATE SET %s
	`, toursFields, values, toursUpdate)

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
			strEscape(tour.Hash), tour.OfferId, tour.RequestId, tour.SourceId,
			strEscape(tour.UpdateDate), tour.Price, tour.CurrencyId, strEscape(tour.Checkin),
			tour.Nights, tour.Adults, tour.Kids, tour.Kid1Age,
			tour.Kid2Age, tour.Kid3Age, tour.HotelId, tour.TownId,
			tour.StarId, tour.MealId, tour.RoomId, strEscape(tour.RoomName),
			tour.HtplaceId, tour.HotelIsInStop, tour.TicketsIncluded, tour.HasEconomTicketsDpt,
			tour.HasEconomTicketsRtn, tour.HasBusinessTicketsDpt, tour.HasBusinessTicketsRtn, strEscape(tour.TourName),
			tour.OriginalPrice, strEscape(tour.TourUrl), tour.PriceType, tour.Flags,
			strEscape(tour.CreateDate), strEscape(tour.UpdateDate), true, tour.DptCityId,
			tour.CountryId, tour.PriceByr, tour.PriceEur, tour.PriceUsd,
		)
	}

	return strings.Join(values, "), (")
}

func strEscape(s string) string {
	return "'" + strings.Replace(s, "'", "''", -1) + "'"
}

func RemoveExistTours(t time.Time) {
	_, err := db.Query("DELETE FROM cached_sletat_tours WHERE created_at >= $1", t)
	if err != nil {
		log.Error.Println(err)
	}
}

func RemoveExpiredTours() {
	t := time.Now().UTC().Add(-36 * time.Hour)

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
