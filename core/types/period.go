package types

import (
	"github.com/golang-module/carbon/v2"
)

type Period struct {
	StartDate carbon.Date
	EndDate   carbon.Date
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

		month = month.Carbon.AddMonth().ToDateStruct()
	}

	return months
}
