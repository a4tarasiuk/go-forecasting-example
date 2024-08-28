package calculation

import (
	"forecasting/app/calculation/distribution_models"
	"forecasting/app/calculation/forecast_models"
	"forecasting/app/domain/models"
	"forecasting/app/providers"
)

type ForecastingService struct {
	mapper BudgetTrafficRecordMapper

	maProvider providers.MonthlyAggregationProvider

	btrProvider providers.BudgetTrafficProvider
}

func NewForecastingService(
	maProvider providers.MonthlyAggregationProvider,
	btrProvider providers.BudgetTrafficProvider,
) ForecastingService {

	return ForecastingService{
		mapper:      NewBudgetTrafficRecordMapper(),
		maProvider:  maProvider,
		btrProvider: btrProvider,
	}
}

func (s *ForecastingService) Evaluate(forecastRule *models.ForecastRule) []models.BudgetTrafficRecord {

	forecastModel := forecast_models.NewManualVolume(s.maProvider)

	forecastRecords, err := forecastModel.Calculate(forecastRule)

	if err != nil {
		return nil
	}

	distributionModel := distribution_models.NewMovingAverage(s.btrProvider)

	distributionRecords, _err := distributionModel.Apply(forecastRule, forecastRecords)

	if _err != nil {
		return nil
	}

	budgetTrafficRecords := make([]models.BudgetTrafficRecord, len(distributionRecords))

	for idx, _distributionRecord := range distributionRecords {
		budgetTrafficRecords[idx] = s.mapper.FromDistributionToBudgetTrafficRecord(forecastRule, _distributionRecord)
	}

	return budgetTrafficRecords
}
