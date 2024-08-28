package repositories

import (
	"database/sql"
	"fmt"
	"log"

	"forecasting/app/budget_defaults"
	"forecasting/app/domain/models"
	"forecasting/core"
	"forecasting/core/types"
	"github.com/golang-module/carbon/v2"
	"github.com/lib/pq"
)

type PostgresForecastRuleRepository struct {
	db *sql.DB
}

func NewPostgresForecastRuleRepository(db *sql.DB) *PostgresForecastRuleRepository {
	return &PostgresForecastRuleRepository{db: db}
}

func (r *PostgresForecastRuleRepository) GetMany() []*models.ForecastRule {
	budgetSnapshotID := r.getBudgetSnapshotID(budget_defaults.BudgetID)

	rows, err := r.db.Query(getManySQLQuery, budget_defaults.BudgetID)
	defer rows.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var ID, budgetID int64
	var sqlHomeOperators, sqlPartnerOperators []sql.NullInt64
	var startDate, endDate string
	var serviceType, trafficDirection, forecastModel int
	var volume float64

	var distributionModel int
	var distributionMovingAverageMonths *int

	var forecastRules []*models.ForecastRule

	for rows.Next() {
		_err := rows.Scan(
			&ID,
			pq.Array(&sqlHomeOperators),
			pq.Array(&sqlPartnerOperators),
			&startDate,
			&endDate,
			&trafficDirection,
			&serviceType,
			&forecastModel,
			&volume,
			&distributionModel,
			&distributionMovingAverageMonths,
			&budgetID,
		)

		if _err != nil {
			log.Println(_err)
		}

		rule := models.ForecastRule{
			ID:               ID,
			BudgetID:         budgetID,
			BudgetSnapshotID: budgetSnapshotID,
			HomeOperators:    mapInt64Array(sqlHomeOperators),
			PartnerOperators: mapInt64Array(sqlPartnerOperators),
			Period: types.NewPeriod(
				carbon.Parse(startDate).ToDateStruct(),
				carbon.Parse(endDate).ToDateStruct(),
			),
			TrafficDirection:                     core.TrafficDirection(trafficDirection),
			ServiceType:                          core.ServiceType(serviceType),
			ForecastModel:                        core.ForecastModel(forecastModel),
			DistributionModel:                    core.DistributionModel(distributionModel),
			Volume:                               volume,
			DistributionModelMovingAverageMonths: distributionMovingAverageMonths,
			LHM:                                  carbon.Parse("2024-02-01").ToDateStruct(), // TODO:
		}

		forecastRules = append(forecastRules, &rule)
	}

	return forecastRules
}

func (r *PostgresForecastRuleRepository) getBudgetSnapshotID(budgetID int) int64 {
	var budgetSnapshotID int64

	// 2 - CALCULATION snapshot
	rows, err := r.db.Query("SELECT id FROM budget_snapshots WHERE budget_id = $1 AND type = 2", budgetID)

	defer rows.Close()

	if err != nil {
		fmt.Println(err)
	}

	rows.Next()
	rows.Scan(&budgetSnapshotID)

	return budgetSnapshotID
}

const getManySQLQuery = `
SELECT 
    id,
    (
        SELECT 
            ARRAY_AGG(DISTINCT fho.operator_id) FILTER ( WHERE fho.operator_id IS NOT NULL )::BIGINT[]
		FROM 
		    forecast_rules_home_operators fho
        WHERE fho.forecastrule_id = forecast_rules.id
    ) AS "home_operators",
    CASE partner_type
        WHEN 1
            THEN (
				SELECT 
				    ARRAY_AGG(DISTINCT fpo.operator_id) FILTER ( WHERE fpo.operator_id IS NOT NULL )::BIGINT[]
				FROM 
				    forecast_rules_partner_operators fpo
				WHERE forecastrule_id = forecast_rules.id
            )
		ELSE (
				SELECT 
				    ARRAY_AGG(DISTINCT op.id) FILTER ( WHERE op.id IS NOT NULL )::BIGINT[]
				FROM 
				    operators op
				INNER JOIN 
				    forecast_rules_partner_countries fpc
					ON 
						fpc.forecastrule_id = forecast_rules.id
				WHERE fpc.country_id = op.country_id
		    )
        END AS "partner_operators",
    start_date,
    end_date,
    traffic_direction,
    service_type,
    model,
    volume,
    distribution_model,
    distribution_moving_average_months,
    budget_id
FROM 
    forecast_rules
WHERE budget_id = $1
ORDER BY start_date`

func mapInt64Array(arr []sql.NullInt64) []int64 {
	var values []int64

	for _, obj := range arr {
		val, _ := obj.Value()
		values = append(values, val.(int64))
	}

	return values
}
