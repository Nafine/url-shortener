package app

import (
	"github.com/gin-gonic/gin"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/db/postgres"
	"url-shortener/internal/logger"
	"url-shortener/internal/web"
)

type App struct {
	router *gin.Engine
}

func Run() {
	cfg := config.Get()

	log := logger.New(cfg)

	log.Info("logger initialized")
	log.Debug("debug messages enabled")

	db, err := postgres.New(cfg.Storage)
	if err != nil {
		log.Error("error initializing database", logger.Err(err))
		os.Exit(1)
	}
	log.Info("database initialized")

	server := web.New(cfg, log, db)

	log.Info("starting server")
	if err := server.Start(); err != nil {
		log.Info("server shutdown", logger.Err(err))
	}
}
