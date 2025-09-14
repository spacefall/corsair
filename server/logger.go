package server

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/bastienwirtz/corsair/config"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

func SetupLogger(cfg config.LoggingConfig, version string) error {
	var level slog.Level
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		return fmt.Errorf("invalid log level: %s", cfg.Level)
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{Level: level}

	switch strings.ToLower(cfg.Format) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	case "pretty":
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
			NoColor:    !isatty.IsTerminal(os.Stdout.Fd()),
		})
	default:
		return fmt.Errorf("invalid log format: %s", cfg.Format)
	}

	logger := slog.New(handler).With("version", version)
	slog.SetDefault(logger)

	return nil
}
