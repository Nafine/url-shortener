package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nafine/url-shortener/internal/db"
	"github.com/nafine/url-shortener/internal/logger"
	"github.com/nafine/url-shortener/internal/web/api"
	"github.com/nafine/url-shortener/internal/web/middleware"
	"log/slog"
	"net/http"
)

type DeleteRequest struct {
	Alias string `json:"alias" validate:"required"`
}

type URLDeleter interface {
	DeleteURL(string) error
}

func Delete(log *slog.Logger, urlDeleter URLDeleter) gin.HandlerFunc {
	const op = "handler.Delete"

	return func(c *gin.Context) {
		log := log.With("operation", op)

		req, ok := c.Get(middleware.RequestBody)
		if !ok {
			log.Info(api.InfoMissingRequestBody)
			c.JSON(http.StatusBadRequest, api.ErrInvalidRequest)
			return
		}

		deleteReq := req.(DeleteRequest)

		if deleteReq.Alias == "" {
			log.Info(api.InfoEmptyAlias)
			c.JSON(http.StatusBadRequest, api.ErrEmptyAlias)
			return
		}

		if err := urlDeleter.DeleteURL(deleteReq.Alias); err != nil {
			if errors.Is(err, db.ErrURLNotFound) {
				log.Info(api.InfoURLNotFound, logger.Err(err))
				c.JSON(http.StatusOK, api.ErrURLNotExist)
				return
			}
			log.Info(api.InfoURLDeletionErr, logger.Err(err))
			c.JSON(http.StatusInternalServerError, api.ErrDeletion)
			return
		}

		log.Info("deletion succeed", slog.String("alias", deleteReq.Alias))
		c.JSON(http.StatusOK, api.Ok())
	}
}
