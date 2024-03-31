package repository

import (
	"github.com/boj/redistore"
	"github.com/gomodule/redigo/redis"
	"github.com/ucho456job/lgtmeme/config"

	"github.com/labstack/echo/v4"
)

type SessionManager interface {
	CacheGeneralAccessToken(c echo.Context, token string) error
	LoadGeneralAccessToken(c echo.Context) (string, error)
}

type sessionManager struct {
	store *redistore.RediStore
	pool  *redis.Pool
}

func NewSessionManager(store *redistore.RediStore, pool *redis.Pool) SessionManager {
	return &sessionManager{
		store: store,
		pool:  pool,
	}
}

func (sm *sessionManager) CacheGeneralAccessToken(c echo.Context, token string) error {
	sess, err := sm.store.Get(c.Request(), config.GENERAL_ACCESS_TOKEN_SESSION_NAME)
	if err != nil {
		return err
	}

	sess.Values["accessToken"] = token

	return sess.Save(c.Request(), c.Response())
}

func (sm *sessionManager) LoadGeneralAccessToken(c echo.Context) (string, error) {
	sess, err := sm.store.Get(c.Request(), config.GENERAL_ACCESS_TOKEN_SESSION_NAME)
	if err != nil {
		return "", err
	}

	token, ok := sess.Values["accessToken"].(string)
	if !ok {
		return "", nil
	}

	return token, nil
}
