package calculation

import "github.com/golang-module/carbon/v2"

type DistributionRecord struct {
	HomeOperatorID    int
	PartnerOperatorID int

	Month carbon.Date

	CallDestination *int
	CalledCountryID *int
	IsPremium       *bool

	TrafficSegmentID *int
	IMSICountType    *int

	VolumeActual float64
}

func (r DistributionRecord) GetVolumeActual() float64 {
	return r.VolumeActual
}
