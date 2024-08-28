package main

import (
	"log"

	"forecasting/app/calculation"
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

	maProvider := providers.NewPostgresMAProvider(_infra.GetDB())

	btrProvider := providers.NewPostgresBudgetTrafficProvider(_infra.GetDB())

	forecastingService := calculation.NewForecastingService(maProvider, btrProvider)

	forecastRuleRepo := repositories.NewPostgresForecastRuleRepository(_infra.GetDB())

	calcCoordinator := calculation.NewForecastRuleCalculationCoordinator(
		forecastRuleRepo,
		forecastingService,
		btrProvider,
	)

	calcCoordinator.CalculateAll()
}
