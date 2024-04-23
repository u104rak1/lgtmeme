package middleware

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
)

func SessionMiddleware() echo.MiddlewareFunc {
	return session.Middleware(config.Store)
}
