package coordination

import (
	"forecasting/app/calculation"
	"forecasting/app/calculation/distribution_models"
	"forecasting/app/calculation/forecast_models"
	"forecasting/app/domain/models"
	"forecasting/app/providers"
)

type Service interface {
	Evaluate(forecastRule *models.ForecastRule) []calculation.DistributionRecord
}

type forecastingService struct {
	mapper BudgetTrafficRecordMapper

	maProvider providers.MonthlyAggregationProvider

	btrProvider providers.BudgetTrafficProvider

	TotalNotCalculatedRules int
}

func NewForecastingService(
	btrMapper BudgetTrafficRecordMapper,
	maProvider providers.MonthlyAggregationProvider,
	btrProvider providers.BudgetTrafficProvider,
) forecastingService {

	return forecastingService{
		mapper:                  btrMapper,
		maProvider:              maProvider,
		btrProvider:             btrProvider,
		TotalNotCalculatedRules: 0,
	}
}

func (s *forecastingService) Evaluate(forecastRule *models.ForecastRule) []models.BudgetTrafficRecord {

	forecastModel := forecast_models.NewManualVolume(s.maProvider)

	forecastRecords, err := forecastModel.Calculate(forecastRule)

	if err != nil {
		s.TotalNotCalculatedRules++
		return nil
	}

	distributionModel := distribution_models.NewMovingAverage(s.btrProvider)

	distributionRecords, _err := distributionModel.Apply(forecastRule, forecastRecords)

	if _err != nil {
		s.TotalNotCalculatedRules++
		return nil
	}

	budgetTrafficRecords := make([]models.BudgetTrafficRecord, len(distributionRecords))

	for idx, _distributionRecord := range distributionRecords {
		budgetTrafficRecords[idx] = s.mapper.FromDistributionToBudgetTrafficRecord(forecastRule, _distributionRecord)
	}

	return budgetTrafficRecords
}
