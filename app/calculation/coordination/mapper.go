package coordination

import (
	"database/sql"

	"forecasting/app/calculation"
	"forecasting/app/domain/models"
	"forecasting/core"
)

type BudgetTrafficRecordMapper struct {
	budgetSnapshotID int64
}

func NewBudgetTrafficRecordMapper(db *sql.DB, budgetID int) BudgetTrafficRecordMapper {
	var budgetSnapshotID int64

	// 2 - CALCULATION snapshot
	rows, _ := db.Query("SELECT id FROM budget_snapshots WHERE budget_id = $1 AND type = 2", budgetID)
	rows.Next()
	rows.Scan(&budgetSnapshotID)

	return BudgetTrafficRecordMapper{budgetSnapshotID: budgetSnapshotID}
}

func (m BudgetTrafficRecordMapper) FromDistributionToBudgetTrafficRecord(
	forecastRule *models.ForecastRule,
	record calculation.DistributionRecord,
) models.BudgetTrafficRecord {

	cd := core.GetDefaultCDByServiceType(forecastRule.ServiceType)

	var cdVal int64

	if cd != nil {
		cdVal = int64(*cd)
	}

	return models.BudgetTrafficRecord{
		BudgetSnapshotID:  m.budgetSnapshotID,
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
