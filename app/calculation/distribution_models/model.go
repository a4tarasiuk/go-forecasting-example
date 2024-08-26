package distribution_models

import (
	"forecasting/app/calculation"
	"forecasting/app/domain/models"
)

type DistributionModel interface {
	Apply(
		forecastRule *models.ForecastRule,
		forecastRecords []calculation.ForecastRecord,
	) ([]calculation.DistributionRecord, error)
}
