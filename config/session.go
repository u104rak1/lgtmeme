package config

import (
	"os"

	"github.com/boj/redistore"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var Store *redistore.RediStore

func initSessionStore() {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	address := host + ":" + port
	secretKey := os.Getenv("SESSION_SECRET_KEY")

	var err error
	Store, err = redistore.NewRediStore(10, "tcp", address, "", []byte(secretKey))
	if err != nil {
		panic(err)
	}
}

func SessionMiddleware() echo.MiddlewareFunc {
	return session.Middleware(Store)
}
