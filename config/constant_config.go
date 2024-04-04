package config

import "time"

// auth server
const (
	// api endpoint
	AUTHZ_ENDPOINT  = "/auth-api/authorize"
	HEALTH_ENDPOINT = "/auth-api/health"
	JWKS_ENDPOINT   = "/auth-api/jwks"
	LOGIN_ENDPOINT  = "/auth-api/login"
	LOGOUT_ENDPOINT = "/auth-api/logout"
	TOKEN_ENDPOINT  = "/auth-api/token"

	// view endpoint
	LOGIN_VIEW_ENDPOINT = "/login"

	// file path
	LOGIN_VIEW_FILEPATH = "view/out/login.html"
)

// client server
const (
	// api endpoint
	CLIENT_AUTH_ENDPOINT          = "/client-api/auth"
	CLIENT_AUTH_CALLBACK_ENDPOINT = "/client-api/auth/callback"
	CLIENT_IMAGES_ENDPOINT        = "/client-api/images"

	// view endpoint
	STATIC_ENDPOINT            = "/"
	ERROR_VIEW_ENDPOINT        = "/error"
	HOME_VIEW_ENDPOINT         = "/"
	CREATE_IMAGE_VIEW_ENDPOINT = "/create-image"
	PASSKEY_VIEW_ENDPOINT      = "/passkey"
	AUTH_VIEW_ENDPOINT         = "/auth"

	// file path
	STATIC_FILEPATH            = "view/out"
	ERROR_VIEW_FILEPATH        = "view/out/error.html"
	HOME_VIEW_FILEPATH         = "view/out/index.html"
	CREATE_IMAGE_VIEW_FILEPATH = "view/out/create-image.html"
	PASSKEY_VIEW_FILEPATH      = "view/out/passkey.html"
	AUTH_VIEW_FILEPATH         = "view/out/auth.html"
)

// resoruce server
const (
	// api endpoint
	RESOURCE_IMAGES_ENDPOINT = "/resource-api/images"
)

// session name
const (
	LOGIN_SESSION_NAME                = "login_session"
	PRE_AUTHN_SESSION_NAME            = "pre_authn_session"
	GENERAL_ACCESS_TOKEN_SESSION_NAME = "general_access_token"
	OWNER_ACCESS_TOKEN_SESSION_NAME   = "owner_access_token"
	STATE_AND_NONCE_SESSION_NAME      = "state_and_nonce"
	REFRESH_TOKEN_SESSION_NAME        = "refresh_token"
)

// session expire
const (
	DEFAULT_SESSION_EXPIRE_SEC       = 60 * 60 * 23
	AUTHZ_CODE_EXPIRE_SEC            = 60
	REFRESH_TOKEN_SESSION_EXPIRE_SEC = 60 * 60 * 24 * 30
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

// scope
const (
	IMAGES_READ_SCOPE   = "images.read"
	IMAGES_CREATE_SCOPE = "images.create"
	IMAGES_UPDATE_SCOPE = "images.update"
	IMAGES_DELETE_SCOPE = "images.delete"
)
