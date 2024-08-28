package calculation

import (
	"log"

	"forecasting/app/domain/models"
	"forecasting/app/domain/repositories"
	"forecasting/app/providers"
)

type ForecastRuleCalculationCoordinator struct {
	forecastRulesRepository repositories.ForecastRuleRepository

	forecastingService forecastingService

	budgetTrafficProvider providers.BudgetTrafficProvider

	budgetTrafficFactory *BudgetTrafficFactory
}

func NewForecastRuleCalculationCoordinator(
	forecastRulesRepository repositories.ForecastRuleRepository,
	forecastingService forecastingService,
	budgetTrafficProvider providers.BudgetTrafficProvider,
) ForecastRuleCalculationCoordinator {

	return ForecastRuleCalculationCoordinator{
		forecastRulesRepository: forecastRulesRepository,
		forecastingService:      forecastingService,
		budgetTrafficProvider:   budgetTrafficProvider,
		budgetTrafficFactory:    NewBudgetTrafficFactory(budgetTrafficProvider),
	}
}

func (c *ForecastRuleCalculationCoordinator) CalculateAll() {
	c.budgetTrafficProvider.ClearForecasted()

	forecastRules := c.loadForecastRules()

	log.Println("Forecast rules application started")

	totalForecastRules := len(forecastRules)
	totalWorkers := 10

	forecastRulesChannel := make(chan *models.ForecastRule, totalForecastRules)
	btrChannel := make(chan []models.BudgetTrafficRecord)

	for range totalWorkers {
		go ForecastWorker(c.forecastingService, forecastRulesChannel, btrChannel)
	}

	for _, forecastRule := range forecastRules {
		forecastRulesChannel <- forecastRule
	}

	close(forecastRulesChannel)

	totalBudgetTrafficRecords := c.handleTraffic(btrChannel, totalForecastRules)

	// logs
	log.Println("Forecast rules application finished")

	c.budgetTrafficProvider.CountForecasted()

	log.Println("Total not calculated rules: ", c.forecastingService.TotalNotCalculatedRules)

	log.Println("Total created records: ", totalBudgetTrafficRecords)

	log.Println("Finished")
}

func (c *ForecastRuleCalculationCoordinator) loadForecastRules() []*models.ForecastRule {
	forecastRules := c.forecastRulesRepository.GetMany()

	log.Println("Forecast rules are retrieved - ", len(forecastRules))

	return forecastRules
}

func (c *ForecastRuleCalculationCoordinator) handleTraffic(
	records <-chan []models.BudgetTrafficRecord,
	totalForecastRules int,
) int {
	totalBudgetTrafficRecords := 0
	totalProcessedForecastRules := 0

	for range totalForecastRules {
		budgetTrafficRecords := <-records

		c.budgetTrafficProvider.CreateMany(budgetTrafficRecords)

		if totalProcessedForecastRules == 5000 {
			log.Println("5000 rules processed")
			totalProcessedForecastRules = 0
		}

		totalBudgetTrafficRecords += len(budgetTrafficRecords)

		totalProcessedForecastRules++
	}

	return totalBudgetTrafficRecords
}
