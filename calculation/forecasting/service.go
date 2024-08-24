package forecasting

import (
	"log"

	"forecasting/calculation"
	"forecasting/calculation/distribution_models"
	"forecasting/calculation/forecast_models"
	"forecasting/rules"
)

type Service interface {
	Evaluate(forecastRule *rules.ForecastRule) []calculation.DistributionRecord
}

type forecastingService struct {
}

func NewService() forecastingService {
	return forecastingService{}
}

func (s *forecastingService) Evaluate(forecastRule *rules.ForecastRule) {

	forecastModel := forecast_models.NewManualVolume()

	forecastRecords, err := forecastModel.Calculate(forecastRule)

	if err != nil {
		log.Print("rule skipped after forecast", forecastRule)
		return
	}

	distributionModel := distribution_models.NewMovingAverage()

	_, _err := distributionModel.Apply(forecastRule, forecastRecords)

	if _err != nil {
		log.Print("rule skipped after distribution", forecastRule)
		return
	}

	// for _, r := range distributionRecords {
	// 	fmt.Printf("%+v\n", r)
	// }

	// TODO: Map to budget traffic records or use interface to that table to write them directly to DB
}
