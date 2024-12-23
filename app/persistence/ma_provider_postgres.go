package persistence

import (
	"database/sql"
	"log"

	"forecasting/app/domain/models"
	"forecasting/core/types"
	"github.com/golang-module/carbon/v2"
	"github.com/lib/pq"
)

type postgresMAProvider struct {
	db *sql.DB
}

func NewPostgresMAProvider(db *sql.DB) *postgresMAProvider {
	return &postgresMAProvider{db: db}
}

func (p *postgresMAProvider) Get(
	forecastRule *models.ForecastRule,
	period types.Period,
) []models.MonthlyAggregationRecord {

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

	var aggregations []models.MonthlyAggregationRecord

	for rows.Next() {
		_err := rows.Scan(&monthStr, &volume)

		if _err != nil {
			log.Println(_err)
		}

		agg := models.MonthlyAggregationRecord{
			VolumeActual: volume,
			Month:        carbon.Parse(monthStr).ToDateStruct(),
		}

		aggregations = append(aggregations, agg)
	}

	return aggregations
}

func (p *postgresMAProvider) GetLast(
	forecastRule *models.ForecastRule,
	period types.Period,
) []models.MonthlyAggregationRecord {
	rows, err := p.db.Query("SELECT start_date FROM budgets WHERE id = $1", forecastRule.BudgetID)
	defer rows.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var budgetStartDateStr string

	rows.Next()
	rows.Scan(&budgetStartDateStr)

	budgetStartDate := carbon.Parse(budgetStartDateStr).ToDateStruct()

	fullPeriod := types.NewPeriod(budgetStartDate, period.EndDate)

	aggregations := p.Get(forecastRule, fullPeriod)

	searchPeriod := period

	if searchPeriod.MoreThenAYear() {
		searchPeriod = searchPeriod.CutToOneYear()
	}

	if searchPeriod.StartDate.Compare("<", budgetStartDate.Carbon) {
		searchPeriod.StartDate = budgetStartDate
	}

	var lastRecords []models.MonthlyAggregationRecord

	for budgetStartDate.Compare("<=", searchPeriod.StartDate.Carbon) {
		lastRecords = []models.MonthlyAggregationRecord{}

		for _, r := range aggregations {
			if searchPeriod.Contains(r.Month) {
				lastRecords = append(lastRecords, r)
			}
		}

		if len(lastRecords) > 0 {
			break
		}

		searchPeriod = searchPeriod.PreviousYear()
	}

	return lastRecords
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
	    btr.budget_snapshot_id = bs.id AND bs.type = 2
WHERE bs.budget_id = $1
  	AND btr.home_operator_id = ANY($2)
  	AND btr.partner_operator_id = ANY($3)
  	AND btr.traffic_direction = $4
  	AND btr.service_type = $5
  	AND btr.traffic_month BETWEEN $6 AND $7
GROUP BY btr.traffic_month
ORDER BY btr.traffic_month
`
