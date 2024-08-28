package dto

import "github.com/golang-module/carbon/v2"

type ForecastRecord struct {
	VolumeActual float64

	Month carbon.Date
}

func (r ForecastRecord) GetVolumeActual() float64 {
	return r.VolumeActual
}
