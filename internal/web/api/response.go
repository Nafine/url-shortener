package api

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

var (
	ErrInvalidRequest = Error("invalid request")
)

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
		default:
			errorMsgs = append(errorMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}
	return Error(strings.Join(errorMsgs, "\n"))
}
