package forecasting

import (
	"log"

	"forecasting/calculation"
	"forecasting/calculation/distribution_models"
	"forecasting/calculation/forecast_models"
	"forecasting/rules"
	"forecasting/traffic"
)

type Service interface {
	Evaluate(forecastRule *rules.ForecastRule) []calculation.DistributionRecord
}

type forecastingService struct {
	mapper BudgetTrafficRecordMapper
}

func NewService(mapper BudgetTrafficRecordMapper) forecastingService {
	return forecastingService{mapper: mapper}
}

func (s *forecastingService) Evaluate(forecastRule *rules.ForecastRule) []traffic.BudgetTrafficRecord {

	forecastModel := forecast_models.NewManualVolume()

	forecastRecords, err := forecastModel.Calculate(forecastRule)

	if err != nil {
		log.Print("rule skipped after forecast", forecastRule)
		return nil
	}

	distributionModel := distribution_models.NewMovingAverage()

	distributionRecords, _err := distributionModel.Apply(forecastRule, forecastRecords)

	if _err != nil {
		log.Print("rule skipped after distribution", forecastRule)
		return nil
	}

	budgetTrafficRecords := make([]traffic.BudgetTrafficRecord, len(distributionRecords))

	for idx, _distributionRecord := range distributionRecords {
		budgetTrafficRecords[idx] = s.mapper.FromDistributionToBudgetTrafficRecord(forecastRule, _distributionRecord)
	}

	return budgetTrafficRecords
}
