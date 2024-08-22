package distribution_models

import (
	"forecasting/calculation"
	"forecasting/rules"
)

type movingAverage struct {
}

func NewMovingAverage() movingAverage {
	return movingAverage{}
}

func (ma *movingAverage) Apply(
	forecastRule *rules.ForecastRule,
	forecastRecords []calculation.ForecastRecord,
) (
	[]calculation.DistributionResultRecord,
	error,
) {
	// TODO:
	//  1. Get historical traffic for calculation/distribution
	//	  1.1. Evaluate n-months period
	//    1.2. Get traffic for period from 1.1 (HISTORICAL only)
	//  2. Distribute forecast records without historical traffic if it is empty or zero-volume
	//		- Even distribution
	//  3. Distribute forecast records based on historical traffic
	//    3.1. Aggregate historical traffic records by key combination
	//	  3.2. Calculate total volume actual (and for all other columns)
	//	  3.3. Map every record with its forecast record (implementation detail from original solution)
	//    3.4. Calculate shares within each month for every combination

	// historicalTrafficRecords := ma.loadHistoricalTraffic(forecastRule)

	return nil, nil
}

func (ma *movingAverage) loadHistoricalTraffic(forecastRule *rules.ForecastRule) []interface{} {
	return nil
}
