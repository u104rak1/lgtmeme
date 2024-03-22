package repository

import (
	"encoding/json"

	"github.com/boj/redistore"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/ucho456job/my_authn_authz/internal/dto"
	"github.com/ucho456job/my_authn_authz/internal/util"

	"github.com/labstack/echo/v4"
)

type SessionManager interface {
	CacheLoginSession(c echo.Context, userID uuid.UUID) error
	LoadLoginSession(c echo.Context) (userID uuid.UUID, isLogin bool, err error)
	CachePreAuthnSession(c echo.Context, q dto.AuthorizationQuery) error
	CacheAuthzCodeWithCtx(c echo.Context, q dto.AuthorizationQuery, authzCode string, userID uuid.UUID) error
	LoadAuthzCodeWithCtx(c echo.Context, code string) (*AuthzCodeContext, error)
	CheckHealthForRedis(c echo.Context, key string) (string, error)
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

func (sm *sessionManager) CacheLoginSession(c echo.Context, userID uuid.UUID) error {
	sess, err := sm.store.Get(c.Request(), util.LOGIN_SESSION_NAME)
	if err != nil {
		return err
	}

	sess.Values["userId"] = userID.String()
	sess.Values["isLogin"] = true

	return sess.Save(c.Request(), c.Response())
}

func (sm *sessionManager) LoadLoginSession(c echo.Context) (userID uuid.UUID, isLogin bool, err error) {
	sess, err := sm.store.Get(c.Request(), util.LOGIN_SESSION_NAME)
	if err != nil {
		return uuid.Nil, false, err
	}

	val, ok := sess.Values["userId"]
	if !ok {
		return uuid.Nil, false, nil
	}

	userIDStr, ok := val.(string)
	if !ok {
		return uuid.Nil, false, nil
	}

	userID, err = uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, false, err
	}

	isLoginVal, ok := sess.Values["isLogin"]
	if !ok {
		return userID, false, nil
	}

	isLogin, ok = isLoginVal.(bool)
	if !ok {
		return userID, false, nil
	}

	return userID, isLogin, nil
}

func (sm *sessionManager) CachePreAuthnSession(c echo.Context, q dto.AuthorizationQuery) error {
	sess, err := sm.store.Get(c.Request(), util.PRE_AUTHN_SESSION_NAME)
	if err != nil {
		return err
	}

	sess.Values["responseType"] = q.ResponseType
	sess.Values["clientId"] = q.ClientID
	sess.Values["redirectUri"] = q.RedirectURI
	sess.Values["scope"] = q.Scope
	sess.Values["state"] = q.State
	sess.Values["nonce"] = q.Nonce

	return sess.Save(c.Request(), c.Response())
}

type AuthzCodeContext struct {
	UserID      uuid.UUID `json:"userId"`
	ClientID    uuid.UUID `json:"clientId"`
	Scope       string    `json:"scope"`
	RedirectURI string    `json:"redirectUri"`
	Nonce       string    `json:"nonce"`
}

func (sm *sessionManager) CacheAuthzCodeWithCtx(c echo.Context, q dto.AuthorizationQuery, authzCode string, userID uuid.UUID) error {
	saveData := AuthzCodeContext{
		UserID:      userID,
		ClientID:    q.ClientID,
		Scope:       q.Scope,
		RedirectURI: q.RedirectURI,
		Nonce:       q.Nonce,
	}

	encodedData, err := json.Marshal(saveData)
	if err != nil {
		return err
	}

	conn := sm.pool.Get()
	defer conn.Close()

	_, err = conn.Do("SET", authzCode, encodedData, "EX", 60)
	if err != nil {
		return err
	}

	return nil
}

func (sm *sessionManager) LoadAuthzCodeWithCtx(c echo.Context, code string) (*AuthzCodeContext, error) {
	conn := sm.pool.Get()
	defer conn.Close()

	value, err := redis.String(conn.Do("GET", code))
	if err != nil {
		return nil, err
	}

	var ctx AuthzCodeContext
	if err := json.Unmarshal([]byte(value), &ctx); err != nil {
		return nil, err
	}

	return &ctx, nil
}

func (sm *sessionManager) CheckHealthForRedis(c echo.Context, key string) (value string, err error) {
	conn := sm.pool.Get()
	defer conn.Close()

	value, err = redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}

	return value, nil
}
