package forecast_models

import (
	"testing"

	"forecasting/app/calculation/dto"
	"forecasting/app/domain/models"
	"forecasting/app/persistence"
	"forecasting/core/types"
	"github.com/stretchr/testify/assert"
)

func Test_when_there_is_no_historical_traffic(t *testing.T) {
	forecastRule := &models.ForecastRule{
		Period: types.Period{
			StartDate: types.ToDate("2024-09-01"),
			EndDate:   types.ToDate("2024-12-01"),
		},
		Volume: 100.0,
		LHM:    types.ToDate("2024-08-01"),
	}

	model := NewManualVolume(persistence.NewMonthlyAggregationInMemoryProvider(nil))

	forecastRecords, err := model.Calculate(forecastRule)

	assert.Nil(t, err)

	assert.Equal(t, forecastRule.Period.GetTotalMonths(), int64(len(forecastRecords)))

	expectedVolume := 25.0

	for _, r := range forecastRecords {
		t.Run(
			r.Month.ToDateString(), func(t *testing.T) {
				assert.Equal(t, expectedVolume, r.VolumeActual)
			},
		)
	}
}

func Test_when_there_is_historical_volume_for_full_previous_year(t *testing.T) {
	forecastRule := &models.ForecastRule{
		Period: types.Period{
			StartDate: types.ToDate("2024-09-01"),
			EndDate:   types.ToDate("2024-10-01"),
		},
		Volume: 1000.0,
		LHM:    types.ToDate("2023-10-01"),
	}

	maRecords := []models.MonthlyAggregationRecord{
		{Month: types.ToDate("2023-09-01"), VolumeActual: 170.636736},
		{Month: types.ToDate("2023-10-01"), VolumeActual: 464.041273},
	}

	model := NewManualVolume(persistence.NewMonthlyAggregationInMemoryProvider(maRecords))

	forecastRecords, err := model.Calculate(forecastRule)

	assert.Nil(t, err)

	expectedForecastRecords := []dto.ForecastRecord{
		{Month: types.ToDate("2024-09-01"), VolumeActual: 268.855599},
		{Month: types.ToDate("2024-10-01"), VolumeActual: 731.144401},
	}

	expectedForecastRecordsMap := forecastRecordsToMap(expectedForecastRecords)

	for _, resultRecord := range forecastRecords {
		t.Run(
			resultRecord.Month.ToDateString(), func(t *testing.T) {
				expectedRecord := expectedForecastRecordsMap[resultRecord.Month.ToDateString()]

				assert.InDelta(t, expectedRecord.VolumeActual, resultRecord.VolumeActual, 0.000001)
			},
		)

	}
}

func Test_when_volume_is_zero_for_previous_year(t *testing.T) {
	forecastRule := &models.ForecastRule{
		Period: types.Period{
			StartDate: types.ToDate("2024-09-01"),
			EndDate:   types.ToDate("2024-10-01"),
		},
		Volume: 10000.0,
		LHM:    types.ToDate("2024-08-01"),
	}

	maRecords := []models.MonthlyAggregationRecord{
		{Month: types.ToDate("2023-09-01"), VolumeActual: 0.0},
		{Month: types.ToDate("2023-10-01"), VolumeActual: 0.0},
	}

	model := NewManualVolume(persistence.NewMonthlyAggregationInMemoryProvider(maRecords))

	forecastRecords, err := model.Calculate(forecastRule)

	assert.Nil(t, err)

	assert.Equal(t, forecastRule.Period.GetTotalMonths(), int64(len(forecastRecords)))

	expectedForecastRecords := []dto.ForecastRecord{
		{Month: types.ToDate("2024-09-01"), VolumeActual: 5000.0},
		{Month: types.ToDate("2024-10-01"), VolumeActual: 5000.0},
	}

	assert.Equal(t, expectedForecastRecords, forecastRecords)
}

func Test_when_no_volume_in_several_months_of_previous_year(t *testing.T) {
	forecastRule := &models.ForecastRule{
		Period: types.Period{
			StartDate: types.ToDate("2022-07-01"),
			EndDate:   types.ToDate("2023-06-01"),
		},
		Volume: 10000.0,
		LHM:    types.ToDate("2022-04-01"),
	}

	maRecords := []models.MonthlyAggregationRecord{
		{Month: types.ToDate("2021-07-01"), VolumeActual: 217.0},
		{Month: types.ToDate("2021-08-01"), VolumeActual: 137.0},
		// {Month: toDate("2021-09-01"), VolumeActual: 0.0},
		{Month: types.ToDate("2021-10-01"), VolumeActual: 253.0},
		{Month: types.ToDate("2021-11-01"), VolumeActual: 126.0},
		{Month: types.ToDate("2021-12-01"), VolumeActual: 432.0},
		{Month: types.ToDate("2022-01-01"), VolumeActual: 127.0},
		// {Month: toDate("2022-02-01"), VolumeActual: 0.0},
		{Month: types.ToDate("2022-03-01"), VolumeActual: 271.0},
		{Month: types.ToDate("2022-04-01"), VolumeActual: 267.0},
		// {Month: toDate("2022-05-01"), VolumeActual: 0.0},
		// {Month: toDate("2022-06-01"), VolumeActual: 0.0},
	}

	model := NewManualVolume(persistence.NewMonthlyAggregationInMemoryProvider(maRecords))

	forecastRecords, err := model.Calculate(forecastRule)

	assert.Nil(t, err)

	assert.Equal(t, forecastRule.Period.GetTotalMonths(), int64(len(forecastRecords)))

	expectedForecastRecords := []dto.ForecastRecord{
		{Month: types.ToDate("2022-07-01"), VolumeActual: 1185.792349},
		{Month: types.ToDate("2022-08-01"), VolumeActual: 748.633879},
		{Month: types.ToDate("2022-10-01"), VolumeActual: 1382.513661},
		{Month: types.ToDate("2022-11-01"), VolumeActual: 688.524590},
		{Month: types.ToDate("2022-12-01"), VolumeActual: 2360.655737},
		{Month: types.ToDate("2023-01-01"), VolumeActual: 693.989071},
		{Month: types.ToDate("2023-03-01"), VolumeActual: 1480.874316},
		{Month: types.ToDate("2023-04-01"), VolumeActual: 1459.016393},
	}

	expectedForecastRecordsMap := forecastRecordsToMap(expectedForecastRecords)

	for _, resultRecord := range forecastRecords {
		t.Run(
			resultRecord.Month.ToDateString(), func(t *testing.T) {
				expectedRecord := expectedForecastRecordsMap[resultRecord.Month.ToDateString()]

				assert.InDelta(t, expectedRecord.VolumeActual, resultRecord.VolumeActual, 0.000001)
			},
		)
	}
}

func Test_when_volume_is_zero_in_several_months_of_previous_year(t *testing.T) {
	forecastRule := &models.ForecastRule{
		Period: types.Period{
			StartDate: types.ToDate("2022-07-01"),
			EndDate:   types.ToDate("2023-06-01"),
		},
		Volume: 10000.0,
		LHM:    types.ToDate("2022-06-01"),
	}

	maRecords := []models.MonthlyAggregationRecord{
		{Month: types.ToDate("2021-07-01"), VolumeActual: 217.0},
		{Month: types.ToDate("2021-08-01"), VolumeActual: 137.0},
		{Month: types.ToDate("2021-09-01"), VolumeActual: 0.0},
		{Month: types.ToDate("2021-10-01"), VolumeActual: 253.0},
		{Month: types.ToDate("2021-11-01"), VolumeActual: 126.0},
		{Month: types.ToDate("2021-12-01"), VolumeActual: 432.0},
		{Month: types.ToDate("2022-01-01"), VolumeActual: 127.0},
		{Month: types.ToDate("2022-02-01"), VolumeActual: 0.0},
		{Month: types.ToDate("2022-03-01"), VolumeActual: 271.0},
		{Month: types.ToDate("2022-04-01"), VolumeActual: 267.0},
		{Month: types.ToDate("2022-05-01"), VolumeActual: 0.0},
		{Month: types.ToDate("2022-06-01"), VolumeActual: 0.0},
	}

	model := NewManualVolume(persistence.NewMonthlyAggregationInMemoryProvider(maRecords))

	forecastRecords, err := model.Calculate(forecastRule)

	assert.Nil(t, err)

	assert.Equal(t, forecastRule.Period.GetTotalMonths(), int64(len(forecastRecords)))

	expectedForecastRecords := []dto.ForecastRecord{
		{Month: types.ToDate("2022-07-01"), VolumeActual: 1185.792349},
		{Month: types.ToDate("2022-08-01"), VolumeActual: 748.633879},
		{Month: types.ToDate("2022-09-01"), VolumeActual: 0},
		{Month: types.ToDate("2022-10-01"), VolumeActual: 1382.513661},
		{Month: types.ToDate("2022-11-01"), VolumeActual: 688.524590},
		{Month: types.ToDate("2022-12-01"), VolumeActual: 2360.655737},
		{Month: types.ToDate("2023-01-01"), VolumeActual: 693.989071},
		{Month: types.ToDate("2023-02-01"), VolumeActual: 0},
		{Month: types.ToDate("2023-03-01"), VolumeActual: 1480.874316},
		{Month: types.ToDate("2023-04-01"), VolumeActual: 1459.016393},
		{Month: types.ToDate("2023-05-01"), VolumeActual: 0},
		{Month: types.ToDate("2023-06-01"), VolumeActual: 0},
	}

	expectedForecastRecordsMap := forecastRecordsToMap(expectedForecastRecords)

	for _, resultRecord := range forecastRecords {
		t.Run(
			resultRecord.Month.ToDateString(), func(t *testing.T) {
				expectedRecord := expectedForecastRecordsMap[resultRecord.Month.ToDateString()]

				assert.InDelta(t, expectedRecord.VolumeActual, resultRecord.VolumeActual, 0.000001)
			},
		)
	}
}

func forecastRecordsToMap(records []dto.ForecastRecord) map[string]dto.ForecastRecord {
	expectedForecastRecordsMap := make(map[string]dto.ForecastRecord, len(records))

	for _, r := range records {
		expectedForecastRecordsMap[r.Month.ToDateString()] = r
	}

	return expectedForecastRecordsMap
}
