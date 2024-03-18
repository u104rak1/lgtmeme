package util

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/ucho456job/my_authn_authz/config"

	"github.com/labstack/echo/v4"
)

type SessionManager interface {
	SaveLoginSession(c echo.Context, userID string) error
	CheckHealthForRedis(key string) (string, error)
}

type sessionManager struct{}

func NewSessionManager() SessionManager {
	return &sessionManager{}
}

func (sm *sessionManager) SaveLoginSession(c echo.Context, userID string) error {
	sess, err := session.Get(LOGIN_SESSION_NAME, c)
	if err != nil {
		return err
	}

	sess.Values["userId"] = userID
	sess.Values["isLogin"] = true
	sess.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
	}

	secure, err := strconv.ParseBool(os.Getenv("COOKIE_SECURE"))
	if err != nil {
		return err
	}

	if secure {
		sess.Options.Secure = true
		sess.Options.SameSite = http.SameSiteNoneMode
	}

	return sess.Save(c.Request(), c.Response())
}

func (sm *sessionManager) CheckHealthForRedis(key string) (value string, err error) {
	conn := config.Store.Pool.Get()
	defer conn.Close()

	value, err = redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}

	return value, nil
}
