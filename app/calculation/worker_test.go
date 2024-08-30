package calculation

import (
	"testing"

	"forecasting/app/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestWorker_Evaluates_ForecastRule(t *testing.T) {
	budgetTrafficRecords := []models.BudgetTrafficRecord{{BudgetSnapshotID: 1}}

	forecastingService := &fakeForecastingService{btrRecordsToReturn: budgetTrafficRecords}

	forecastRule := &models.ForecastRule{ID: 1}

	forecastRulesCh := make(chan *models.ForecastRule)

	budgetTrafficRecordsCh := make(chan []models.BudgetTrafficRecord)

	go ForecastWorker(forecastingService, forecastRulesCh, budgetTrafficRecordsCh)

	forecastRulesCh <- forecastRule

	select {
	case resultRecords := <-budgetTrafficRecordsCh:
		assert.Equal(t, budgetTrafficRecords, resultRecords)
	}
}
