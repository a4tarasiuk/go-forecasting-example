package providers

import (
	"forecasting/app/domain/models"
	"forecasting/core/types"
)

type MonthlyAggregationInMemoryProvider struct {
	Records []models.MonthlyAggregationRecord
}

func NewMonthlyAggregationInMemoryProvider(
	records []models.MonthlyAggregationRecord,
) *MonthlyAggregationInMemoryProvider {

	return &MonthlyAggregationInMemoryProvider{Records: records}
}

func (p *MonthlyAggregationInMemoryProvider) GetLast(
	forecastRule *models.ForecastRule,
	period types.Period,
) []models.MonthlyAggregationRecord {

	return p.Get(forecastRule, period)
}

func (p *MonthlyAggregationInMemoryProvider) Get(
	forecastRule *models.ForecastRule,
	period types.Period,
) []models.MonthlyAggregationRecord {

	var filteredRecords []models.MonthlyAggregationRecord

	for _, r := range p.Records {
		if period.Contains(r.Month) {
			filteredRecords = append(filteredRecords, r)
		}
	}

	return filteredRecords
}
