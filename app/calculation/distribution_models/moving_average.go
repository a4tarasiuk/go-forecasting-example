package distribution_models

import (
	"forecasting/app/calculation/calc_utils"
	"forecasting/app/calculation/dto"
	"forecasting/app/domain/models"
	"forecasting/app/providers"
	"forecasting/core"
	"forecasting/core/types"
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

type movingAverage struct {
	budgetTrafficProvider providers.BudgetTrafficProvider
}

func NewMovingAverage(provider providers.BudgetTrafficProvider) movingAverage {
	return movingAverage{budgetTrafficProvider: provider}
}

func (ma *movingAverage) Apply(
	forecastRule *models.ForecastRule,
	forecastRecords []dto.ForecastRecord,
) (
	[]dto.DistributionRecord,
	error,
) {
	historicalTrafficRecords := ma.loadHistoricalTraffic(forecastRule)

	if calc_utils.ShouldCalculateWithoutTraffic(historicalTrafficRecords) {
		return ma.calculateWithoutTraffic(forecastRule, forecastRecords), nil
	}

	return ma.calculateWithTraffic(forecastRecords, historicalTrafficRecords), nil
}

func (ma *movingAverage) loadHistoricalTraffic(forecastRule *models.ForecastRule) []models.BudgetTrafficRecord {
	nMonthPeriodEndDate := forecastRule.LHM

	monthsToSub := *forecastRule.DistributionModelMovingAverageMonths - 1

	nMonthPeriodStartDate := nMonthPeriodEndDate.SubMonths(monthsToSub).ToDateStruct()

	nMonthPeriod := types.NewPeriod(nMonthPeriodStartDate, nMonthPeriodEndDate)

	options := providers.BudgetTrafficOptions{
		ForecastRule:   forecastRule,
		Period:         &nMonthPeriod,
		HistoricalOnly: true,
	}

	budgetTrafficRecords := ma.budgetTrafficProvider.Get(options)

	return budgetTrafficRecords
}

func (ma *movingAverage) calculateWithoutTraffic(
	forecastRule *models.ForecastRule,
	forecastRecords []dto.ForecastRecord,
) []dto.DistributionRecord {

	totalHomeOperators := int64(len(forecastRule.HomeOperators))
	totalPartnerOperators := int64(len(forecastRule.PartnerOperators))
	totalMonths := int64(len(forecastRecords))

	totalResultRecords := totalHomeOperators * totalPartnerOperators * totalMonths

	distributionRecords := make([]dto.DistributionRecord, totalResultRecords)

	idx := 0

	for _, forecastRecord := range forecastRecords {
		totalRecordsWithinMonth := totalHomeOperators * totalPartnerOperators

		for _, homeOperatorID := range forecastRule.HomeOperators {
			for _, partnerOperatorID := range forecastRule.PartnerOperators {

				_distributionRecord := dto.DistributionRecord{
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
	forecastRecords []dto.ForecastRecord,
	historicalRecords []models.BudgetTrafficRecord,
) []dto.DistributionRecord {

	rawHistoricalTrafficRecords := make([]map[string]interface{}, len(historicalRecords))

	for idx, _historicalRecord := range historicalRecords {
		rawHistoricalTrafficRecords[idx] = _historicalRecord.Serialize()
	}

	typesMap := map[string]series.Type{
		"TrafficType":      series.Int,
		"TrafficDirection": series.Int,
		"TrafficSegmentID": series.Int,
		"ServiceType":      series.Int,
		"CalledCountryID":  series.Int,
	}

	historicalTrafficDF := dataframe.LoadMaps(rawHistoricalTrafficRecords, dataframe.WithTypes(typesMap))

	// Aggregation
	historicalTrafficDF = historicalTrafficDF.GroupBy(
		"HomeOperatorID",
		"PartnerOperatorID",
		"CallDestination",
		"CalledCountryID",
		"IsPremium",
		"TrafficSegmentID",
		"IMSICountType",
	).Aggregation(
		[]dataframe.AggregationType{dataframe.Aggregation_SUM}, []string{"VolumeActual"},
	)

	// Coefficient calculation
	totalHistoricalVolume := calc_utils.CalculateTotalHistoricalVolume(historicalRecords)

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
	forecastDfMap := make(map[dto.ForecastRecord]dataframe.DataFrame, len(forecastRecords))

	for _, forecastRecord := range forecastRecords {
		forecastDfMap[forecastRecord] = historicalTrafficDF.Copy()
	}

	var distributionRecords []dto.DistributionRecord

	for forecastRecord, df := range forecastDfMap {

		for _, row := range df.Maps() {
			forecastedVolumeActual := row["ForecastCoefficient"].(float64) * forecastRecord.VolumeActual

			_distributionRecord := dto.DistributionRecord{
				HomeOperatorID:    int64(row["HomeOperatorID"].(int)),
				PartnerOperatorID: int64(row["PartnerOperatorID"].(int)),
				Month:             forecastRecord.Month,
				CallDestination:   mapOptionalInt("CallDestination", row),
				CalledCountryID:   mapOptionalInt("CalledCountryID", row),
				IsPremium:         mapBool("IsPremium", row),
				TrafficSegmentID:  mapOptionalInt("TrafficSegmentID", row),
				IMSICountType:     mapOptionalInt("IMSICountType", row),
				VolumeActual:      forecastedVolumeActual,
			}

			distributionRecords = append(distributionRecords, _distributionRecord)
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
