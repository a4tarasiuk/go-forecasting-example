package persistence

import (
	"forecasting/app/domain/models"
	"forecasting/app/providers"
)

type BudgetTrafficInMemoryProvider struct {
	records []models.BudgetTrafficRecord
}

func NewBudgetTrafficInMemoryProvider(records []models.BudgetTrafficRecord) *BudgetTrafficInMemoryProvider {
	return &BudgetTrafficInMemoryProvider{records: records}
}

func (b *BudgetTrafficInMemoryProvider) Get(options providers.BudgetTrafficOptions) []models.BudgetTrafficRecord {
	return b.records
}

func (b *BudgetTrafficInMemoryProvider) CreateMany(records []models.BudgetTrafficRecord) {
	b.records = append(b.records, records...)
}

func (b *BudgetTrafficInMemoryProvider) ClearForecasted() {
	filteredRecords := make([]models.BudgetTrafficRecord, 0)

	for _, record := range b.records {
		if record.TrafficType != 2 { // 2 - FORECASTED
			filteredRecords = append(filteredRecords, record)
		}
	}

	b.records = filteredRecords
}

func (b *BudgetTrafficInMemoryProvider) CountForecasted() int64 {
	total := int64(0)

	for _, r := range b.records {
		if r.TrafficType == 2 {
			total++
		}
	}

	return total
}
