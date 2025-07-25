package db

import "errors"

var (
	ErrURLAlreadyExists = errors.New("URL already exists")
	ErrURLNotFound      = errors.New("URL not found")
)
