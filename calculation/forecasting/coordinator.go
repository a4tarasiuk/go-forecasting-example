package forecasting

import (
	"log"

	"forecasting/budget_defaults"
	"forecasting/rules/persistence"
	"github.com/golang-module/carbon/v2"
)

type ForecastRuleCalculationCoordinator struct {
	forecastRulesRepository persistence.ForecastRuleRepository

	forecastingService forecastingService
}

func NewForecastRuleCalculationCoordinator() ForecastRuleCalculationCoordinator {
	return ForecastRuleCalculationCoordinator{
		forecastRulesRepository: persistence.NewPostgresForecastRuleRepository(),
		forecastingService:      NewService(),
	}
}

func (c ForecastRuleCalculationCoordinator) CalculateAll() {
	forecastRules := c.forecastRulesRepository.GetMany()

	var lhm *carbon.Date

	if budget_defaults.BudgetLHM != nil {
		dt := carbon.Parse(*budget_defaults.BudgetLHM).ToDateStruct()
		lhm = &dt
	}

	for _, forecastRule := range forecastRules {
		log.Print("start apply rule ", forecastRule.ID)

		forecastRule.LHM = lhm

		c.forecastingService.Evaluate(forecastRule)
	}
}
