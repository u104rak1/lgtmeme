package config

import (
	"io"
	"os"

	"github.com/labstack/echo/v4"
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

func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		Logger.Info("Received request", "method", c.Request().Method, "uri", c.Request().RequestURI)
		return next(c)
	}
}
