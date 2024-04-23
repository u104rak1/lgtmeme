package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
)

func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		config.Logger.Info("Received request", "method", c.Request().Method, "uri", c.Request().RequestURI)
		return next(c)
	}
}
