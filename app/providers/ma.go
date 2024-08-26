package providers

import (
	"forecasting/app/domain/models"
	"forecasting/core/types"
)

type MonthlyAggregationProvider interface {
	GetLast(forecastRule *models.ForecastRule, period types.Period) []models.MonthlyAggregationRecord

	Get(forecastRule *models.ForecastRule, period types.Period) []models.MonthlyAggregationRecord
}
