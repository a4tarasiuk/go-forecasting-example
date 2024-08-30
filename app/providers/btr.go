package providers

import (
	"forecasting/app/domain/models"
	"forecasting/core/types"
)

type (
	BudgetTrafficOptions struct {
		ForecastRule   *models.ForecastRule
		Period         *types.Period
		HistoricalOnly bool
	}

	BudgetTrafficProvider interface {
		Get(options BudgetTrafficOptions) []models.BudgetTrafficRecord

		CreateMany(records []models.BudgetTrafficRecord)

		ClearForecasted()

		CountForecasted() int64

		Count() int64
	}
)
