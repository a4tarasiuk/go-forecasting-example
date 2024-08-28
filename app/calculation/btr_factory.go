package calculation

import (
	"forecasting/app/domain/models"
	"forecasting/app/providers"
)

const btrChunkSize = 5000

type BudgetTrafficFactory struct {
	records []models.BudgetTrafficRecord

	budgetTrafficProvider providers.BudgetTrafficProvider
}

func NewBudgetTrafficFactory(budgetTrafficProvider providers.BudgetTrafficProvider) *BudgetTrafficFactory {
	return &BudgetTrafficFactory{budgetTrafficProvider: budgetTrafficProvider}
}

func (f *BudgetTrafficFactory) AddMany(records []models.BudgetTrafficRecord) {
	if len(records) == 0 {
		return
	}

	if len(f.records) > 0 {
		lastRecordBeforeInsert, firstCommitingRecord := f.records[len(f.records)-1], records[0]

		recordsFromDifferentYears := lastRecordBeforeInsert.Month.Year() != firstCommitingRecord.Month.Year()

		if len(f.records) > btrChunkSize || recordsFromDifferentYears {
			f.Commit()
		}
	}

	f.records = append(f.records, records...)
}

func (f *BudgetTrafficFactory) Commit() {
	f.budgetTrafficProvider.CreateMany(f.records)

	f.records = make([]models.BudgetTrafficRecord, 0)
}
