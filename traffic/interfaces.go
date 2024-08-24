package traffic

import (
	"forecasting/core/types"
	"forecasting/rules"
)

type MonthlyAggregationProvider interface {
	GetLast(forecastRule *rules.ForecastRule, period types.Period) []MonthlyAggregationRecord

	Get(forecastRule *rules.ForecastRule, period types.Period) []MonthlyAggregationRecord
}

type (
	BudgetTrafficOptions struct {
		ForecastRule   *rules.ForecastRule
		Period         *types.Period
		HistoricalOnly bool
	}

	BudgetTrafficProvider interface {
		Get(options BudgetTrafficOptions) []BudgetTrafficRecord

		CreateMany(records []BudgetTrafficRecord)

		ClearForecasted()
	}
)
