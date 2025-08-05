package handler

import (
	"testing"
	"url-shortener/internal/validation"
)

func TestMain(m *testing.M) {
	validation.Init()
	m.Run()
}
