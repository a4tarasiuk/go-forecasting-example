package main

import (
	"fmt"

	"forecasting/internal/config"
	"forecasting/internal/infra"
)

func main() {
	cfg, _ := config.InitConfig()

	_infra, _ := infra.NewInfra(cfg)

	defer _infra.Shutdown()

	fmt.Println(_infra.GetDB())
}
