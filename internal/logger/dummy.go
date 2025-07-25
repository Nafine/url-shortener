package logger

import (
	"context"
	"log/slog"
)

type DummyHandler struct{}

func (l *DummyHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (l *DummyHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (l *DummyHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return l
}

func (l *DummyHandler) WithGroup(_ string) slog.Handler {
	return l
}

func NewDummyLogger() *slog.Logger {
	return slog.New(&DummyHandler{})
}
