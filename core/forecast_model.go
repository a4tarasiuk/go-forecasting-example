package core

type ForecastModel byte

const (
	ManualVolumeForecastModel ForecastModel = iota
	MovingAverageForecastModel
	PPIForecastModel
	RollingAverageForecastModel
	RollingAverageNoSeasonForecastModel
)
