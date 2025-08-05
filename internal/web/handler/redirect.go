package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"url-shortener/internal/db"
	"url-shortener/internal/logger"
	"url-shortener/internal/web/api"
	"url-shortener/internal/web/middleware"
)

type URLGetter interface {
	GetURL(string) (string, error)
}

type Alias struct {
	Alias string `uri:"alias" binding:"required" validate:"required,urlsegment"`
}

func Redirect(log *slog.Logger, urlGetter URLGetter) gin.HandlerFunc {
	const op = "handler.Redirect"

	return func(c *gin.Context) {
		log := log.With("operation", op)

		segment, ok := c.Get(middleware.SegmentBody)
		if !ok {
			log.Info(api.InfoMissingSegment)
			c.JSON(http.StatusBadRequest, api.ErrInvalidRequest)
			return
		}

		alias := segment.(Alias)

		url, err := urlGetter.GetURL(alias.Alias)

		if err != nil {
			if errors.Is(err, db.ErrURLNotFound) {
				log.Info(api.InfoURLNotFound, logger.Err(err))
				c.JSON(http.StatusNotFound, api.ErrNotFound)
				return
			}
			log.Info(api.InfoRedirectErr, logger.Err(err))
			c.JSON(http.StatusInternalServerError, api.ErrRedirect)
			return
		}

		log.Info("redirection succeed", slog.String("alias", alias.Alias), slog.String("url", url))
		c.Redirect(http.StatusFound, url)
	}
}
