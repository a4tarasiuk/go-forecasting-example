package rules

import (
	"forecasting/core"
	"forecasting/core/types"
	"github.com/golang-module/carbon/v2"
)

type ForecastRule struct {
	ID int

	HomeOperators    []int
	PartnerOperators []int

	Period types.Period

	TrafficDirection core.TrafficDirection

	ServiceType core.ServiceType

	ForecastModel     core.ForecastModel
	DistributionModel core.DistributionModel

	Volume float64

	DistributionModelMovingAverageMonths *int

	LHM *carbon.Date
}
