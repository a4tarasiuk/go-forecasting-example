package forecast_models

import (
	"forecasting/calculation"
	"forecasting/rules"
)

type ForecastModel interface {
	Calculate(forecastRule *rules.ForecastRule) ([]calculation.ForecastRecord, error)
}
