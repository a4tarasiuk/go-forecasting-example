package repositories

import "forecasting/app/domain/models"

type ForecastRuleRepository interface {
	GetMany() []*models.ForecastRule
}
