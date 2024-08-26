package infra

import (
	"context"
	"database/sql"

	"forecasting/internal/config"
)

type Infrastructure interface {
	GetConfig() config.AppConfig

	GetDB() *sql.DB

	Shutdown()
}

type Module interface {
	Startup(context.Context, Infrastructure) error
}
