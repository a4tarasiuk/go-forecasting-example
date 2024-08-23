package traffic

func CalculateTotalHistoricalVolume[R volumedRecord](trafficRecords []R) float64 {
	total := 0.0

	for _, record := range trafficRecords {
		total += record.GetVolumeActual()
	}

	return total
}

func ShouldCalculateWithoutTraffic[R volumedRecord](trafficRecords []R) bool {
	trafficIsEmpty := len(trafficRecords) == 0

	trafficIsZeroVolume := CalculateTotalHistoricalVolume(trafficRecords) == 0

	return trafficIsEmpty || trafficIsZeroVolume
}
