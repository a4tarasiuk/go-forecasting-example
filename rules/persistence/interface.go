package persistence

import "forecasting/rules"

type ForecastRuleRepository interface {
	GetMany() []rules.ForecastRule
}
