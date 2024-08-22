package traffic

import "github.com/golang-module/carbon/v2"

type MonthlyAggregationRecord struct {
	Month carbon.Date

	VolumeActual float64
}
