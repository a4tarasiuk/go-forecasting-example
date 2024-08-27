package coordination

import (
	"forecasting/app/domain/models"
)

func ForecastWorker(
	service forecastingService,
	forecastRules <-chan *models.ForecastRule,
	budgetTrafficRecords chan<- []models.BudgetTrafficRecord,
) {
	for forecastRule := range forecastRules {
		records := service.Evaluate(forecastRule)

		budgetTrafficRecords <- records
	}
}
