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
	if forecastRule.LHM == nil {
		// LHM is not set then it will not be calculated
		return model.calculateWithoutTraffic(forecastRule)
	}

	trafficPeriod, err := createTrafficPeriodFromForecasted(forecastRule)

	if err != nil {
		return nil, err
	}

	trafficRecords := model.trafficProvider.GetLast(forecastRule, trafficPeriod)

	if traffic.ShouldCalculateWithoutTraffic(trafficRecords) {
		return model.calculateWithoutTraffic(forecastRule)
	}

	if forecastRule.Period.Contains(*forecastRule.LHM) {
		adjustedForecastVolume, _err := model.extractHistoricalVolumeFromForecasted(forecastRule)

		if _err != nil {
			return nil, _err
		}

		forecastRule.Volume = adjustedForecastVolume
	}

	forecastRecords := model.calculateWithTraffic(forecastRule, trafficRecords)

	return forecastRecords, nil
}

func (model *manualVolume) calculateWithoutTraffic(forecastRule *rules.ForecastRule) (
	[]calculation.ForecastRecord,
	error,
) {
	validatedPeriod, err := forecastRule.GetValidatedPeriod()

	if err != nil {
		return nil, err
	}

	months := validatedPeriod.GetMonths()

	totalMonths := len(months)

	volumePerRecord := forecastRule.Volume / float64(len(months))

	forecastRecords := make([]calculation.ForecastRecord, totalMonths)

	for idx, month := range months {
		record := calculation.ForecastRecord{VolumeActual: volumePerRecord, Month: month}

		forecastRecords[idx] = record
	}

	return forecastRecords, nil
}

func (model *manualVolume) extractHistoricalVolumeFromForecasted(forecastRule *rules.ForecastRule) (float64, error) {
	historicalPeriodInForecasted := types.NewPeriod(forecastRule.Period.StartDate, *forecastRule.LHM)

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

	validatedPeriod, _ := forecastRule.GetValidatedPeriod()

	totalForecastedMonths := validatedPeriod.GetTotalMonths()

	historicalTrafficMonthMap := make(map[string]traffic.MonthlyAggregationRecord, totalForecastedMonths)

	for _, record := range trafficRecords {
		historicalTrafficMonthMap[record.Month.ToDateString()] = record
	}

	forecastRecords := make([]calculation.ForecastRecord, totalForecastedMonths)

	for idx, forecastedMonth := range validatedPeriod.GetMonths() {
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

func createTrafficPeriodFromForecasted(rule *rules.ForecastRule) (types.Period, error) {
	forecastPeriod, err := rule.GetValidatedPeriod()

	if err != nil {
		return types.Period{}, nil
	}

	budgetTrafficPeriod := types.Period{
		StartDate: forecastPeriod.StartDate.SubYear().ToDateStruct(),
		EndDate:   forecastPeriod.EndDate.SubYear().ToDateStruct(),
	}

	if forecastPeriod.GetTotalMonths() > 12 {
		budgetTrafficPeriod = types.Period{
			StartDate: budgetTrafficPeriod.StartDate,
			EndDate:   budgetTrafficPeriod.StartDate.AddMonths(12 - 1).ToDateStruct(),
		}
	}

	return budgetTrafficPeriod, nil
}
