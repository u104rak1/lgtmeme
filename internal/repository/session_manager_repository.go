package repository

// mockgen -source=internal/repository/session_manager_repository.go -destination=testutil/mock/repository/mock_session_manager_repository.go -package=repository_mock

import (
	"encoding/json"
	"os"

	"github.com/boj/redistore"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/dto"

	"github.com/labstack/echo/v4"
)

type SessionManagerRepository interface {
	CacheLoginSession(c echo.Context, userID uuid.UUID) error
	LoadLoginSession(c echo.Context) (userID uuid.UUID, isLogin bool, err error)

	CachePreAuthnSession(c echo.Context, q dto.AuthzQuery) error
	LoadPreAuthnSession(c echo.Context) (query *dto.AuthzQuery, exists bool, err error)

	CacheAuthzCodeCtx(c echo.Context, q dto.AuthzQuery, authzCode string, userID uuid.UUID) error
	LoadAuthzCodeCtx(c echo.Context, code string) (*AuthzCodeCtx, error)

	Logout(c echo.Context) error

	CacheToken(c echo.Context, token, sessionName string) error
	LoadToken(c echo.Context, sessionName string) (string, error)

	CacheGeneralAccessToken(c echo.Context, token string) error
	LoadGeneralAccessToken(c echo.Context) (string, error)

	CacheStateAndNonce(c echo.Context, state, nonce string) error
	LoadStateAndNonce(c echo.Context) (state string, nonce string, err error)

	CachePublicKey(c echo.Context, keySet jwk.Set) error
	LoadPublicKey(c echo.Context) (keySet jwk.Set, err error)

	CheckRedis(c echo.Context, key string) (string, error)
}

type sessionManagerRepository struct {
	store *redistore.RediStore
	pool  *redis.Pool
}

func NewSessionManagerRepository(store *redistore.RediStore, pool *redis.Pool) SessionManagerRepository {
	return &sessionManagerRepository{
		store: store,
		pool:  pool,
	}
}

func (r *sessionManagerRepository) CacheLoginSession(c echo.Context, userID uuid.UUID) error {
	sess, err := r.store.Get(c.Request(), config.LOGIN_SESSION_NAME)
	if err != nil {
		return err
	}

	sess.Values["userId"] = userID.String()
	sess.Values["isLogin"] = true

	return sess.Save(c.Request(), c.Response())
}

func (r *sessionManagerRepository) LoadLoginSession(c echo.Context) (userID uuid.UUID, isLogin bool, err error) {
	sess, err := r.store.Get(c.Request(), config.LOGIN_SESSION_NAME)
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

func (r *sessionManagerRepository) CachePreAuthnSession(c echo.Context, q dto.AuthzQuery) error {
	sess, err := r.store.Get(c.Request(), config.PRE_AUTHN_SESSION_NAME)
	if err != nil {
		return err
	}

	sess.Values["responseType"] = q.ResponseType
	sess.Values["clientID"] = q.ClientID.String()
	sess.Values["redirectURI"] = q.RedirectURI
	sess.Values["scope"] = q.Scope
	sess.Values["state"] = q.State
	sess.Values["nonce"] = q.Nonce

	return sess.Save(c.Request(), c.Response())
}

func (r *sessionManagerRepository) LoadPreAuthnSession(c echo.Context) (query *dto.AuthzQuery, exists bool, err error) {
	sess, err := r.store.Get(c.Request(), config.PRE_AUTHN_SESSION_NAME)
	if err != nil {
		return nil, false, err
	}

	responseType, ok := sess.Values["responseType"].(string)
	if !ok {
		return nil, false, nil
	}

	clientIDStr, ok := sess.Values["clientID"].(string)
	if !ok {
		return nil, false, nil
	}

	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		return nil, false, nil
	}

	redirectURI, ok := sess.Values["redirectURI"].(string)
	if !ok {
		return nil, false, nil
	}

	scope, ok := sess.Values["scope"].(string)
	if !ok {
		scope = ""
	}

	state, ok := sess.Values["state"].(string)
	if !ok {
		return nil, false, nil
	}

	nonce, ok := sess.Values["nonce"].(string)
	if !ok {
		nonce = ""
	}

	query = &dto.AuthzQuery{
		ResponseType: responseType,
		ClientID:     clientID,
		RedirectURI:  redirectURI,
		Scope:        scope,
		State:        state,
		Nonce:        nonce,
	}

	if err := r.clearSession(c, config.PRE_AUTHN_SESSION_NAME); err != nil {
		return nil, false, err
	}

	return query, true, nil
}

type AuthzCodeCtx struct {
	UserID      uuid.UUID `json:"userId"`
	ClientID    uuid.UUID `json:"clientId"`
	Scope       string    `json:"scope"`
	RedirectURI string    `json:"redirectUri"`
	Nonce       string    `json:"nonce"`
}

func (r *sessionManagerRepository) CacheAuthzCodeCtx(c echo.Context, q dto.AuthzQuery, authzCode string, userID uuid.UUID) error {
	saveData := AuthzCodeCtx{
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

	conn := r.pool.Get()
	defer conn.Close()

	_, err = conn.Do("SET", authzCode, encodedData, "EX", config.AUTHZ_CODE_EXPIRE_SEC)
	if err != nil {
		return err
	}

	return nil
}

func (r *sessionManagerRepository) LoadAuthzCodeCtx(c echo.Context, code string) (*AuthzCodeCtx, error) {
	conn := r.pool.Get()
	defer conn.Close()

	value, err := redis.String(conn.Do("GET", code))
	if err != nil {
		return nil, err
	}

	var ctx AuthzCodeCtx
	if err := json.Unmarshal([]byte(value), &ctx); err != nil {
		return nil, err
	}

	if err := r.clearSession(c, code); err != nil {
		return nil, err
	}

	return &ctx, nil
}

func (r *sessionManagerRepository) Logout(c echo.Context) error {
	if err := r.clearSession(c, config.LOGIN_SESSION_NAME); err != nil {
		return err
	}
	if err := r.clearSession(c, config.ADMIN_ACCESS_TOKEN_SESSION_NAME); err != nil {
		return err
	}
	return nil
}

func (r *sessionManagerRepository) clearSession(c echo.Context, sessionName string) error {
	sess, err := r.store.Get(c.Request(), sessionName)
	if err != nil {
		return err
	}

	sess.Options.MaxAge = -1

	return sess.Save(c.Request(), c.Response())
}

func (r *sessionManagerRepository) CacheToken(c echo.Context, token, sessionName string) error {
	sess, err := r.store.Get(c.Request(), sessionName)
	if err != nil {
		return err
	}

	if sessionName == config.REFRESH_TOKEN_SESSION_NAME {
		sess.Options.MaxAge = config.REFRESH_TOKEN_SESSION_EXPIRE_SEC
	}

	sess.Values[sessionName] = token

	return sess.Save(c.Request(), c.Response())
}

func (r *sessionManagerRepository) LoadToken(c echo.Context, sessionName string) (string, error) {
	sess, err := r.store.Get(c.Request(), sessionName)
	if err != nil {
		return "", err
	}

	token, ok := sess.Values[sessionName].(string)
	if !ok {
		return "", nil
	}

	return token, nil
}

func (r *sessionManagerRepository) CacheGeneralAccessToken(c echo.Context, token string) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", os.Getenv("GENERAL_ACCESS_TOKEN_REDIS_KEY"), token, "EX", config.DEFAULT_SESSION_EXPIRE_SEC)
	if err != nil {
		return err
	}

	return nil
}

func (r *sessionManagerRepository) LoadGeneralAccessToken(c echo.Context) (string, error) {
	conn := r.pool.Get()
	defer conn.Close()

	token, err := redis.String(conn.Do("GET", os.Getenv("GENERAL_ACCESS_TOKEN_REDIS_KEY")))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *sessionManagerRepository) CacheStateAndNonce(c echo.Context, state, nonce string) error {
	sess, err := r.store.Get(c.Request(), config.STATE_AND_NONCE_SESSION_NAME)
	if err != nil {
		return err
	}

	sess.Values["state"] = state
	sess.Values["nonce"] = nonce

	return sess.Save(c.Request(), c.Response())
}

func (r *sessionManagerRepository) LoadStateAndNonce(c echo.Context) (state string, nonce string, err error) {
	sess, err := r.store.Get(c.Request(), config.STATE_AND_NONCE_SESSION_NAME)
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

	if err := r.clearSession(c, config.STATE_AND_NONCE_SESSION_NAME); err != nil {
		return "", "", err
	}

	return state, nonce, nil
}

func (r *sessionManagerRepository) CachePublicKey(c echo.Context, keySet jwk.Set) error {
	jsonKeySet, err := json.Marshal(keySet)
	if err != nil {
		return err
	}

	conn := r.pool.Get()
	defer conn.Close()

	_, err = conn.Do("SET", os.Getenv("JWKS_REDIS_KEY"), jsonKeySet, "EX", config.DEFAULT_SESSION_EXPIRE_SEC)
	if err != nil {
		return err
	}

	return nil
}

func (r *sessionManagerRepository) LoadPublicKey(c echo.Context) (keySet jwk.Set, err error) {
	conn := r.pool.Get()
	defer conn.Close()

	value, err := redis.String(conn.Do("GET", os.Getenv("JWKS_REDIS_KEY")))
	if err != nil {
		return nil, err
	}

	keySet, err = jwk.Parse([]byte(value))
	if err != nil {
		return nil, err
	}

	return keySet, nil
}

func (r *sessionManagerRepository) CheckRedis(c echo.Context, key string) (value string, err error) {
	conn := r.pool.Get()
	defer conn.Close()

	value, err = redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}

	return value, nil
}
