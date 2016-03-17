package db

import (
	"log"
)

type Operator struct {
	Id              int
	ExchangeRateUsd interface{}
	ExchangeRateEur interface{}
	ExchangeRateRur interface{}
}

func QueryOperators() ([]Operator, error) {
	rows, err := db.Query("SELECT sletat_tour_operator_id, exchange_rate_usd, exchange_rate_eur, exchange_rate_rur FROM sletat_tour_operators")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	operators := []Operator{}

	for rows.Next() {
		var id int
		var exchangeRateUsd interface{}
		var exchangeRateEur interface{}
		var exchangeRateRur interface{}

		err = rows.Scan(
			&id,
			&exchangeRateUsd,
			&exchangeRateEur,
			&exchangeRateRur,
		)
		if err != nil {
			log.Println(err)
		}

		operator := Operator{
			Id:              id,
			ExchangeRateUsd: exchangeRateUsd,
			ExchangeRateEur: exchangeRateEur,
			ExchangeRateRur: exchangeRateRur,
		}

		operators = append(operators, operator)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return operators, nil
}
