package coordination

import (
	"log"

	"forecasting/budget_defaults"
	"forecasting/postgres"
	"forecasting/rules/persistence"
	"forecasting/traffic"
	persistence2 "forecasting/traffic/persistence"
	"github.com/golang-module/carbon/v2"
)

type ForecastRuleCalculationCoordinator struct {
	forecastRulesRepository persistence.ForecastRuleRepository

	forecastingService forecastingService

	budgetTrafficProvider traffic.BudgetTrafficProvider
}

func NewForecastRuleCalculationCoordinator() ForecastRuleCalculationCoordinator {
	return ForecastRuleCalculationCoordinator{
		forecastRulesRepository: persistence.NewPostgresForecastRuleRepository(),
		forecastingService:      NewService(NewBudgetTrafficRecordMapper(postgres.DB, budget_defaults.BudgetID)),
		budgetTrafficProvider:   persistence2.NewPostgresBudgetTrafficProvider(),
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

	log.Println("Finished")
}
