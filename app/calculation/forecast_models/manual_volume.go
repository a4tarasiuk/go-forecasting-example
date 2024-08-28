package forecast_models

import (
	"errors"

	"forecasting/app/calculation/calc_utils"
	"forecasting/app/calculation/dto"
	"forecasting/app/domain/models"
	"forecasting/app/providers"
	"forecasting/core/types"
)

type manualVolume struct {
	trafficProvider providers.MonthlyAggregationProvider
}

func NewManualVolume(trafficProvider providers.MonthlyAggregationProvider) manualVolume {
	return manualVolume{trafficProvider: trafficProvider}
}

func (model *manualVolume) Calculate(forecastRule *models.ForecastRule) ([]dto.ForecastRecord, error) {
	trafficPeriod, err := createTrafficPeriodFromForecasted(forecastRule)

	if err != nil {
		return nil, err
	}

	trafficRecords := model.trafficProvider.GetLast(forecastRule, trafficPeriod)

	if calc_utils.ShouldCalculateWithoutTraffic(trafficRecords) {
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

func (model *manualVolume) calculateWithoutTraffic(forecastRule *models.ForecastRule) []dto.ForecastRecord {
	validatedPeriod, _ := forecastRule.GetValidatedPeriod()

	months := validatedPeriod.GetMonths()

	totalMonths := len(months)

	volumePerRecord := forecastRule.Volume / float64(totalMonths)

	forecastRecords := make([]dto.ForecastRecord, totalMonths)

	for idx, month := range months {
		record := dto.ForecastRecord{VolumeActual: volumePerRecord, Month: month}

		forecastRecords[idx] = record
	}

	return forecastRecords
}

func (model *manualVolume) extractHistoricalVolumeFromForecasted(forecastRule *models.ForecastRule) (float64, error) {
	historicalPeriodInForecasted := types.NewPeriod(forecastRule.Period.StartDate, forecastRule.LHM)

	trafficRecords := model.trafficProvider.Get(forecastRule, historicalPeriodInForecasted)

	totalHistoricalVolume := calc_utils.CalculateTotalHistoricalVolume(trafficRecords)

	if totalHistoricalVolume > forecastRule.Volume {
		return 0, errors.New("historical volume exceeds forecasted volume")
	}

	adjustedForecastVolume := forecastRule.Volume - totalHistoricalVolume

	return adjustedForecastVolume, nil
}

func (model *manualVolume) calculateWithTraffic(
	forecastRule *models.ForecastRule,
	trafficRecords []models.MonthlyAggregationRecord,
) []dto.ForecastRecord {

	totalHistoricalVolume := calc_utils.CalculateTotalHistoricalVolume(trafficRecords)

	totalForecastedVolume := forecastRule.Volume

	historicalTrafficMonthMap := make(map[string]models.MonthlyAggregationRecord, len(trafficRecords))

	for _, record := range trafficRecords {
		historicalTrafficMonthMap[record.Month.ToDateString()] = record
	}

	forecastPeriod, _ := forecastRule.GetValidatedPeriod()

	totalForecastedMonths := forecastPeriod.GetTotalMonths()

	forecastRecords := make([]dto.ForecastRecord, totalForecastedMonths)

	for idx, forecastMonth := range forecastPeriod.GetMonths() {
		historicalRecord := historicalTrafficMonthMap[forecastMonth.SubYear().ToDateString()]

		forecastedVolumeActual := (historicalRecord.VolumeActual / totalHistoricalVolume) * totalForecastedVolume

		forecastRecord := dto.ForecastRecord{
			VolumeActual: forecastedVolumeActual,
			Month:        forecastMonth,
		}

		forecastRecords[idx] = forecastRecord
	}

	return forecastRecords
}

func createTrafficPeriodFromForecasted(rule *models.ForecastRule) (types.Period, error) {
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
