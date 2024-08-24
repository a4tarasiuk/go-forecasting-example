package main

import (
	"forecasting/calculation/coordination"
)

func main() {
	calcCoordinator := coordination.NewForecastRuleCalculationCoordinator()

	calcCoordinator.CalculateAll()
}
