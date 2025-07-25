package web

import (
	"github.com/gin-gonic/gin"
	"log/slog"
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
	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	return &Server{
		address: cfg.Address,
		router:  newRouter(log, db),
	}
}

func newRouter(log *slog.Logger, db *postgres.Storage) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.Logger(log))

	{
		v1 := router.Group("/url")
		v1.POST("", handler.Save(log, db))
		{
			v12 := v1.Group("/delete")
			v12.Use(gin.BasicAuth(gin.Accounts{
				"bob": "123",
			}))
			v12.DELETE("", handler.Delete(log, db))
		}
	}

	router.GET("/:alias", handler.Redirect(log, db))

	return router
}

func (s *Server) Start() error {
	return s.router.Run(s.address)
}
