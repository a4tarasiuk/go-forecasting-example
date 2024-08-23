package traffic

import "forecasting/core"

type BudgetTrafficRecord struct {
	MonthlyAggregationRecord

	ID               int
	BudgetSnapshotID int

	HomeOperatorID    int
	PartnerOperatorID int

	TrafficType      byte // Can be enum
	TrafficDirection core.TrafficDirection

	ServiceType core.ServiceType

	CallDestination *int
	CalledCountryID *int
	IsPremium       *bool

	IMSICountType int
}

func (r BudgetTrafficRecord) GetVolumeActual() float64 { return r.VolumeActual }
