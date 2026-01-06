package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"url-shortener/internal/logger"
	"url-shortener/internal/validation"
	"url-shortener/internal/web/api"
)

const RequestBody = "requestBody"
const SegmentBody = "segmentBody"

func JSONParser[T any](log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := parseJSON[T](api.Context{log, c})
		if err == nil {
			c.Set(RequestBody, req)
			c.Next()
		}
	}
}

func parseJSON[T any](parseCtx api.Context) (T, error) {
	var req T
	if err := parseCtx.Ctx.ShouldBindJSON(&req); err != nil {
		abortInvalidRequest(parseCtx, err)
		return req, err
	} else if err := validation.Validate.Struct(req); err != nil {
		abortInvalidRequestFields(parseCtx, err)
		return req, err
	}
	return req, nil
}

func SegmentParser[T any](log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		segment, err := parseUri[T](api.Context{log, c})
		if err == nil {
			c.Set(SegmentBody, segment)
			c.Next()
		}
	}
}

func parseUri[T any](parseCtx api.Context) (T, error) {
	var uriSegment T
	if err := parseCtx.Ctx.ShouldBindUri(&uriSegment); err != nil {
		abortInvalidRequest(parseCtx, err)
		return uriSegment, err
	} else if err := validation.Validate.Struct(uriSegment); err != nil {
		abortInvalidRequestFields(parseCtx, err)
		return uriSegment, err
	}
	return uriSegment, nil
}

func abortInvalidRequest(parseCtx api.Context, err error) {
	parseCtx.Log.Info(api.InfoInvalidRequest, logger.Err(err))
	writeAbort(parseCtx.Ctx, api.ErrInvalidRequest)
}

func abortInvalidRequestFields(parseCtx api.Context, err error) {
	validationErrors := err.(validator.ValidationErrors)
	parseCtx.Log.Info(api.InfoInvalidRequestFields, logger.Err(err))
	writeAbort(parseCtx.Ctx, api.ValidationError(validationErrors))
}

func writeAbort(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusBadRequest, data)
	ctx.Abort()
}
