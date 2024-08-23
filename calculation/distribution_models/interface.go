package distribution_models

import (
	"forecasting/calculation"
	"forecasting/rules"
)

type DistributionModel interface {
	Apply(
		forecastRule *rules.ForecastRule,
		forecastRecords []calculation.ForecastRecord,
	) ([]calculation.DistributionRecord, error)
}
