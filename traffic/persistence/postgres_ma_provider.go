package persistence

import (
	"database/sql"
	"fmt"
	"log"

	"forecasting/core/types"
	"forecasting/postgres"
	"forecasting/rules"
	"forecasting/traffic"
	"github.com/golang-module/carbon/v2"
	"github.com/lib/pq"
)

type postgresMAProvider struct {
	db *sql.DB
}

func NewPostgresMAProvider() *postgresMAProvider {
	db := postgres.CreateDBConnection()

	return &postgresMAProvider{db: db}
}

func (p *postgresMAProvider) Get(
	forecastRule *rules.ForecastRule,
	period types.Period,
) []traffic.MonthlyAggregationRecord {

	rows, err := p.db.Query(
		getManyAggSQLQuery,
		forecastRule.BudgetID,
		pq.Array(forecastRule.HomeOperators),
		pq.Array(forecastRule.PartnerOperators),
		forecastRule.TrafficDirection,
		forecastRule.ServiceType,
		period.StartDate.ToDateString(),
		period.EndDate.ToDateString(),
	)
	defer rows.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var monthStr string
	var volume float64

	aggregations := make([]traffic.MonthlyAggregationRecord, 0)

	for rows.Next() {
		_err := rows.Scan(&monthStr, &volume)

		if _err != nil {
			log.Fatal(_err)
		}

		agg := traffic.MonthlyAggregationRecord{
			VolumeActual: volume,
			Month:        carbon.Parse(monthStr).ToDateStruct(),
		}

		fmt.Printf("%+v\n", agg)

		aggregations = append(aggregations, agg)
	}

	return aggregations
}

const getManyAggSQLQuery = `
SELECT 
    btr.traffic_month, 
    SUM(btr.volume_actual) AS "volume_actual"
FROM 
	budget_traffic_records btr
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
GROUP BY btr.traffic_month
ORDER BY btr.traffic_month
`
