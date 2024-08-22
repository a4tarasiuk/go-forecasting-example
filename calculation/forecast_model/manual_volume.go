package forecast_model

import (
	"forecasting/calculation"
	"forecasting/rules"
	"forecasting/traffic"
	"github.com/go-gota/gota/dataframe"
	"github.com/golang-module/carbon/v2"
)

type manualVolume struct {
}

func NewManualVolume() manualVolume {
	return manualVolume{}
}

func (model *manualVolume) Calculate(forecastRule *rules.ForecastRule) []calculation.ForecastRecord {
	traffic := model.loadTraffic(forecastRule)

	df := dataframe.LoadStructs(traffic)

	if model.shouldCalculateWithoutTraffic(df) {
		return model.calculateWithoutTraffic(forecastRule)
	}

	return []calculation.ForecastRecord{
		{VolumeActual: 50, Month: carbon.Parse("2024-10-01").ToDateStruct()},
	}
}

func (model *manualVolume) loadTraffic(forecastRule *rules.ForecastRule) []traffic.MonthlyAggregationRecord {
	return []traffic.MonthlyAggregationRecord{
		{
			VolumeActual: 0.0,
			// Month:        time.Date(2023, 8, 1, 0, 0, 0, 0, tz),
			Month: "2023-08-01",
		},
		{
			VolumeActual: 0.0,
			// Month:        time.Date(2023, 9, 1, 0, 0, 0, 0, tz),
			Month: "2023-09-01",
		},
		{
			VolumeActual: 0.0,
			// Month:        time.Date(2023, 10, 1, 0, 0, 0, 0, tz),
			Month: "2023-10-01",
		},
	}
}

func (model *manualVolume) shouldCalculateWithoutTraffic(df dataframe.DataFrame) bool {
	rows, _ := df.Dims()

	trafficIsEmpty := rows == 0

	trafficIsZeroVolume := df.Col("VolumeActual").Sum() == 0

	return trafficIsEmpty || trafficIsZeroVolume
}

func (model *manualVolume) calculateWithoutTraffic(forecastRule *rules.ForecastRule) []calculation.ForecastRecord {
	months := forecastRule.Period.GetMonths()

	totalMonths := len(months)

	volumePerRecord := forecastRule.Volume / float64(len(months))

	forecastRecords := make([]calculation.ForecastRecord, totalMonths)

	for idx, month := range months {
		record := calculation.ForecastRecord{VolumeActual: volumePerRecord, Month: month}

		forecastRecords[idx] = record
	}

	return forecastRecords
}
