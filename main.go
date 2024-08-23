package main

import (
	"fmt"

	"forecasting/calculation"
	"forecasting/calculation/distribution_models"
	"forecasting/core"
	"forecasting/core/types"
	"forecasting/rules"
	"github.com/golang-module/carbon/v2"
)

func main() {
	startDate := carbon.Parse("2024-08-05").ToDateStruct()
	endDate := carbon.Parse("2024-09-11").ToDateStruct()

	lhm := carbon.Parse("2024-08-01").ToDateStruct()

	nMonths := 3

	forecastRule := &rules.ForecastRule{
		ID:                                   1,
		HomeOperators:                        []int{1, 2},
		PartnerOperators:                     []int{3, 4},
		Period:                               types.NewPeriod(startDate, endDate),
		TrafficDirection:                     core.InboundTrafficDirection,
		ServiceType:                          core.VoiceMO,
		ForecastModel:                        core.ManualVolumeForecastModel,
		DistributionModel:                    core.MovingAverageDistributionModel,
		Volume:                               1200.0,
		DistributionModelMovingAverageMonths: &nMonths,
		LHM:                                  &lhm,
	}

	// manualVolumeModel := forecast_models.NewManualVolume()
	//
	// trafficRecords, err := manualVolumeModel.Calculate(forecastRule)
	//

	movingAverageDistributionModel := distribution_models.NewMovingAverage()

	forecastRecords := []calculation.ForecastRecord{
		{VolumeActual: 500.0, Month: carbon.Parse("2024-08-01").ToDateStruct()},
		{VolumeActual: 100.0, Month: carbon.Parse("2024-09-01").ToDateStruct()},
	}

	trafficRecords, err := movingAverageDistributionModel.Apply(forecastRule, forecastRecords)

	fmt.Println(err)

	for _, record := range trafficRecords {
		fmt.Printf("%+v\n", record)
	}
}
