package coordination

import (
	"forecasting/app/calculation"
	"forecasting/app/calculation/distribution_models"
	"forecasting/app/calculation/forecast_models"
	"forecasting/app/domain/models"
)

type Service interface {
	Evaluate(forecastRule *models.ForecastRule) []calculation.DistributionRecord
}

type forecastingService struct {
	mapper BudgetTrafficRecordMapper

	TotalNotCalculatedRules int
}

func NewService(mapper BudgetTrafficRecordMapper) forecastingService {
	return forecastingService{mapper: mapper, TotalNotCalculatedRules: 0}
}

func (s *forecastingService) Evaluate(forecastRule *models.ForecastRule) []models.BudgetTrafficRecord {

	forecastModel := forecast_models.NewManualVolume()

	forecastRecords, err := forecastModel.Calculate(forecastRule)

	if err != nil {
		s.TotalNotCalculatedRules++
		return nil
	}

	distributionModel := distribution_models.NewMovingAverage()

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
