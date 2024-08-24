package traffic

import (
	"forecasting/core"
	"github.com/golang-module/carbon/v2"
)

type BudgetTrafficRecord struct {
	BudgetSnapshotID int64

	HomeOperatorID    int64
	PartnerOperatorID int64

	TrafficType      byte // Can be enum
	TrafficDirection core.TrafficDirection

	Month carbon.Date

	ServiceType core.ServiceType

	CallDestination *byte
	CalledCountryID *int64
	IsPremium       *bool

	IMSICountType *byte

	TrafficSegmentID *int64

	VolumeActual float64
}

func (r BudgetTrafficRecord) GetVolumeActual() float64 { return r.VolumeActual }

func (r BudgetTrafficRecord) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"BudgetSnapshotID":  r.BudgetSnapshotID,
		"HomeOperatorID":    r.HomeOperatorID,
		"PartnerOperatorID": r.PartnerOperatorID,
		"TrafficType":       r.TrafficType,
		"TrafficDirection":  r.TrafficDirection,
		"ServiceType":       r.ServiceType,
		"CallDestination":   serializeNullableNumber(r.CallDestination),
		"CalledCountryID":   serializeNullableNumber(r.CalledCountryID),
		"IsPremium":         serializeNullableBool(r.IsPremium),
		"IMSICountType":     serializeNullableNumber(r.IMSICountType),
		"TrafficSegmentID":  serializeNullableNumber(r.TrafficSegmentID),
		"Month":             r.Month.ToDateString(),
		"VolumeActual":      r.VolumeActual,
	}
}

func serializeNullableBool(val *bool) int {
	sVal := -1

	if val != nil && *val == true {
		one := 1
		sVal = one
	} else if val != nil && *val == false {
		zero := 0
		sVal = zero
	}

	return sVal
}

func serializeNullableNumber[V int | int64 | byte](val *V) int {
	if val == nil {
		return -1
	}

	return int(*val)
}
