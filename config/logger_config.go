package config

import (
	"os"

	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
)

var Logger *slog.Logger

func InitLogger() {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		Logger.Info("Received request", "method", c.Request().Method, "uri", c.Request().RequestURI)
		return next(c)
	}
}
