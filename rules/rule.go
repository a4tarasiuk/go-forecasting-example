package rules

import (
	"forecasting/core"
	"forecasting/core/types"
	"github.com/golang-module/carbon/v2"
)

type ForecastRule struct {
	ID int64

	BudgetID int64

	HomeOperators    []int64
	PartnerOperators []int64

	Period types.Period

	TrafficDirection core.TrafficDirection

	ServiceType core.ServiceType

	ForecastModel     core.ForecastModel
	DistributionModel core.DistributionModel

	Volume float64

	DistributionModelMovingAverageMonths *int

	LHM *carbon.Date
}
