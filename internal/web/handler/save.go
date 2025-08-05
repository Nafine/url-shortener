package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"math/rand"
	"net/http"
	"url-shortener/internal/db"
	"url-shortener/internal/logger"
	"url-shortener/internal/random"
	"url-shortener/internal/web/api"
	"url-shortener/internal/web/middleware"
)

type SaveRequest struct {
	Url   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	api.StatusResponse
	Alias string `json:"alias"`
}

type URLSaver interface {
	SaveURL(string, string) error
}

func Save(log *slog.Logger, urlSaver URLSaver) gin.HandlerFunc {
	const op = "handlers.Save"
	return func(c *gin.Context) {
		log := log.With("operation", op)

		req, ok := c.Get(middleware.RequestBody)
		if !ok {
			log.Info(api.InfoMissingRequestBody)
			c.JSON(http.StatusBadRequest, api.ErrInvalidRequest)
			return
		}

		saveRequest := req.(SaveRequest)

		if err := validator.New().Struct(req); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			log.Info(api.InfoInvalidRequestFields, logger.Err(err))
			c.JSON(http.StatusBadRequest, api.ValidationError(validationErrors))
			return
		}

		alias := saveRequest.Alias
		if alias == "" {
			alias = random.NewRandomPath(rand.Intn(12) + 4)
		}

		if err := urlSaver.SaveURL(saveRequest.Url, alias); err != nil {
			if errors.Is(err, db.ErrURLAlreadyExists) {
				log.Info(api.InfoURLExists, logger.Err(err))
				c.JSON(http.StatusOK, api.ErrURLExists)
				return
			}
			log.Info(api.InfoURLSaveErr, logger.Err(err))
			c.JSON(http.StatusInternalServerError, api.ErrSave)
			return
		}

		log.Info("successfully saved url", slog.String("url", saveRequest.Url), slog.String("alias", alias))
		c.JSON(http.StatusOK, Response{
			StatusResponse: api.Ok(),
			Alias:          alias,
		})
	}
}
