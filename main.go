package main

import (
	"remy/internal/config"
	"remy/internal/infrastructure"
	"remy/internal/logging"
)

func main() {
	config.LoadEnv()

	cfg, err := config.LoadConfig()
	if err != nil {
		logging.Logger.Fatalf("failed to load configuration: %v", err)
	}

	logging.InitLogger(cfg)

	db, err := config.SetupDatabase(cfg)
	if err != nil {
		logging.Logger.Fatalf("failed to connect to database: %v", err)
	}

	router := infrastructure.NewRouter()
	infrastructure.SetupRoutes(router, db)

	server := infrastructure.NewServer(router, cfg)
	infrastructure.StartServer(server)
}
