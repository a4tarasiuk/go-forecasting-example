package types

import (
	"github.com/golang-module/carbon/v2"
)

type Period struct {
	StartDate carbon.Date
	EndDate   carbon.Date
}

func NewPeriod(StartDate carbon.Date, EndDate carbon.Date) Period {
	return Period{
		StartDate: StartDate.StartOfMonth().ToDateStruct(),
		EndDate:   EndDate.EndOfMonth().ToDateStruct(),
	}
}

func (p *Period) GetTotalMonths() int64 {
	diffInMonths := p.StartDate.DiffInMonths(p.EndDate.Carbon)

	return diffInMonths + 1
}

func (p *Period) GetMonths() []carbon.Date {
	month := p.StartDate

	totalMonths := p.GetTotalMonths()

	months := make([]carbon.Date, totalMonths)

	for idx := range totalMonths {
		months[idx] = month

		month = month.Carbon.AddMonth().StartOfMonth().ToDateStruct()
	}

	return months
}

func (p *Period) Contains(month carbon.Date) bool {
	months := make(map[string]string, p.GetTotalMonths())

	for _, _month := range p.GetMonths() {
		months[_month.ToDateString()] = _month.ToDateString()
	}

	_, exists := months[month.StartOfMonth().ToDateString()]

	return exists
}

func (p *Period) MoreThenAYear() bool {
	return p.GetTotalMonths() > carbon.MonthsPerYear
}

func (p *Period) CutToOneYear() Period {
	if !p.MoreThenAYear() {
		return *p
	}

	cutPeriod := NewPeriod(
		p.StartDate,
		p.StartDate.AddMonths(carbon.MonthsPerYear-1).ToDateStruct(),
	)

	return cutPeriod
}

func (p *Period) PreviousYear() Period {
	previousYearPeriod := NewPeriod(
		p.StartDate.SubYear().ToDateStruct(),
		p.EndDate.SubYear().ToDateStruct(),
	)

	return previousYearPeriod
}

func ToDate(dtStr string) carbon.Date {
	return carbon.Parse(dtStr).ToDateStruct()
}
