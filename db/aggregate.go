package db

import (
	"github.com/fellah/tcache/log"
)

func MakeAggregation() {
	fields := []string{
		"country_id",
		"town_id",
		"hotel_id",
	}

	for _, field := range fields {
		makeAggregation(field)
	}
}

func makeAggregation(aggrField string) {
	rows, err := db.Query(`
	SELECT
		` + aggrField + `,
		MIN(price) as price,
		price_byr,
		price_eur,
		price_usd
	FROM cached_sletat_tours
	GROUP BY ` + aggrField + `, price_byr, price_eur, price_usd
	ORDER BY min(price) ASC`)
	if err != nil {
		log.Error.Println(err)
		return
	}

	defer rows.Close()

	var aggrId int
	var price int
	var priceByr int64
	var priceEur int
	var priceUsd int

	for rows.Next() {
		err = rows.Scan(&aggrId, &price, &priceByr, &priceEur, &priceUsd)
		if err != nil {
			log.Error.Println(err)
			continue
		}

		saveAggregation(aggrId, aggrField, price, priceByr, priceEur, priceUsd)
	}
}

func saveAggregation(aggr_id int, t string, price int, priceByr int64, priceEur, priceUsd int) {
	var id int

	err := db.QueryRow(`
	INSERT INTO agregate_data_for_cached_sletat_tours(
		agregate_item_id,
		agregate_for_type,
		price,
		price_byr,
		price_eur,
		price_usd
	) VALUES($1, $2, $3, $4, $5, $6)
	ON CONFLICT (agregate_item_id, agregate_for_type) DO UPDATE SET
		price = $3,
		price_byr = $4,
		price_eur = $5,
		price_usd = $6
	RETURNING id`, aggr_id, t, price, priceByr, priceEur, priceUsd).Scan(&id)
	if err != nil {
		log.Error.Println(err)
	}
}
