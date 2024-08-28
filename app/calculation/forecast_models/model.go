package forecast_models

import (
	"forecasting/app/calculation/dto"
	"forecasting/app/domain/models"
)

type ForecastModel interface {
	Calculate(forecastRule *models.ForecastRule) ([]dto.ForecastRecord, error)
}
