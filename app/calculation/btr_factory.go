package calculation

import (
	"forecasting/app/domain/models"
	"forecasting/app/providers"
)

type BudgetTrafficFactory struct {
	records []models.BudgetTrafficRecord

	budgetTrafficProvider providers.BudgetTrafficProvider

	chunkSize int
}

func NewBudgetTrafficFactory(
	budgetTrafficProvider providers.BudgetTrafficProvider,
	chunkSize *int,
) *BudgetTrafficFactory {
	defaultChunkSize := 5000

	if chunkSize == nil {
		chunkSize = &defaultChunkSize
	}

	return &BudgetTrafficFactory{budgetTrafficProvider: budgetTrafficProvider, chunkSize: *chunkSize}
}

func (f *BudgetTrafficFactory) AddMany(records []models.BudgetTrafficRecord) {
	if len(records) == 0 {
		return
	}

	f.records = append(f.records, records...)

	if len(f.records) > f.chunkSize {
		f.Commit()
	}
}

func (f *BudgetTrafficFactory) Commit() {
	f.budgetTrafficProvider.CreateMany(f.records)

	f.records = make([]models.BudgetTrafficRecord, 0)
}
