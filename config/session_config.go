package config

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/boj/redistore"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var Store *redistore.RediStore
var Pool *redis.Pool

func NewSessionStore() {
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

	secure, err := strconv.ParseBool(os.Getenv("COOKIE_SECURE"))
	if err != nil {
		panic(err)
	}

	Store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   DEFAULT_SESSION_EXPIRE_SEC,
		Secure:   secure,
	}
	if secure {
		Store.Options.SameSite = http.SameSiteNoneMode
	}

	Pool = Store.Pool
}

func SessionMiddleware() echo.MiddlewareFunc {
	return session.Middleware(Store)
}
