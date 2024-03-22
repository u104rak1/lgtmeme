package util

import "time"

// api endpoint
const (
	HEALTH_ENDPOINT         = "/api/health"
	LOGIN_ENDPOINT          = "/api/login"
	AUTHORAIZETION_ENDPOINT = "/api/connect/authorize"
	TOKEN_ENDPOINT          = "/api/connect/token"
)

// screen endpoint & file path
const (
	LOGIN_SCREEN_ENDPOINT = "/login"
	LOGIN_SCREEN_FILEPATH = "view/out/login.html"
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
)

// number
const (
	REFRESH_TOKEN_SIZE = 64
)
