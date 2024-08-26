package calculation

import "forecasting/app/domain/models"

func CalculateTotalHistoricalVolume[R models.VolumedRecord](trafficRecords []R) float64 {
	total := 0.0

	for _, record := range trafficRecords {
		total += record.GetVolumeActual()
	}

	return total
}

func ShouldCalculateWithoutTraffic[R models.VolumedRecord](trafficRecords []R) bool {
	trafficIsEmpty := len(trafficRecords) == 0

	trafficIsZeroVolume := CalculateTotalHistoricalVolume(trafficRecords) == 0

	return trafficIsEmpty || trafficIsZeroVolume
}
