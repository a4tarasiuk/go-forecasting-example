package distribution_models

import (
	"forecasting/calculation"
	"forecasting/core"
	"forecasting/core/types"
	"forecasting/rules"
	"forecasting/traffic"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
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
	if forecastRule.LHM == nil {
		return ma.calculateWithoutTraffic(forecastRule, forecastRecords), nil
	}

	historicalTrafficRecords := ma.loadHistoricalTraffic(forecastRule)

	if traffic.ShouldCalculateWithoutTraffic(historicalTrafficRecords) {
		return ma.calculateWithoutTraffic(forecastRule, forecastRecords), nil
	}

	return ma.calculateWithTraffic(forecastRule, forecastRecords, historicalTrafficRecords), nil
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

func (ma *movingAverage) calculateWithTraffic(
	forecastRule *rules.ForecastRule,
	forecastRecords []calculation.ForecastRecord,
	historicalRecords []traffic.BudgetTrafficRecord,
) []calculation.DistributionRecord {

	rawHistoricalTrafficRecords := make([]map[string]interface{}, len(historicalRecords))

	for idx, _historicalRecord := range historicalRecords {
		rawHistoricalTrafficRecords[idx] = _historicalRecord.Serialize()
	}

	typesMap := map[string]series.Type{
		"TrafficType":      series.Int,
		"TrafficDirection": series.Int,
		"ServiceType":      series.Int,
		"CalledCountryID":  series.Int,
	}

	historicalTrafficDF := dataframe.LoadMaps(rawHistoricalTrafficRecords, dataframe.WithTypes(typesMap))

	// Aggregation
	historicalTrafficDF = historicalTrafficDF.GroupBy(
		"HomeOperatorID",
		"PartnerOperatorID",
		"TrafficDirection",
		"ServiceType",
		"CallDestination",
		"CalledCountryID",
		"IsPremium",
		"TrafficSegmentID",
		"IMSICountType",
	).Aggregation(
		[]dataframe.AggregationType{dataframe.Aggregation_SUM}, []string{"VolumeActual"},
	)

	// Coefficient calculation
	totalHistoricalVolume := traffic.CalculateTotalHistoricalVolume(historicalRecords)

	divideOnTotal := func(el series.Element) series.Element {
		el.Set(el.Val().(float64) / totalHistoricalVolume)
		return el
	}

	historicalTrafficDF = historicalTrafficDF.Mutate(
		series.New(
			historicalTrafficDF.Col("VolumeActual_SUM").Map(divideOnTotal),
			series.Float,
			"ForecastCoefficient",
		),
	)

	// Duplicate records for every forecasted month
	forecastDfMap := make(map[calculation.ForecastRecord]dataframe.DataFrame)

	for _, forecastRecord := range forecastRecords {
		forecastDfMap[forecastRecord] = historicalTrafficDF.Copy()
	}

	totalResultRecords := int64(historicalTrafficDF.Nrow()) * forecastRule.Period.GetTotalMonths()

	distributionRecords := make([]calculation.DistributionRecord, totalResultRecords)

	idx := 0

	for forecastRecord, df := range forecastDfMap {

		for _, row := range df.Maps() {
			forecastedVolumeActual := row["ForecastCoefficient"].(float64) * forecastRecord.VolumeActual

			_distributionRecord := calculation.DistributionRecord{
				HomeOperatorID:    row["HomeOperatorID"].(int),
				PartnerOperatorID: row["PartnerOperatorID"].(int),
				Month:             forecastRecord.Month,
				CallDestination:   mapOptionalInt("CallDestination", row),
				CalledCountryID:   mapOptionalInt("CalledCountryID", row),
				IsPremium:         mapBool("IsPremium", row),
				TrafficSegmentID:  mapOptionalInt("TrafficSegmentID", row),
				IMSICountType:     mapOptionalInt("IMSICountType", row),
				VolumeActual:      forecastedVolumeActual,
			}

			distributionRecords[idx] = _distributionRecord

			idx++
		}
	}

	return distributionRecords
}

func mapOptionalInt(colName string, row map[string]interface{}) *int {

	var value *int

	rawValue := row[colName].(int)

	if rawValue == -1 {
		value = nil
	} else {
		value = &rawValue
	}

	return value
}

func mapBool(colName string, row map[string]interface{}) *bool {
	var value *bool

	if row[colName].(int) == -1 {
		value = nil
	} else if row[colName].(int) == 0 {
		f := false
		value = &f
	} else {
		t := true
		value = &t
	}

	return value
}