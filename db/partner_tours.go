package db

import (
	"github.com/fellah/tcache/log"
	"strconv"
	"database/sql"
)

func SavePartnerTour(group_hash string, tour map[string]string, transaction *sql.Tx) {
	if tour["checkin"] == "" {
		return
	}

	query := `
		INSERT INTO partners_tours as cst (
			group_hash,
			nights, adults, kids, kid1age, kid2age, kid3age, checkin, dpt_city_id, town_id,
			meal_present, operator_id,
		   	price, hotel_id, tickets_included, has_econom_tickets_dpt, has_econom_tickets_rtn,
		   	hotel_is_in_stop, sletat_request_id, sletat_offer_id,
		   	few_econom_tickets_dpt, few_econom_tickets_rtn, few_places_in_hotel, flags,
		   	description, tour_url, room_name, receiving_party,
		   	created_at, updated_at
		)
		VALUES (
			$1,
			$2, $3, $4, $5, $6, $7, $8, $9,	$10,
			$11, $12,
			$13, $14, $15, $16, $17,
			$18, $19, $20,
			$21, $22, $23, $24,
			$25, $26, $27, $28,
			NOW(), NOW()
		)
		ON CONFLICT (
			group_hash
		)
		DO UPDATE SET
			price = EXCLUDED.price,
			hotel_id = EXCLUDED.hotel_id,
			tickets_included = EXCLUDED.tickets_included,
			has_econom_tickets_dpt = EXCLUDED.has_econom_tickets_dpt,
			has_econom_tickets_rtn = EXCLUDED.has_econom_tickets_rtn,
			hotel_is_in_stop = EXCLUDED.hotel_is_in_stop,
			sletat_request_id = EXCLUDED.sletat_request_id,
			sletat_offer_id = EXCLUDED.sletat_offer_id,
			few_econom_tickets_dpt = EXCLUDED.few_econom_tickets_dpt,
			few_econom_tickets_rtn = EXCLUDED.few_econom_tickets_rtn,
			few_places_in_hotel = EXCLUDED.few_places_in_hotel,
			flags = EXCLUDED.flags,
		   	description = EXCLUDED.description,
		   	tour_url = EXCLUDED.tour_url,
		   	room_name = EXCLUDED.room_name,
		   	receiving_party = EXCLUDED.receiving_party,
			updated_at = NOW()
	`
	err := SendQueryParamsRaw(transaction, query, "\\x" + group_hash,
		a2i(tour["nights"]), a2i(tour["adults"]), a2i(tour["kids"]),
		a2i(tour["kid1age"]), a2i(tour["kid2age"]), a2i(tour["kid3age"]),
		tour["checkin"],
		a2i(tour["dpt_city_id"]), a2i(tour["town_id"]),
		tour["meal_present"],
		a2i(tour["operator_id"]), a2i(tour["price"]),
		a2i(tour["hotel_id"]), a2i(tour["tickets_included"]),
		a2i(tour["has_econom_tickets_dpt"]), a2i(tour["has_econom_tickets_rtn"]),
		a2i(tour["hotel_is_in_stop"]), a2i(tour["sletat_request_id"]),
		a2i64(tour["sletat_offer_id"]), a2i(tour["few_econom_tickets_dpt"]),
		a2i(tour["few_econom_tickets_rtn"]), a2i(tour["few_places_in_hotel"]),
		a2i64(tour["flags"]),
		tour["description"], tour["tour_url"], tour["room_name"], tour["receiving_party"],
	)

	if err != nil {
		log.Error.Println(err)
	}
}

func CleanPartnerTours() {
	tx, err := StartTransaction()
	if err != nil {
		log.Error.Println(err)
		return
	}

	SendQueryParamsRaw(tx, "DELETE FROM partners_tours WHERE checkin < (NOW() - '1 day'::interval)")
	SendQueryParamsRaw(tx, "DELETE FROM partners_tours WHERE updated_at < (NOW() - '1 hour'::interval)")

	CommitTransaction(tx)
}

func a2i(str string) (int) {
	i, err := strconv.Atoi(str)
	if err == nil {
		return i
	}

	return 0
}

func a2i64(str string) (int64) {
	i, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		return i
	}

	return 0
}
