package util

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

// screen endpoint & file path
const (
	LOGIN_SCREEN_ENDPOINT   = "/login"
	LOGIN_SCREEN_FILEPATH   = "view/out/login.html"
	PASSKEY_SCREEN_ENDPOINT = "/passkey"
	PASSKEY_SCREEN_FILEPATH = "view/out/passkey.html"
)

// session name
const (
	LOGIN_SESSION_NAME     = "login_session"
	PRE_AUTHN_SESSION_NAME = "pre_authn_session"
)

// time
const (
	ACCESS_TOKEN_EXPIRES_IN = time.Hour * 24
	ID_TOKEN_EXPIRES_IN     = time.Minute * 10

	AUTHZ_CODE_EXPIRE_SEC = 60
)

// number
const (
	REFRESH_TOKEN_SIZE = 64
)
