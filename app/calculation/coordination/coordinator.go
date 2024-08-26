package coordination

import (
	"log"

	"forecasting/app/domain/repositories"
	"forecasting/app/providers"
	"github.com/golang-module/carbon/v2"
)

type ForecastRuleCalculationCoordinator struct {
	forecastRulesRepository repositories.ForecastRuleRepository

	forecastingService forecastingService

	budgetTrafficProvider providers.BudgetTrafficProvider
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
	}
}

func (c ForecastRuleCalculationCoordinator) CalculateAll() {
	log.Println("Start traffic clearing")

	c.budgetTrafficProvider.ClearForecasted()

	log.Println("Traffic cleared")

	log.Println("Start forecast rules retrieving")

	forecastRules := c.forecastRulesRepository.GetMany()

	log.Println("Forecast rules are retrieved - ", len(forecastRules))

	log.Println("Forecast rules application started")

	ruleCounter := 0

	for _, forecastRule := range forecastRules {
		forecastRule.LHM = carbon.Parse("2024-02-01").ToDateStruct()

		budgetTrafficRecords := c.forecastingService.Evaluate(forecastRule)

		c.budgetTrafficProvider.CreateMany(budgetTrafficRecords)

		if ruleCounter == 5000 {
			log.Println("5000 rules calculated")
			ruleCounter = 0
		}

		ruleCounter++
	}

	log.Println("Forecast rules application finished")

	c.budgetTrafficProvider.CountForecasted()

	log.Println("Total not calculated rules: ", c.forecastingService.TotalNotCalculatedRules)

	log.Println("Finished")
}
