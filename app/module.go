package app

import (
	"context"

	"forecasting/internal/infra"
)

type Module struct {
}

func (m *Module) Startup(context.Context, infra.Infrastructure) error {
	return nil
}
