package main

import (
	"github.com/ChaPerx64/dobby/apps/backend/internal/adapters/api"
	"github.com/ChaPerx64/dobby/apps/backend/internal/config"
)

func main() {
	cfg := config.Load()
	api.RunServer(cfg)
}
