package handler

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"url-shortener/internal/logger"
	"url-shortener/internal/web/api"
)

type URLGetter interface {
	GetURL(string) (string, error)
}

type Alias struct {
	Alias string `uri:"alias" binding:"required" validate:"required,url"`
}

var (
	ErrEmptyAlias   = api.Error("empty alias")
	ErrInvalidAlias = api.Error("invalid alias")
	ErrNotFound     = api.Error("url not found")
)

func Redirect(log *slog.Logger, urlGetter URLGetter) gin.HandlerFunc {
	const op = "handler.Redirect"

	return func(c *gin.Context) {
		log := log.With("operation", op)

		var alias Alias

		if err := c.ShouldBindUri(&alias); err != nil {
			log.Info(InfoInvalidAlias, logger.Err(err))
			c.JSON(http.StatusBadRequest, ErrInvalidAlias)
			return
		}

		url, err := urlGetter.GetURL(alias.Alias)

		if err != nil {
			log.Info(InfoURLNotFound, logger.Err(err))
			c.JSON(http.StatusNotFound, ErrNotFound)
			return
		}

		log.Info("redirection succeed", slog.String("alias", alias.Alias), slog.String("url", url))
		c.Redirect(http.StatusFound, url)
	}
}
