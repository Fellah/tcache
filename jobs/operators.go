package jobs

import (
	"log"

	"github.com/fellah/tcache/db"
)

type Operator struct {
	ExchangeRateUsd float64
	ExchangeRateEur float64
	ExchangeRateRur float64
}

var operators map[int]Operator

func QueryOperators() {
	rawOperators, err := db.QueryOperators()
	if err != nil {
		log.Fatal(err)
	}

	operators = make(map[int]Operator)
	for _, rawOperator := range rawOperators {
		operators[rawOperator.Id] = Operator{
			ExchangeRateUsd: parseExchangeRateValue(rawOperator.ExchangeRateUsd),
			ExchangeRateEur: parseExchangeRateValue(rawOperator.ExchangeRateEur),
			ExchangeRateRur: parseExchangeRateValue(rawOperator.ExchangeRateRur),
		}
	}
}

func parseExchangeRateValue(v interface{}) float64 {
	switch v.(type) {
	case float64:
		return v.(float64)
	default:
		return 0
	}
}
