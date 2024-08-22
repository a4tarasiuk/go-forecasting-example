package forecast_models

import (
	"errors"

	"forecasting/calculation"
	"forecasting/core/types"
	"forecasting/rules"
	"forecasting/traffic"
)

type manualVolume struct {
	trafficProvider traffic.MonthlyAggregationProvider
}

func NewManualVolume() manualVolume {
	return manualVolume{trafficProvider: traffic.NewMonthlyAggregationProvider()}
}

func (model *manualVolume) Calculate(forecastRule *rules.ForecastRule) ([]calculation.ForecastRecord, error) {
	if forecastRule.LHM == nil {
		return model.calculateWithoutTraffic(forecastRule), nil
	}

	trafficRecords := model.trafficProvider.GetLast(forecastRule, nil)

	if traffic.ShouldCalculateWithoutTraffic(trafficRecords) {
		return model.calculateWithoutTraffic(forecastRule), nil
	}

	if forecastRule.Period.Contains(*forecastRule.LHM) {
		adjustedForecastedVolume, err := model.extractHistoricalVolumeFromForecasted(forecastRule)

		if err != nil {
			return nil, err
		}

		forecastRule.Volume = adjustedForecastedVolume
	}

	forecastRecords := model.calculateWithTraffic(forecastRule, trafficRecords)

	return forecastRecords, nil
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

func (model *manualVolume) extractHistoricalVolumeFromForecasted(forecastRule *rules.ForecastRule) (float64, error) {
	historicalPeriodInForecasted := types.NewPeriod(forecastRule.Period.StartDate, *forecastRule.LHM)

	trafficRecords := model.trafficProvider.Get(forecastRule, &historicalPeriodInForecasted)

	totalHistoricalVolume := traffic.CalculateTotalHistoricalVolume(trafficRecords)

	if totalHistoricalVolume > forecastRule.Volume {
		return 0, errors.New("historical volume exceeds forecasted volume")
	}

	adjustedForecastVolume := forecastRule.Volume - totalHistoricalVolume

	return adjustedForecastVolume, nil
}

func (model *manualVolume) calculateWithTraffic(
	forecastRule *rules.ForecastRule,
	trafficRecords []traffic.MonthlyAggregationRecord,
) []calculation.ForecastRecord {

	totalHistoricalVolume := traffic.CalculateTotalHistoricalVolume(trafficRecords)

	totalForecastedVolume := forecastRule.Volume

	totalForecastedMonths := forecastRule.Period.GetTotalMonths()

	historicalTrafficMonthMap := make(map[string]traffic.MonthlyAggregationRecord, totalForecastedMonths)

	for _, record := range trafficRecords {
		historicalTrafficMonthMap[record.Month.ToDateString()] = record
	}

	forecastRecords := make([]calculation.ForecastRecord, totalForecastedMonths)

	for idx, forecastedMonth := range forecastRule.Period.GetMonths() {
		historicalRecord := historicalTrafficMonthMap[forecastedMonth.SubYear().ToDateString()]

		forecastedVolumeActual := (historicalRecord.VolumeActual / totalHistoricalVolume) * totalForecastedVolume

		forecastRecord := calculation.ForecastRecord{
			VolumeActual: forecastedVolumeActual,
			Month:        forecastedMonth,
		}

		forecastRecords[idx] = forecastRecord
	}

	return forecastRecords
}
