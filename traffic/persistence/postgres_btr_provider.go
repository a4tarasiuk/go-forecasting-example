package persistence

import (
	"database/sql"
	"log"
	"strconv"

	"forecasting/budget_defaults"
	"forecasting/core"
	"forecasting/postgres"
	"forecasting/traffic"
	"github.com/golang-module/carbon/v2"
	"github.com/lib/pq"
)

type postgresBudgetTrafficProvider struct {
	db *sql.DB
}

func NewPostgresBudgetTrafficProvider() *postgresBudgetTrafficProvider {
	return &postgresBudgetTrafficProvider{db: postgres.DB}
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
	var trafficSegmentID, calledCountryID, callDestination sql.NullInt64
	var trafficType, trafficDirection, serviceType byte
	var imsiCountType sql.NullByte
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

		var cdValue *int64
		if callDestination.Valid {
			cdValue = &callDestination.Int64
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

func (p *postgresBudgetTrafficProvider) CreateMany(records []traffic.BudgetTrafficRecord) {
	insertManySQLQuery := `
INSERT INTO
	budget_traffic_records (
	                        budget_snapshot_id,
	                    	home_operator_id,
	                        partner_operator_id,
	                        traffic_month,
	                        traffic_type,
	                        traffic_direction,
	                        service_type,
	                        call_destination,
	                        called_country_id,
	                        is_premium,
	                        traffic_segment_id,
	                        imsi_count_type,
	                        volume_actual,
	                        volume_billed,
	                        tap_charge_net,
	                        tap_charge_gross,
	                        charge_net,
	                        charge_gross,
	                        created_at,
	                        updated_at
	                        ) VALUES `

	vals := []interface{}{}

	for i, r := range records {
		vals = append(
			vals,
			r.BudgetSnapshotID,
			r.HomeOperatorID,
			r.PartnerOperatorID,
			r.Month.ToDateString(),
			r.TrafficType,
			r.TrafficDirection,
			r.ServiceType,
			r.CallDestination,
			r.CalledCountryID,
			r.IsPremium,
			r.TrafficSegmentID,
			r.IMSICountType,
			r.VolumeActual,
			0,
			0,
			0,
			0,
			0,
			carbon.Now().ToDateString(),
			carbon.Now().ToDateString(),
		)

		numFields := 20 // the number of fields you are inserting
		n := i * numFields

		insertManySQLQuery += `(`
		for j := 0; j < numFields; j++ {
			insertManySQLQuery += `$` + strconv.Itoa(n+j+1) + `,`
		}
		insertManySQLQuery = insertManySQLQuery[:len(insertManySQLQuery)-1] + `),`
	}

	insertManySQLQuery = insertManySQLQuery[0 : len(insertManySQLQuery)-1]

	stmt, _ := p.db.Prepare(insertManySQLQuery)

	_, err := stmt.Exec(vals...)

	if err != nil {
		log.Fatal(err)
	}
}

func (p *postgresBudgetTrafficProvider) ClearForecasted() {
	p.db.Query(deleteManySQLQuery, budget_defaults.BudgetID)
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

const deleteManySQLQuery = `
DELETE FROM budget_traffic_records btr
USING budget_snapshots bs
WHERE bs.id = btr.budget_snapshot_id 
  AND bs.type = 2 
  AND bs.budget_id = $1
  AND btr.traffic_type = 2
`
