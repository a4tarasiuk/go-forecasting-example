package persistence

import (
	"database/sql"
	"log"

	"forecasting/core"
	"forecasting/traffic"
	"github.com/golang-module/carbon/v2"
	"github.com/lib/pq"
)

type postgresBudgetTrafficProvider struct {
	db *sql.DB
}

func NewPostgresBudgetTrafficProvider() *postgresBudgetTrafficProvider {
	// TODO: Unify db session

	connStr := "postgresql://postgres:postgres@localhost/test?sslmode=disable" // TODO: Env

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return &postgresBudgetTrafficProvider{db: db}
}

func (p *postgresBudgetTrafficProvider) Get(options traffic.BudgetTrafficOptions) []traffic.BudgetTrafficRecord {

	// BudgetTrafficOptions.HistoricalOnly is enabled by default. It is hardcoded in SQL query

	rows, err := p.db.Query(
		getManySQLQuery,
		options.ForecastRule.BudgetID,
		pq.Array(options.ForecastRule.HomeOperators),
		pq.Array(options.ForecastRule.PartnerOperators),
		options.ForecastRule.TrafficDirection,
		options.ForecastRule.ServiceType,
		options.ForecastRule.Period.StartDate.ToDateString(),
		options.ForecastRule.Period.EndDate.ToDateString(),
	)
	defer rows.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var budgetSnapshotID, homeOperatorID, partnerOperatorID int64
	var trafficSegmentID, calledCountryID sql.NullInt64
	var trafficType, trafficDirection, serviceType byte
	var imsiCountType, callDestination sql.NullByte
	var sqlTrafficMonth string
	var isPremium sql.NullBool
	var volume float64

	budgetTrafficRecords := make([]traffic.BudgetTrafficRecord, 0)

	for rows.Next() {
		_err := rows.Scan(
			&budgetSnapshotID,
			&homeOperatorID,
			&partnerOperatorID,
			&trafficType,
			&sqlTrafficMonth,
			&trafficDirection,
			&trafficSegmentID,
			&serviceType,
			&volume,
			&callDestination,
			&calledCountryID,
			&isPremium,
			&imsiCountType,
		)

		if _err != nil {
			log.Fatal(_err)
		}

		var cdValue *byte
		if callDestination.Valid {
			cdValue = &callDestination.Byte
		}

		var calledCountryIDValue *int64
		if calledCountryID.Valid {
			calledCountryIDValue = &calledCountryID.Int64
		}

		var isPremiumValue *bool
		if isPremium.Valid {
			isPremiumValue = &isPremium.Bool
		}

		var imsiCountTypeValue *byte
		if imsiCountType.Valid {
			imsiCountTypeValue = &imsiCountType.Byte
		}

		var trafficSegmentIDValue *int64
		if trafficSegmentID.Valid {
			trafficSegmentIDValue = &trafficSegmentID.Int64
		}

		record := traffic.BudgetTrafficRecord{
			BudgetSnapshotID:  budgetSnapshotID,
			HomeOperatorID:    homeOperatorID,
			PartnerOperatorID: partnerOperatorID,
			TrafficType:       trafficType,
			TrafficDirection:  core.TrafficDirection(trafficDirection),
			Month:             carbon.Parse(sqlTrafficMonth).ToDateStruct(),
			ServiceType:       core.ServiceType(serviceType),
			CallDestination:   cdValue,
			CalledCountryID:   calledCountryIDValue,
			IsPremium:         isPremiumValue,
			IMSICountType:     imsiCountTypeValue,
			TrafficSegmentID:  trafficSegmentIDValue,
			VolumeActual:      volume,
		}

		budgetTrafficRecords = append(budgetTrafficRecords, record)
	}

	return budgetTrafficRecords
}

const getManySQLQuery = `
SELECT
    btr.budget_snapshot_id,
	btr.home_operator_id,
	btr.partner_operator_id,
	btr.traffic_type,
    btr.traffic_month,
    btr.traffic_direction,
	btr.traffic_segment_id,
    btr.service_type,
    btr.volume_actual,
    btr.call_destination,
    btr.called_country_id,
    btr.is_premium,
    btr.imsi_count_type
FROM budget_traffic_records btr
	INNER JOIN 
    	budget_snapshots bs
	ON 
	    btr.budget_snapshot_id = bs.id AND bs.type = 1
WHERE bs.budget_id = $1
  	AND btr.home_operator_id = ANY($2)
  	AND btr.partner_operator_id = ANY($3)
  	AND btr.traffic_direction = $4
  	AND btr.service_type = $5
  	AND btr.traffic_month BETWEEN $6 AND $7
  	AND btr.traffic_type = 1
`
