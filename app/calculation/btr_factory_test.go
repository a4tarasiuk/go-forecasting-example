package calculation

import (
	"testing"

	"forecasting/app/domain/models"
	"forecasting/app/persistence"
	"forecasting/core/types"
	"github.com/stretchr/testify/assert"
)

func Test_AddMany_when_records_are_not_created_when_total_is_less_then_chunk_size(t *testing.T) {
	records := []models.BudgetTrafficRecord{
		{BudgetSnapshotID: 1, Month: types.ToDate("2024-01-01")},
		{BudgetSnapshotID: 2, Month: types.ToDate("2024-01-01")},
	}

	btrProvider := persistence.NewBudgetTrafficInMemoryProvider(nil)

	chunkSize := 5

	factory := NewBudgetTrafficFactory(btrProvider, &chunkSize)

	factory.AddMany(records)

	assert.Zero(t, btrProvider.Count())

	assert.Equal(t, len(records), len(factory.records))
}

func Test_AddMany_when_records_are_committed_after_total_is_more_then_chunk_size(t *testing.T) {
	records := []models.BudgetTrafficRecord{{BudgetSnapshotID: 1}, {BudgetSnapshotID: 2}, {BudgetSnapshotID: 3}}

	btrProvider := persistence.NewBudgetTrafficInMemoryProvider(nil)

	chunkSize := 2

	factory := NewBudgetTrafficFactory(btrProvider, &chunkSize)

	factory.AddMany(records)

	assert.Zero(t, len(factory.records))

	assert.Equal(t, int64(len(records)), btrProvider.Count())
}

func Test_AddMany_when_passed_records_are_nil(t *testing.T) {
	btrProvider := persistence.NewBudgetTrafficInMemoryProvider(nil)

	factory := NewBudgetTrafficFactory(btrProvider, nil)

	factory.AddMany(nil)

	assert.Zero(t, len(factory.records))

	assert.Zero(t, btrProvider.Count())
}

func Test_Commit(t *testing.T) {
	records := []models.BudgetTrafficRecord{{BudgetSnapshotID: 1}, {BudgetSnapshotID: 2}}

	btrProvider := persistence.NewBudgetTrafficInMemoryProvider(nil)

	factory := NewBudgetTrafficFactory(btrProvider, nil)

	factory.records = records

	factory.Commit()

	assert.Zero(t, len(factory.records))

	assert.Equal(t, int64(len(records)), btrProvider.Count())
}
