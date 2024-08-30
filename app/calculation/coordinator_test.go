package calculation

import (
	"testing"

	"forecasting/app/domain/models"
	"forecasting/app/persistence"
	"forecasting/app/providers"
	"github.com/stretchr/testify/assert"
)

type fakeForecastingService struct {
	calculatedRuleIDs []int64

	btrRecordsToReturn []models.BudgetTrafficRecord
}

func (s *fakeForecastingService) Evaluate(forecastRule *models.ForecastRule) []models.BudgetTrafficRecord {
	s.calculatedRuleIDs = append(s.calculatedRuleIDs, forecastRule.ID)

	return s.btrRecordsToReturn
}

func Test_CalculateAll_Rules_Are_Calculated(t *testing.T) {
	forecastRule1 := &models.ForecastRule{ID: 1}
	forecastRule2 := &models.ForecastRule{ID: 2}
	forecastRule3 := &models.ForecastRule{ID: 3}

	ruleRepo := persistence.NewForecastRuleInMemoryRepository(
		[]*models.ForecastRule{
			forecastRule1,
			forecastRule2,
			forecastRule3,
		},
	)

	forecastingService := &fakeForecastingService{}

	budgetTrafficProvider := persistence.NewBudgetTrafficInMemoryProvider(nil)

	coordinator := NewForecastRuleCalculationCoordinator(
		ruleRepo,
		forecastingService,
		budgetTrafficProvider,
	)

	coordinator.CalculateAll()

	expectedCalculatedRuleIDs := []int64{forecastRule1.ID, forecastRule2.ID, forecastRule3.ID}

	assert.Equal(t, expectedCalculatedRuleIDs, forecastingService.calculatedRuleIDs)
}

func Test_CalculateAll_BTRs_Are_Created(t *testing.T) {
	forecastRule := &models.ForecastRule{ID: 1}

	budgetTrafficRecords := []models.BudgetTrafficRecord{
		{},
		{},
	}

	budgetTrafficProvider := persistence.NewBudgetTrafficInMemoryProvider(nil)

	assert.Zero(t, len(budgetTrafficProvider.Get(providers.BudgetTrafficOptions{})))

	coordinator := NewForecastRuleCalculationCoordinator(
		persistence.NewForecastRuleInMemoryRepository([]*models.ForecastRule{forecastRule}),
		&fakeForecastingService{btrRecordsToReturn: budgetTrafficRecords},
		budgetTrafficProvider,
	)

	coordinator.CalculateAll()

	expectedTotalRecords := 2

	assert.Equal(t, expectedTotalRecords, len(budgetTrafficProvider.Get(providers.BudgetTrafficOptions{})))
}

func Test_Forecast_Records_Are_Cleared_Before_Calculation(t *testing.T) {
	forecastRule := &models.ForecastRule{ID: 1}

	budgetTrafficRecords := []models.BudgetTrafficRecord{
		{TrafficType: 1},
		{TrafficType: 2},
		{TrafficType: 1},
		{TrafficType: 2},
	}

	budgetTrafficProvider := persistence.NewBudgetTrafficInMemoryProvider(budgetTrafficRecords)

	coordinator := NewForecastRuleCalculationCoordinator(
		persistence.NewForecastRuleInMemoryRepository([]*models.ForecastRule{forecastRule}),
		&fakeForecastingService{},
		budgetTrafficProvider,
	)

	coordinator.CalculateAll()

	assert.Zero(t, budgetTrafficProvider.CountForecasted())
}
