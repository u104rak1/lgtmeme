package repository

import (
	"github.com/boj/redistore"
	"github.com/gomodule/redigo/redis"
	"github.com/ucho456job/lgtmeme/config"

	"github.com/labstack/echo/v4"
)

type SessionManager interface {
	CacheToken(c echo.Context, token, sessionName string) error
	LoadToken(c echo.Context, sessionName string) (string, error)
	CacheStateAndNonce(c echo.Context, state, nonce string) error
	LoadStateAndNonce(c echo.Context) (state string, nonce string, err error)
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

func (sm *sessionManager) CacheToken(c echo.Context, token, sessionName string) error {
	sess, err := sm.store.Get(c.Request(), sessionName)
	if err != nil {
		return err
	}

	if sessionName == config.REFRESH_TOKEN_SESSION_NAME {
		sess.Options.MaxAge = config.REFRESH_TOKEN_SESSION_EXPIRE_SEC
	}

	sess.Values[sessionName] = token

	return sess.Save(c.Request(), c.Response())
}

func (sm *sessionManager) LoadToken(c echo.Context, sessionName string) (string, error) {
	sess, err := sm.store.Get(c.Request(), sessionName)
	if err != nil {
		return "", err
	}

	token, ok := sess.Values[sessionName].(string)
	if !ok {
		return "", nil
	}

	return token, nil
}

func (sm *sessionManager) CacheStateAndNonce(c echo.Context, state, nonce string) error {
	sess, err := sm.store.Get(c.Request(), config.STATE_AND_NONCE_SESSION_NAME)
	if err != nil {
		return err
	}

	sess.Values["state"] = state
	sess.Values["nonce"] = nonce

	return sess.Save(c.Request(), c.Response())
}

func (sm *sessionManager) LoadStateAndNonce(c echo.Context) (state string, nonce string, err error) {
	sess, err := sm.store.Get(c.Request(), config.STATE_AND_NONCE_SESSION_NAME)
	if err != nil {
		return "", "", err
	}

	state, ok := sess.Values["state"].(string)
	if !ok {
		return "", "", nil
	}

	nonce, ok = sess.Values["nonce"].(string)
	if !ok {
		return "", "", nil
	}

	return state, nonce, nil
}
