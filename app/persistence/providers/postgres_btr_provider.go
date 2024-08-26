package providers

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"forecasting/app/budget_defaults"
	"forecasting/app/domain/models"
	"forecasting/app/providers"
	"forecasting/core"
	"github.com/golang-module/carbon/v2"
	"github.com/lib/pq"
)

type postgresBudgetTrafficProvider struct {
	db *sql.DB
}

func NewPostgresBudgetTrafficProvider(db *sql.DB) *postgresBudgetTrafficProvider {
	return &postgresBudgetTrafficProvider{db: db}
}

func (p *postgresBudgetTrafficProvider) Get(options providers.BudgetTrafficOptions) []models.BudgetTrafficRecord {

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

	var budgetTrafficRecords []models.BudgetTrafficRecord

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
			log.Println(_err)
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

		record := models.BudgetTrafficRecord{
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

func (p *postgresBudgetTrafficProvider) CreateMany(records []models.BudgetTrafficRecord) {
	defer recover()

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
	                        is_premium,
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

	if len(records) == 0 {
		return
	}

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
			0,
			r.IsPremium,
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

		numFields := 18 // the number of fields you are inserting
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
		log.Println(err)
	}
}

func (p *postgresBudgetTrafficProvider) ClearForecasted() {
	p.db.Query(deleteManySQLQuery, budget_defaults.BudgetID)
}

func (p *postgresBudgetTrafficProvider) CountForecasted() {
	rows, _ := p.db.Query("SELECT COUNT(id) FROM budget_traffic_records WHERE budget_snapshot_id = 498 AND traffic_type = 2")

	defer rows.Close()

	var totalRecords int64

	rows.Next()
	rows.Scan(&totalRecords)

	fmt.Println("total created forecasted records - ", totalRecords)
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
	    btr.budget_snapshot_id = bs.id AND bs.type = 2
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
