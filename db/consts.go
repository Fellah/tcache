package db

const (
	toursFields = `
		source_id, price, currency_id, checkin,
		nights, adults, kids, hotel_id,
		town_id, meal_id,
		created_at, updated_at,
		dpt_city_id, country_id, price_byr,
		price_eur, price_usd, active, kid1age,
		kid2age, kid3age,
		tickets_included, has_econom_tickets_dpt, has_econom_tickets_rtn, hotel_is_in_stop
	`

	toursFieldsPartition = toursFields + `, updated_price`

	toursFieldsEHI = `
		source_id, price, currency_id, checkin,
		nights, adults, kids,
		town_id, meal_id, created_at,
		updated_at, dpt_city_id, country_id, price_byr,
		price_eur, price_usd, active, kid1age,
		kid2age, kid3age
	`

	toursUnique = `
		hotel_id, checkin, dpt_city_id, nights,
		adults, kids, meal_id,
		kid1age, kid2age, kid3age
	`

	toursUniqueEHI = `
		country_id, town_id, checkin, dpt_city_id,
		nights, adults, meal_id, kids,
		kid1age, kid2age, kid3age
	`

	toursValues = `
		%d, %d, %d, '%s',
		%d, %d, %d, %d,
		%d, %d,
		NOW(), NOW(),
		%d, %d, %d, %d,
		%d, %t, %d, %d,
		%d,
		'%d', '%d', '%d', '%d'
	`

	toursValuesPartition = toursValues + `, %s`

	toursValuesEHI = `
		%d, %d, %d, '%s',
		%d, %d, %d,
		%d, %d, '%s', '%s',
		%d, %d, %d, %d,
		%d, %t, %d, %d,
		%d
	`

	toursUpdate = `
		price = EXCLUDED.price,
		source_id = EXCLUDED.source_id,
		tickets_included = EXCLUDED.tickets_included,
		has_econom_tickets_dpt = EXCLUDED.has_econom_tickets_dpt,
		has_econom_tickets_rtn = EXCLUDED.has_econom_tickets_rtn,
		hotel_is_in_stop = EXCLUDED.hotel_is_in_stop,
		updated_at = NOW()
	`
)
