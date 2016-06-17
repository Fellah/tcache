package db

const (
	toursFields = `
		tour_hash, offer_id, request_id, source_id,
		update_date, price, currency_id, checkin,
		nights, adults, kids, kid1age,
		kid2age, kid3age, hotel_id, town_id,
		star_id, meal_id, room_id, room_name,
		htplace_id, hotel_is_in_stop, tickets_included, has_econom_tickets_dpt,
		has_econom_tickets_rtn, has_business_tickets_dpt, has_business_tickets_rtn, tour_name,
		original_price, tour_url, price_type, flags,
		created_at, updated_at, active, dpt_city_id,
		country_id, price_byr, price_eur, price_usd
	`

	toursValues = `
		%s, %d, %d, %d, %s, %d, %d, %s,
		%d, %d, %d, %d, %d, %d, %d, %d,
		%d, %d, %d, %s, %d, %d, %d, %d,
		%d, %d, %d, %s, %d, %s, %d, %d,
		%s, %s, %t, %d, %d, %d, %d, %d
	`

	/*toursValues = `
		%q, %d, %d, %d, %q, %d, %d, %q,
		%d, %d, %d, %d, %d, %d, %d, %d,
		%d, %d, %d, %q, %d, %d, %d, %d,
		%d, %d, %d, %q, %d, %q, %d, %d,
		%q, %q, %t, %d, %d, %d, %d, %d
	`*/

	toursUpdate = `
		tour_hash = EXCLUDED.tour_hash, source_id = EXCLUDED.source_id,
		update_date = EXCLUDED.update_date, price = EXCLUDED.price,
		currency_id = EXCLUDED.currency_id, checkin = EXCLUDED.checkin,
		nights = EXCLUDED.nights, adults = EXCLUDED.adults,
		kids = EXCLUDED.kids, kid1age = EXCLUDED.kid1age,
		kid2age = EXCLUDED.kid2age, kid3age = EXCLUDED.kid3age,
		hotel_id = EXCLUDED.hotel_id, town_id = EXCLUDED.town_id,
		star_id = EXCLUDED.star_id, meal_id = EXCLUDED.meal_id,
		room_id = EXCLUDED.room_id, room_name = EXCLUDED.room_name,
		htplace_id = EXCLUDED.htplace_id, hotel_is_in_stop = EXCLUDED.hotel_is_in_stop,
		tickets_included = EXCLUDED.tickets_included, has_econom_tickets_dpt = EXCLUDED.has_econom_tickets_dpt,
		has_econom_tickets_rtn = EXCLUDED.has_econom_tickets_rtn, has_business_tickets_dpt = EXCLUDED.has_business_tickets_dpt,
		has_business_tickets_rtn = EXCLUDED.has_business_tickets_rtn, tour_name = EXCLUDED.tour_name,
		original_price = EXCLUDED.original_price, tour_url = EXCLUDED.tour_url,
		price_type = EXCLUDED.price_type, flags = EXCLUDED.flags,
		created_at = EXCLUDED.created_at, updated_at = EXCLUDED.updated_at,
		active = EXCLUDED.active, dpt_city_id = EXCLUDED.dpt_city_id,
		country_id = EXCLUDED.country_id, price_byr = EXCLUDED.price_byr,
		price_eur = EXCLUDED.price_eur, price_usd = EXCLUDED.price_usd
	`
)
