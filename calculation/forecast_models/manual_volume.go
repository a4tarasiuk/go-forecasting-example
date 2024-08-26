package forecast_models

import (
	"errors"

	"forecasting/calculation"
	"forecasting/core/types"
	"forecasting/rules"
	"forecasting/traffic"
	"forecasting/traffic/persistence"
)

type manualVolume struct {
	trafficProvider traffic.MonthlyAggregationProvider
}

func NewManualVolume() manualVolume {
	return manualVolume{trafficProvider: persistence.NewPostgresMAProvider()}
}

func (model *manualVolume) Calculate(forecastRule *rules.ForecastRule) ([]calculation.ForecastRecord, error) {
	trafficPeriod, err := createTrafficPeriodFromForecasted(forecastRule)

	if err != nil {
		return nil, err
	}

	trafficRecords := model.trafficProvider.GetLast(forecastRule, trafficPeriod)

	if traffic.ShouldCalculateWithoutTraffic(trafficRecords) {
		return model.calculateWithoutTraffic(forecastRule), nil
	}

	if forecastRule.Period.Contains(forecastRule.LHM) {
		adjustedForecastVolume, _err := model.extractHistoricalVolumeFromForecasted(forecastRule)

		if _err != nil {
			return nil, _err
		}

		forecastRule.Volume = adjustedForecastVolume
	}

	forecastRecords := model.calculateWithTraffic(forecastRule, trafficRecords)

	return forecastRecords, nil
}

func (model *manualVolume) calculateWithoutTraffic(forecastRule *rules.ForecastRule) []calculation.ForecastRecord {
	validatedPeriod, _ := forecastRule.GetValidatedPeriod()

	months := validatedPeriod.GetMonths()

	totalMonths := len(months)

	volumePerRecord := forecastRule.Volume / float64(totalMonths)

	forecastRecords := make([]calculation.ForecastRecord, totalMonths)

	for idx, month := range months {
		record := calculation.ForecastRecord{VolumeActual: volumePerRecord, Month: month}

		forecastRecords[idx] = record
	}

	return forecastRecords
}

func (model *manualVolume) extractHistoricalVolumeFromForecasted(forecastRule *rules.ForecastRule) (float64, error) {
	historicalPeriodInForecasted := types.NewPeriod(forecastRule.Period.StartDate, forecastRule.LHM)

	trafficRecords := model.trafficProvider.Get(forecastRule, historicalPeriodInForecasted)

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

	historicalTrafficMonthMap := make(map[string]traffic.MonthlyAggregationRecord, len(trafficRecords))

	for _, record := range trafficRecords {
		historicalTrafficMonthMap[record.Month.ToDateString()] = record
	}

	forecastPeriod, _ := forecastRule.GetValidatedPeriod()

	totalForecastedMonths := forecastPeriod.GetTotalMonths()

	forecastRecords := make([]calculation.ForecastRecord, totalForecastedMonths)

	for idx, forecastMonth := range forecastPeriod.GetMonths() {
		historicalRecord := historicalTrafficMonthMap[forecastMonth.SubYear().ToDateString()]

		forecastedVolumeActual := (historicalRecord.VolumeActual / totalHistoricalVolume) * totalForecastedVolume

		forecastRecord := calculation.ForecastRecord{
			VolumeActual: forecastedVolumeActual,
			Month:        forecastMonth,
		}

		forecastRecords[idx] = forecastRecord
	}

	return forecastRecords
}

func createTrafficPeriodFromForecasted(rule *rules.ForecastRule) (types.Period, error) {
	forecastPeriod, err := rule.GetValidatedPeriod()

	if err != nil {
		return types.Period{}, err
	}

	budgetTrafficPeriod := types.NewPeriod(
		forecastPeriod.StartDate.SubYear().ToDateStruct(),
		forecastPeriod.EndDate.SubYear().ToDateStruct(),
	)

	if forecastPeriod.GetTotalMonths() > 12 {
		budgetTrafficPeriod = types.NewPeriod(
			budgetTrafficPeriod.StartDate,
			budgetTrafficPeriod.StartDate.AddMonths(12-1).ToDateStruct(),
		)
	}

	return budgetTrafficPeriod, nil
}
