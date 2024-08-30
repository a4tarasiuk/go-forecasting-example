package persistence

import "forecasting/app/domain/models"

type ForecastRuleInMemoryRepository struct {
	rules []*models.ForecastRule
}

func NewForecastRuleInMemoryRepository(rules []*models.ForecastRule) *ForecastRuleInMemoryRepository {
	return &ForecastRuleInMemoryRepository{rules: rules}
}

func (r *ForecastRuleInMemoryRepository) GetMany() []*models.ForecastRule {
	return r.rules
}
