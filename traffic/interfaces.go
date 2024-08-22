package traffic

import (
	"forecasting/core/types"
	"forecasting/rules"
)

type MonthlyAggregationProvider interface {
	GetLast(forecastRule *rules.ForecastRule, period *types.Period) []MonthlyAggregationRecord

	Get(forecastRule *rules.ForecastRule, period *types.Period) []MonthlyAggregationRecord
}
