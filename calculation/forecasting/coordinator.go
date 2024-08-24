package forecasting

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
	log.Println("Start forecast rules retrieving")

	forecastRules := c.forecastRulesRepository.GetMany()

	log.Println("Forecast rules are retrieved")

	var lhm *carbon.Date

	if budget_defaults.BudgetLHM != nil {
		dt := carbon.Parse(*budget_defaults.BudgetLHM).ToDateStruct()
		lhm = &dt
	}

	log.Println("Forecast rules application started")

	for _, forecastRule := range forecastRules {
		forecastRule.LHM = lhm

		budgetTrafficRecords := c.forecastingService.Evaluate(forecastRule)

		c.budgetTrafficProvider.CreateMany(budgetTrafficRecords)
	}

	log.Println("Forecast rules application finished")
}
