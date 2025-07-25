package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"url-shortener/internal/db"
	"url-shortener/internal/logger"
	"url-shortener/internal/web/api"
)

type DeleteRequest struct {
	Alias string `json:"alias" validate:"required"`
}

type URLDeleter interface {
	DeleteURL(string) error
}

var (
	ErrDeletion    = api.Error("error deleting url")
	ErrURLNotExist = api.Error("URL does not exist")
)

var (
	InfoURLNotFound    = "url was not found"
	InfoURLDeletionErr = "url deletion error"
)

func Delete(log *slog.Logger, urlDeleter URLDeleter) gin.HandlerFunc {
	const op = "handler.Delete"

	return func(c *gin.Context) {
		log := log.With("operation", op)

		var req DeleteRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Info(InfoInvalidRequest, logger.Err(err))
			c.JSON(http.StatusBadRequest, api.ErrInvalidRequest)
			return
		}

		if req.Alias == "" {
			log.Info(InfoEmptyAlias)
			c.JSON(http.StatusBadRequest, ErrEmptyAlias)
			return
		}

		if err := urlDeleter.DeleteURL(req.Alias); err != nil {
			if errors.Is(err, db.ErrURLNotFound) {
				log.Info(InfoURLNotFound, logger.Err(err))
				c.JSON(http.StatusOK, ErrURLNotExist)
				return
			}
			log.Info(InfoURLDeletionErr, logger.Err(err))
			c.JSON(http.StatusInternalServerError, ErrDeletion)
			return
		}

		log.Info("deletion succeed", slog.String("alias", req.Alias))
		c.JSON(http.StatusOK, api.Ok())
	}
}
