package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nafine/url-shortener/internal/auth"
	"github.com/nafine/url-shortener/internal/config"
	"github.com/nafine/url-shortener/internal/db/postgres"
	"github.com/nafine/url-shortener/internal/web/handler"
	"github.com/nafine/url-shortener/internal/web/middleware"
	"log/slog"
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
		address: fmt.Sprintf("%v:%v", cfg.Host, cfg.Port),
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
