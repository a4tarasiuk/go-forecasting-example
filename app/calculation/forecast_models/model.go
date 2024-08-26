package forecast_models

import (
	"forecasting/app/calculation"
	"forecasting/app/domain/models"
)

type ForecastModel interface {
	Calculate(forecastRule *models.ForecastRule) ([]calculation.ForecastRecord, error)
}
