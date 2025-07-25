package logger

import (
	"fmt"
	"log/slog"
	"os"
	"url-shortener/internal/config"
)

func New(cfg *config.Config) *slog.Logger {
	switch cfg.Env {
	case "local":
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		panic(fmt.Sprintf("invalid env from cfg: %s", cfg.Env))
	}
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
