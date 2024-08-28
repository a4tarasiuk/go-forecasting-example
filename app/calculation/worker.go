package calculation

import (
	"forecasting/app/domain/models"
)

func ForecastWorker(
	service ForecastingService,
	forecastRules <-chan *models.ForecastRule,
	budgetTrafficRecords chan<- []models.BudgetTrafficRecord,
) {
	for forecastRule := range forecastRules {
		records := service.Evaluate(forecastRule)

		budgetTrafficRecords <- records
	}
}
