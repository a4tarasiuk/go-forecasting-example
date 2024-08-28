package calculation

import (
	"forecasting/app/calculation/dto"
	"forecasting/app/domain/models"
	"forecasting/core"
)

type BudgetTrafficRecordMapper struct {
}

func NewBudgetTrafficRecordMapper() BudgetTrafficRecordMapper {
	return BudgetTrafficRecordMapper{}
}

func (m BudgetTrafficRecordMapper) FromDistributionToBudgetTrafficRecord(
	forecastRule *models.ForecastRule,
	record dto.DistributionRecord,
) models.BudgetTrafficRecord {

	cd := core.GetDefaultCDByServiceType(forecastRule.ServiceType)

	var cdVal int64

	if cd != nil {
		cdVal = int64(*cd)
	}

	return models.BudgetTrafficRecord{
		BudgetSnapshotID:  record.BudgetSnapshotID,
		HomeOperatorID:    record.HomeOperatorID,
		PartnerOperatorID: record.PartnerOperatorID,
		TrafficType:       2, // FORECASTED
		TrafficDirection:  forecastRule.TrafficDirection,
		Month:             record.Month,
		ServiceType:       forecastRule.ServiceType,
		CallDestination:   &cdVal,
		CalledCountryID:   nil,
		IsPremium:         nil,
		IMSICountType:     nil,
		TrafficSegmentID:  nil,
		VolumeActual:      record.VolumeActual,
	}
}
