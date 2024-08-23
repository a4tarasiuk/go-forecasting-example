package distribution_models

import (
	"forecasting/calculation"
	"forecasting/core"
	"forecasting/core/types"
	"forecasting/rules"
	"forecasting/traffic"
)

type movingAverage struct {
	budgetTrafficProvider traffic.BudgetTrafficProvider
}

func NewMovingAverage() movingAverage {
	return movingAverage{budgetTrafficProvider: traffic.NewBudgetTrafficProvider()}
}

func (ma *movingAverage) Apply(
	forecastRule *rules.ForecastRule,
	forecastRecords []calculation.ForecastRecord,
) (
	[]calculation.DistributionRecord,
	error,
) {
	// TODO:
	//  3. Distribute forecast records based on historical traffic
	//    3.1. Aggregate historical traffic records by key combination
	//	  3.2. Calculate total volume actual (and for all other columns)
	//	  3.3. Map every record with its forecast record (implementation detail from original solution)
	//    3.4. Calculate shares within each month for every combination

	historicalTrafficRecords := ma.loadHistoricalTraffic(forecastRule)

	if traffic.ShouldCalculateWithoutTraffic(historicalTrafficRecords) {
		return ma.calculateWithoutTraffic(forecastRule, forecastRecords), nil
	}

	return nil, nil
}

func (ma *movingAverage) loadHistoricalTraffic(forecastRule *rules.ForecastRule) []traffic.BudgetTrafficRecord {
	nMonthPeriodEndDate := forecastRule.LHM.SubMonth().ToDateStruct()

	monthsToSub := *forecastRule.DistributionModelMovingAverageMonths - 1

	nMonthPeriodStartDate := nMonthPeriodEndDate.SubMonths(monthsToSub).ToDateStruct()

	nMonthPeriod := types.NewPeriod(nMonthPeriodStartDate, nMonthPeriodEndDate)

	options := traffic.BudgetTrafficOptions{
		ForecastRule:   forecastRule,
		Period:         &nMonthPeriod,
		HistoricalOnly: true,
	}

	budgetTrafficRecords := ma.budgetTrafficProvider.Get(options)

	return budgetTrafficRecords
}

func (ma *movingAverage) calculateWithoutTraffic(
	forecastRule *rules.ForecastRule,
	forecastRecords []calculation.ForecastRecord,
) []calculation.DistributionRecord {

	totalHomeOperators := int64(len(forecastRule.HomeOperators))
	totalPartnerOperators := int64(len(forecastRule.PartnerOperators))
	totalMonths := forecastRule.Period.GetTotalMonths()

	totalResultRecords := totalHomeOperators * totalPartnerOperators * totalMonths

	distributionRecords := make([]calculation.DistributionRecord, totalResultRecords)

	idx := 0

	for _, forecastRecord := range forecastRecords {
		totalRecordsWithinMonth := totalHomeOperators * totalPartnerOperators

		for _, homeOperatorID := range forecastRule.HomeOperators {
			for _, partnerOperatorID := range forecastRule.PartnerOperators {

				_distributionRecord := calculation.DistributionRecord{
					HomeOperatorID:    homeOperatorID,
					PartnerOperatorID: partnerOperatorID,
					Month:             forecastRecord.Month,
					CallDestination:   core.GetDefaultCDByServiceType(forecastRule.ServiceType),
					CalledCountryID:   nil,
					IsPremium:         nil,
					TrafficSegmentID:  nil,
					IMSICountType:     nil,
					VolumeActual:      forecastRecord.VolumeActual / float64(totalRecordsWithinMonth),
				}

				distributionRecords[idx] = _distributionRecord

				idx++
			}
		}
	}

	return distributionRecords
}
