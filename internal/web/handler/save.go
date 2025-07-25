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
	"url-shortener/internal/pkg/random"
	"url-shortener/internal/web/api"
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

var (
	ErrURLExists   = api.Error("URL already exists")
	ErrFailSaveURL = api.Error("failed to save url")
)

var (
	InfoUrlSaveErr = "failed to save url"
)

func Save(log *slog.Logger, urlSaver URLSaver) gin.HandlerFunc {
	const op = "handlers.Save"
	return func(c *gin.Context) {
		log := log.With("operation", op)

		var req SaveRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Info(InfoInvalidRequest, logger.Err(err))
			c.JSON(http.StatusBadRequest, api.ErrInvalidRequest)
			return
		}

		if err := validator.New().Struct(req); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			log.Info(InfoInvalidRequestFields, logger.Err(err))
			c.JSON(http.StatusBadRequest, api.ValidationError(validationErrors))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomPath(rand.Intn(12) + 4)
		}

		if err := urlSaver.SaveURL(req.Url, alias); err != nil {
			if errors.Is(err, db.ErrURLAlreadyExists) {
				log.Info(InfoURLExists, logger.Err(err))
				c.JSON(http.StatusOK, ErrURLExists)
				return
			}
			log.Info(InfoUrlSaveErr, logger.Err(err))
			c.JSON(http.StatusInternalServerError, ErrFailSaveURL)
			return
		}

		log.Info("successfully saved url", slog.String("url", req.Url), slog.String("alias", alias))
		c.JSON(http.StatusOK, Response{
			StatusResponse: api.Ok(),
			Alias:          alias,
		})
	}
}
