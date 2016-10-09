package db

import (
	"github.com/fellah/tcache/log"
)

func SavePartnerTour(group_hash string, tour map[string]string) {
	query := `
		INSERT INTO cached_sletat_tours as cst (
			group_hash,
			nights, adults, kids, kid1age, kid2age, kid3age, checkin, dpt_city_id, town_id,
			meal_present, operator_id,
		   	price, hotel_id, tickets_included, has_econom_tickets_dpt, has_econom_tickets_rtn,
		   	hotel_is_in_stop, sletat_request_id, sletat_offer_id,
		   	created_at, updated_at
		)
		VALUES (
			$1,
			$2, $3, $4, $5, $6, $7, $8, $9,	$10,
			$11, $12,
			$13, $14, $15, $16, $17,
			$18, $19, $20,
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
			updated_at = NOW()
	`

	var values []string = []string{
		"\\x"+group_hash,
		tour["nights"],	tour["adults"], tour["kids"], tour["kid1age"], tour["kid2age"],
		tour["kid3age"], tour["checkin"], tour["dpt_city_id"], tour["town_id"],
		tour["meal_present"], tour["operator_id"],
		tour["price"], tour["hotel_id"], tour["tickets_included"], tour["has_econom_tickets_dpt"],
		tour["has_econom_tickets_rtn"], tour["hotel_is_in_stop"], tour["sletat_request_id"],
		tour["sletat_offer_id"],
	}

	err := sendQueryParams(query, values)
	if err != nil {
		log.Error.Println(err)
	}
}
