package forecast_model

import (
	"forecasting/calculation"
	"forecasting/rules"
)

type ForecastModel interface {
	Calculate(forecastRule rules.ForecastRule) []calculation.ForecastRecord
}
