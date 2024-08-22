package main

import (
	"fmt"

	"forecasting/calculation/forecast_model"
	"forecasting/core"
	"forecasting/core/types"
	"forecasting/rules"
	"github.com/golang-module/carbon/v2"
)

func main() {
	startDate := carbon.Parse("2024-08-01").ToDateStruct()
	endDate := carbon.Parse("2024-10-01").ToDateStruct()

	forecastRule := &rules.ForecastRule{
		ID:                                   1,
		HomeOperators:                        []int{1, 2, 3, 4, 5},
		PartnerOperators:                     []int{6, 7},
		Period:                               types.Period{StartDate: startDate, EndDate: endDate},
		TrafficDirection:                     core.InboundTrafficDirection,
		ServiceType:                          core.VoiceMO,
		ForecastModel:                        core.ManualVolumeForecastModel,
		DistributionModel:                    core.MovingAverageDistributionModel,
		Volume:                               1200.0,
		DistributionModelMovingAverageMonths: nil,
	}

	manualVolumeModel := forecast_model.NewManualVolume()

	trafficRecords := manualVolumeModel.Calculate(forecastRule)

	for _, record := range trafficRecords {
		fmt.Printf("%+v\n", record)
	}
}
