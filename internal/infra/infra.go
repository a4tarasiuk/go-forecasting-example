package infra

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"forecasting/internal/config"
)

type Infra struct {
	cfg config.AppConfig

	db *sql.DB
}

func NewInfra(cfg config.AppConfig) (*Infra, error) {
	infra := &Infra{cfg: cfg}

	if err := infra.initDB(); err != nil {
		return nil, err
	}

	return infra, nil
}

func (i *Infra) initDB() (err error) {
	dbCfg := i.cfg.DB

	connStr := fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?sslmode=disable",
		dbCfg.Schema,
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Name,
	)

	i.db, err = sql.Open("postgres", connStr)

	return err
}

func (i *Infra) GetConfig() config.AppConfig {
	return i.cfg
}

func (i *Infra) GetDB() *sql.DB {
	return i.db
}

func (i *Infra) Shutdown() (err error) {
	err = i.db.Close()
	return err
}
