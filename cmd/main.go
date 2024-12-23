package main

import (
	"log"

	"forecasting/app/calculation"
	"forecasting/app/persistence"
	"forecasting/internal/config"
	"forecasting/internal/infra"
)

func main() {
	cfg, err := config.InitConfig()

	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	var _infra *infra.Infra

	_infra, err = infra.NewInfra(cfg)

	if err != nil {
		log.Println("Error creating infra", err)
	}

	defer _infra.Shutdown()

	maProvider := persistence.NewPostgresMAProvider(_infra.GetDB())

	btrProvider := persistence.NewPostgresBudgetTrafficProvider(_infra.GetDB())

	forecastingService := calculation.NewForecastingService(maProvider, btrProvider)

	forecastRuleRepo := persistence.NewPostgresForecastRuleRepository(_infra.GetDB())

	calcCoordinator := calculation.NewForecastRuleCalculationCoordinator(
		forecastRuleRepo,
		forecastingService,
		btrProvider,
	)

	calcCoordinator.CalculateAll()
}
