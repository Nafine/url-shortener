package web

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"url-shortener/internal/auth"
	"url-shortener/internal/config"
	"url-shortener/internal/db/postgres"
	"url-shortener/internal/web/handler"
	"url-shortener/internal/web/middleware"
)

type Server struct {
	address string
	router  *gin.Engine
}

func New(cfg *config.Config, log *slog.Logger, db *postgres.Storage) *Server {
	if cfg.AppEnv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	return &Server{
		address: cfg.Address,
		router:  newRouter(log, db, cfg),
	}
}

func newRouter(log *slog.Logger, db *postgres.Storage, cfg *config.Config) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.Logger(log))

	{
		root := router.Group("/url")
		root.POST("", middleware.JSONParser[handler.SaveRequest](log), handler.Save(log, db))
		{
			authoritative := root.Group("/delete")
			authoritative.Use(middleware.APIKeyAuth(log, auth.LoadKeys()))
			authoritative.DELETE("", middleware.JSONParser[handler.DeleteRequest](log), handler.Delete(log, db))
		}
	}

	router.GET("/:alias", middleware.SegmentParser[handler.Alias](log), handler.Redirect(log, db))

	return router
}

func (s *Server) Start() error {
	return s.router.Run(s.address)
}
