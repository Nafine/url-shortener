package main

import (
	"fmt"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/db/postgres"
	"url-shortener/internal/logger"
	"url-shortener/internal/validation"
	"url-shortener/internal/web"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		fmt.Println("Error loading config: ", err)
		os.Exit(1)
	}

	log := logger.New(cfg)
	log.Info("logger initialized")
	log.Debug("debug messages enabled")

	db, err := postgres.New(cfg.StorageDSN)
	if err != nil {
		log.Error("error initializing database", logger.Err(err))
		os.Exit(1)
	}
	log.Info("database initialized")

	validation.Init()
	server := web.New(cfg, log, db)

	log.Info("starting server")
	if err := server.Start(); err != nil {
		log.Info("server shutdown", logger.Err(err))
	}
}
