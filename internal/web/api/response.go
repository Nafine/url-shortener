package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"strings"
)

var (
	ErrEmptyAlias     = Error("empty alias")
	ErrNotFound       = Error("url not found")
	ErrURLNotExist    = Error("URL does not exist")
	ErrURLExists      = Error("URL already exists")
	ErrSave           = Error("failed to save url")
	ErrDeletion       = Error("error deleting url")
	ErrRedirect       = Error("redirect failed")
	ErrInvalidRequest = Error("invalid request")
)

var (
	InfoEmptyAlias           = "alias is empty"
	InfoInvalidRequest       = "invalid request body"
	InfoInvalidRequestFields = "invalid request fields"
	InfoURLExists            = "url already exists"
	InfoURLNotFound          = "url was not found"
	InfoURLDeletionErr       = "url deletion error"
	InfoURLSaveErr           = "failed to save url"
	InfoRedirectErr          = "redirect failed"
	InfoMissingRequestBody   = "missing request body"
	InfoMissingSegment       = "missing uri segment"
)

type Context struct {
	Log *slog.Logger
	Ctx *gin.Context
}

type StatusResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func Ok() StatusResponse {
	return StatusResponse{
		Status: "ok",
	}
}

func Error(msg string) StatusResponse {
	return StatusResponse{
		Status: "error",
		Error:  msg,
	}
}

func ValidationError(errors validator.ValidationErrors) StatusResponse {
	var errorMsgs []string

	for _, err := range errors {
		switch err.ActualTag() {
		case "required":
			errorMsgs = append(errorMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errorMsgs = append(errorMsgs, fmt.Sprintf("field %s is an invalid URL", err.Field()))
		case "urlsegment":
			errorMsgs = append(errorMsgs, fmt.Sprintf("field %s is an invalid URL segment", err.Field()))
		default:
			errorMsgs = append(errorMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}
	return Error(strings.Join(errorMsgs, "\n"))
}
