package config

import "time"

// api endpoint
const (
	AUTHORAIZETION_ENDPOINT = "/api/connect/authorize"
	HEALTH_ENDPOINT         = "/api/health"
	JWKS_ENDPOINT           = "/api/jwks"
	LOGIN_ENDPOINT          = "/api/login"
	LOGOUT_ENDPOINT         = "/api/logout"
	TOKEN_ENDPOINT          = "/api/connect/token"
)

// view endpoint & file path
const (
	STATIC_ENDPOINT       = "/"
	STATIC_FILEPATH       = "view/out"
	ERROR_VIEW_ENDPOINT   = "/error"
	ERROR_VIEW_FILEPATH   = "view/out/error.html"
	HOME_VIEW_ENDPOINT    = "/"
	HOME_VIEW_FILEPATH    = "view/out/index.html"
	LOGIN_VIEW_ENDPOINT   = "/login"
	LOGIN_VIEW_FILEPATH   = "view/out/login.html"
	PASSKEY_VIEW_ENDPOINT = "/passkey"
	PASSKEY_VIEW_FILEPATH = "view/out/passkey.html"
)

// session name
const (
	LOGIN_SESSION_NAME                           = "login_session"
	PRE_AUTHN_SESSION_NAME                       = "pre_authn_session"
	CLIENT_CREDENTIALS_ACCESS_TOKEN_SESSION_NAME = "client_credentials_access_token"
)

// session expire
const (
	DEFAULT_SESSION_EXPIRE_SEC = 60 * 60 * 23
	AUTHZ_CODE_EXPIRE_SEC      = 60
)

// token expire
const (
	ACCESS_TOKEN_EXPIRES_IN = time.Hour * 24
	ID_TOKEN_EXPIRES_IN     = time.Minute * 10
)

// number
const (
	REFRESH_TOKEN_SIZE = 64
)
