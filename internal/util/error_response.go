package util

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/my_authn_authz/config"
	"github.com/ucho456job/my_authn_authz/internal/dto"
)

// 302 Found with error
func RedirectWithErrorForAuthz(c echo.Context, q dto.AuthorizationQuery, errCode, errDescription string) error {
	redirectURL, err := url.Parse(q.RedirectURI)
	if err != nil {
		return InternalServerErrorResponse(c, err)
	}

	query := redirectURL.Query()
	query.Set("error", errCode)
	query.Set("error_description", errDescription)
	if q.State != "" {
		query.Set("state", q.State)
	}
	if q.Nonce != "" {
		query.Set("nonce", q.Nonce)
	}
	redirectURL.RawQuery = query.Encode()

	return c.Redirect(http.StatusFound, redirectURL.String())
}

// 401 Unauthorized
func UnauthorizedErrorResponse(c echo.Context, err error) error {
	config.Logger.Warn("Unauthorized", "error", err.Error())
	return c.JSON(http.StatusUnauthorized, map[string]string{
		"errorCode":    "unauthorized",
		"errorMessage": err.Error(),
	})
}

// 500 Internal Server Error
func InternalServerErrorResponse(c echo.Context, err error) error {
	config.Logger.Error("Internal server error", "error", err.Error())
	return c.JSON(http.StatusInternalServerError, map[string]string{
		"errorCode":    "internal_server_error",
		"errorMessage": err.Error(),
	})
}
