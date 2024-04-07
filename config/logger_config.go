package config

import (
	"io"
	"os"

	"golang.org/x/exp/slog"
)

var Logger *slog.Logger

func NewLogger() {
	logLevel := os.Getenv("LOG_LEVEL")
	var handler slog.Handler

	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)

	if logLevel == "SILENT" {
		handler = slog.NewJSONHandler(io.Discard, nil)
	} else {
		handler = jsonHandler
	}

	Logger = slog.New(handler)
}
