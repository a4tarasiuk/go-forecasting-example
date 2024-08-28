package distribution_models

import (
	"forecasting/app/calculation/dto"
	"forecasting/app/domain/models"
)

type DistributionModel interface {
	Apply(
		forecastRule *models.ForecastRule,
		forecastRecords []dto.ForecastRecord,
	) ([]dto.DistributionRecord, error)
}
