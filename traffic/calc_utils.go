package traffic

func CalculateTotalHistoricalVolume(trafficRecords []MonthlyAggregationRecord) float64 {
	total := 0.0

	for _, record := range trafficRecords {
		total += record.VolumeActual
	}

	return total
}

func ShouldCalculateWithoutTraffic(trafficRecords []MonthlyAggregationRecord) bool {
	trafficIsEmpty := len(trafficRecords) == 0

	trafficIsZeroVolume := CalculateTotalHistoricalVolume(trafficRecords) == 0

	return trafficIsEmpty || trafficIsZeroVolume
}
