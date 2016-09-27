package jobs

import (
	"github.com/fellah/go-helpers/log"

	"github.com/fellah/tcache/db"
	"github.com/fellah/tcache/data"
)

var (
	operators map[int]data.Operator
	activeOperatorsIds []int
)

func queryOperators() {
	rawOperators, err := db.QueryOperators()
	if err != nil {
		log.Error.Fatal(err)
	}

	operators = make(map[int]data.Operator)
	for _, rawOperator := range rawOperators {
		operators[rawOperator.Id] = data.Operator{
			ExchangeRateUsd: parseExchangeRateValue(rawOperator.ExchangeRateUsd),
			ExchangeRateEur: parseExchangeRateValue(rawOperator.ExchangeRateEur),
			ExchangeRateRur: parseExchangeRateValue(rawOperator.ExchangeRateRur),
		}
	}

	activeOperatorsIds, err = db.QueryActiveOperators()
	if err != nil {
		log.Error.Fatal(err)
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

func isOperatorActive(operatorId int) bool {
	for _, activeOperatorId := range activeOperatorsIds {
		if activeOperatorId == operatorId {
			return true
		}
	}

	return false
}