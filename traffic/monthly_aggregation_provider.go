package traffic

import (
	"forecasting/core/types"
	"forecasting/rules"
	"github.com/golang-module/carbon/v2"
)

type monthlyAggregationProvider struct {
}

func NewMonthlyAggregationProvider() *monthlyAggregationProvider {
	return &monthlyAggregationProvider{}
}

func (p *monthlyAggregationProvider) GetLast(
	forecastRule *rules.ForecastRule,
	period types.Period,
) []MonthlyAggregationRecord {

	return _TrafficRecordsWithVolume
}

func (p *monthlyAggregationProvider) Get(
	forecastRule *rules.ForecastRule,
	period types.Period,
) []MonthlyAggregationRecord {
	return _HistoricalRecordsInForecastedPeriod
}

var _TrafficRecordsWithVolume = []MonthlyAggregationRecord{
	{
		VolumeActual: 50.0,
		Month:        carbon.Parse("2023-08-01").ToDateStruct(),
	},
	{
		VolumeActual: 100.0,
		Month:        carbon.Parse("2023-09-01").ToDateStruct(),
	},
	{
		VolumeActual: 20.0,
		Month:        carbon.Parse("2023-10-01").ToDateStruct(),
	},
}

var _ZeroVolumeTrafficRecords = []MonthlyAggregationRecord{
	{
		VolumeActual: 0.0,
		Month:        carbon.Parse("2023-08-01").ToDateStruct(),
	},
	{
		VolumeActual: 0.0,
		Month:        carbon.Parse("2023-09-01").ToDateStruct(),
	},
	{
		VolumeActual: 0.0,
		Month:        carbon.Parse("2023-10-01").ToDateStruct(),
	},
}

var _HistoricalRecordsInForecastedPeriod = []MonthlyAggregationRecord{
	{
		VolumeActual: 500.0,
		Month:        carbon.Parse("2024-08-01").ToDateStruct(),
	},
}
