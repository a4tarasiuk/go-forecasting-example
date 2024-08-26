package models

import "github.com/golang-module/carbon/v2"

type MonthlyAggregationRecord struct {
	Month carbon.Date

	VolumeActual float64
}

func (r MonthlyAggregationRecord) GetVolumeActual() float64 { return r.VolumeActual }
