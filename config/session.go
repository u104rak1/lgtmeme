package config

import (
	"fmt"
	"os"

	"github.com/boj/redistore"
	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var Store *redistore.RediStore

func InitSessionStore() {
	var err error
	secretKey := os.Getenv("SESSION_SECRET_KEY")

	if os.Getenv("ECHO_MODE") == "production" {
		redisURL := os.Getenv("REDIS_URL")

		Store, err = redistore.NewRediStoreWithPool(&redis.Pool{
			Dial: func() (redis.Conn, error) {
				return redis.DialURL(redisURL)
			},
		}, []byte(secretKey))
	} else {
		host := os.Getenv("REDIS_HOST")
		port := os.Getenv("REDIS_PORT")
		address := fmt.Sprintf("%s:%s", host, port)

		Store, err = redistore.NewRediStore(10, "tcp", address, "", []byte(secretKey))
	}

	if err != nil {
		panic(err)
	}
}

func SessionMiddleware() echo.MiddlewareFunc {
	return session.Middleware(Store)
}
