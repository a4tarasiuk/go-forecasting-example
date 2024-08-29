package types

import (
	"strings"
	"testing"

	"github.com/golang-module/carbon/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewPeriod(t *testing.T) {
	startDate := ToDate("2024-08-01")
	endDate := ToDate("2024-10-01")

	period := NewPeriod(startDate, endDate)

	assert.Equal(t, startDate, period.StartDate)
	assert.Equal(t, endDate.EndOfMonth().ToDateStruct(), period.EndDate)
}

func TestPeriod_GetTotalMonths(t *testing.T) {
	params := []struct {
		startDate              string
		endDate                string
		expectedNumberOfMonths int64
	}{
		{startDate: "2024-01-01", endDate: "2024-06-01", expectedNumberOfMonths: 6},
		{startDate: "2024-06-01", endDate: "2024-12-01", expectedNumberOfMonths: 7},
		{startDate: "2024-01-01", endDate: "2024-12-01", expectedNumberOfMonths: 12},
		{startDate: "2024-06-01", endDate: "2025-06-01", expectedNumberOfMonths: 13},
		{startDate: "2024-07-01", endDate: "2025-12-01", expectedNumberOfMonths: 18},
		{startDate: "2024-11-01", endDate: "2025-02-01", expectedNumberOfMonths: 4},
	}

	for _, tt := range params {
		t.Run(
			strings.Join([]string{tt.startDate, tt.endDate}, "-"), func(t *testing.T) {
				period := NewPeriod(ToDate(tt.startDate), ToDate(tt.endDate))

				assert.Equal(t, tt.expectedNumberOfMonths, period.GetTotalMonths())
			},
		)
	}
}

func TestPeriod_MoreThenAYear(t *testing.T) {
	params := []struct {
		startDate     string
		endDate       string
		moreThanAYear bool
	}{
		{startDate: "2024-01-01", endDate: "2024-01-01", moreThanAYear: false},
		{startDate: "2024-01-01", endDate: "2024-02-01", moreThanAYear: false},
		{startDate: "2024-01-01", endDate: "2024-06-01", moreThanAYear: false},
		{startDate: "2024-01-01", endDate: "2024-11-01", moreThanAYear: false},
		{startDate: "2024-01-01", endDate: "2024-12-01", moreThanAYear: false},
		{startDate: "2024-01-01", endDate: "2025-01-01", moreThanAYear: true},
		{startDate: "2024-01-01", endDate: "2025-02-01", moreThanAYear: true},
		{startDate: "2024-01-01", endDate: "2025-12-01", moreThanAYear: true},
		{startDate: "2024-01-01", endDate: "2026-01-01", moreThanAYear: true},
		{startDate: "2024-01-01", endDate: "2026-04-01", moreThanAYear: true},
		{startDate: "2024-01-01", endDate: "2026-12-01", moreThanAYear: true},
	}

	for _, tt := range params {
		t.Run(
			strings.Join([]string{tt.startDate, tt.endDate}, "-"), func(t *testing.T) {
				period := NewPeriod(ToDate(tt.startDate), ToDate(tt.endDate))

				assert.Equal(t, tt.moreThanAYear, period.MoreThenAYear())
			},
		)
	}
}

func TestPeriod_Contains(t *testing.T) {
	params := []struct {
		startDate   string
		endDate     string
		checkedDate string
		contains    bool
	}{
		{startDate: "2024-03-01", endDate: "2024-07-01", checkedDate: "2024-01-01", contains: false},
		{startDate: "2024-03-01", endDate: "2024-07-01", checkedDate: "2024-02-01", contains: false},
		{startDate: "2024-03-01", endDate: "2024-07-01", checkedDate: "2024-08-01", contains: false},
		{startDate: "2024-03-01", endDate: "2024-07-01", checkedDate: "2024-09-01", contains: false},
		{startDate: "2024-03-01", endDate: "2024-07-01", checkedDate: "2024-12-01", contains: false},
		{startDate: "2024-03-01", endDate: "2024-07-01", checkedDate: "2024-03-01", contains: true},
		{startDate: "2024-03-01", endDate: "2024-07-01", checkedDate: "2024-07-01", contains: true},
		{startDate: "2024-03-01", endDate: "2024-07-01", checkedDate: "2024-04-01", contains: true},
		{startDate: "2024-03-01", endDate: "2024-07-01", checkedDate: "2024-06-01", contains: true},
		{startDate: "2024-03-01", endDate: "2024-07-01", checkedDate: "2024-05-01", contains: true},
	}

	for _, tt := range params {
		t.Run(
			strings.Join([]string{tt.startDate, tt.endDate}, "-"), func(t *testing.T) {
				period := NewPeriod(ToDate(tt.startDate), ToDate(tt.endDate))

				periodContainsMonth := period.Contains(ToDate(tt.checkedDate))

				assert.Equal(t, tt.contains, periodContainsMonth)
			},
		)
	}
}

func TestPeriod_CutToOneYear_When_Period_Is_Less_Then_A_Year(t *testing.T) {
	period := NewPeriod(ToDate("2024-08-01"), ToDate("2024-10-01"))

	cutPeriod := period.CutToOneYear()

	assert.Equal(t, period, cutPeriod)
}

func TestPeriod_CutToOneYear_When_Period_Is_More_Then_A_Year(t *testing.T) {
	period := NewPeriod(ToDate("2024-08-01"), ToDate("2025-10-01"))

	cutPeriod := period.CutToOneYear()

	expectedPeriod := NewPeriod(period.StartDate, ToDate("2025-07-01"))

	assert.Equal(t, expectedPeriod, cutPeriod)
}

func TestPeriod_PreviousYear(t *testing.T) {
	period := NewPeriod(ToDate("2024-08-01"), ToDate("2025-10-01"))

	expectedPeriod := NewPeriod(ToDate("2023-08-01"), ToDate("2024-10-01"))

	previousYearPeriod := period.PreviousYear()

	assert.Equal(t, expectedPeriod, previousYearPeriod)
}

func TestPeriod_GetMonths(t *testing.T) {
	period := NewPeriod(ToDate("2024-08-01"), ToDate("2024-12-01"))

	expectedMonths := []carbon.Date{
		ToDate("2024-08-01").ToDateStruct(),
		ToDate("2024-09-01").ToDateStruct(),
		ToDate("2024-10-01").ToDateStruct(),
		ToDate("2024-11-01").ToDateStruct(),
		ToDate("2024-12-01").ToDateStruct(),
	}

	assert.Equal(t, expectedMonths, period.GetMonths())
}
