package main

import (
	"fmt"

	"forecasting/calculation"
	"forecasting/calculation/distribution_models"
	"forecasting/core"
	"forecasting/core/types"
	"forecasting/rules"
	"forecasting/traffic"
	"github.com/golang-module/carbon/v2"
)

func main() {
	startDate := carbon.Parse("2023-01-01").ToDateStruct()
	endDate := carbon.Parse("2023-01-01").ToDateStruct()

	lhm := carbon.Parse("2022-12-01").ToDateStruct()

	nMonths := 3

	forecastRule := &rules.ForecastRule{
		ID:                                   1,
		HomeOperators:                        []int{1},
		PartnerOperators:                     []int{2, 3},
		Period:                               types.NewPeriod(startDate, endDate),
		TrafficDirection:                     core.InboundTrafficDirection,
		ServiceType:                          core.VoiceMO,
		ForecastModel:                        core.ManualVolumeForecastModel,
		DistributionModel:                    core.MovingAverageDistributionModel,
		Volume:                               2500.0,
		DistributionModelMovingAverageMonths: &nMonths,
		LHM:                                  &lhm,
	}

	// manualVolumeModel := forecast_models.NewManualVolume()
	//
	// trafficRecords, err := manualVolumeModel.Calculate(forecastRule)
	//

	movingAverageDistributionModel := distribution_models.NewMovingAverage()

	forecastRecords := []calculation.ForecastRecord{
		{VolumeActual: forecastRule.Volume, Month: carbon.Parse("2023-01-01").ToDateStruct()},
	}

	trafficRecords, err := movingAverageDistributionModel.Apply(forecastRule, forecastRecords)

	fmt.Println(err)

	totalForecastVolumeAfterDistribution := traffic.CalculateTotalHistoricalVolume(trafficRecords)

	fmt.Println("Total F volume", totalForecastVolumeAfterDistribution)

	for _, record := range trafficRecords {
		fmt.Printf("%+v\n", record)
	}
}
