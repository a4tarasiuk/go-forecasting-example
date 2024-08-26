package forecast_models

import (
	"forecasting/app/calculation"
	"forecasting/rules/models"
)

type ForecastModel interface {
	Calculate(forecastRule *models.ForecastRule) ([]calculation.ForecastRecord, error)
}
