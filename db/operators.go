package db

import (
	"github.com/fellah/tcache/log"
)

type Operator struct {
	Id              int
	ExchangeRateUsd interface{}
	ExchangeRateEur interface{}
	ExchangeRateRur interface{}
}

func QueryOperators() ([]Operator, error) {
	rows, err := db.Query(`
	SELECT
		sletat_tour_operator_id,
		exchange_rate_usd,
		exchange_rate_eur,
		exchange_rate_rur
	FROM
	sletat_tour_operators`)
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
			log.Error.Println(err)
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

func QueryActiveOperators() ([]int, error) {
	rows, err := db.Query("SELECT sletat_tour_operator_id FROM sletat_tour_operators WHERE active = true")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var operatorId int
	operatorIds := make([]int, 0)

	for rows.Next() {
		err = rows.Scan(&operatorId)
		if err != nil {
			log.Error.Println(err)
		}

		operatorIds = append(operatorIds, operatorId)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return operatorIds, nil
}