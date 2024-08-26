package main

import (
	"log"

	"forecasting/app/budget_defaults"
	"forecasting/app/calculation/coordination"
	"forecasting/app/persistence/providers"
	"forecasting/app/persistence/repositories"
	"forecasting/internal/config"
	"forecasting/internal/infra"
)

func main() {
	cfg, err := config.InitConfig()

	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	_infra, _ := infra.NewInfra(cfg)

	defer _infra.Shutdown()

	btrMapper := coordination.NewBudgetTrafficRecordMapper(_infra.GetDB(), budget_defaults.BudgetID)

	maProvider := providers.NewPostgresMAProvider(_infra.GetDB())

	btrProvider := providers.NewPostgresBudgetTrafficProvider(_infra.GetDB())

	forecastingService := coordination.NewForecastingService(
		btrMapper,
		maProvider,
		btrProvider,
	)

	forecastRuleRepo := repositories.NewPostgresForecastRuleRepository(_infra.GetDB())

	calcCoordinator := coordination.NewForecastRuleCalculationCoordinator(
		forecastRuleRepo,
		forecastingService,
		btrProvider,
	)

	calcCoordinator.CalculateAll()
}
