package util

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	authDto "github.com/ucho456job/lgtmeme/internal/auth/dto"
)

// 302 Found with error
func RedirectWithErrorForAuthz(c echo.Context, q authDto.AuthorizationQuery, errCode, errDescription string) error {
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

// 400 Bad Request
func BadRequestResponse(c echo.Context, err error) error {
	config.Logger.Warn("Bad request", "error", err.Error())
	return c.JSON(http.StatusBadRequest, map[string]string{
		"errorCode":    "bad_request",
		"errorMessage": err.Error(),
	})
}

// 401 Unauthorized
func UnauthorizedErrorResponse(c echo.Context, err error) error {
	config.Logger.Warn("Unauthorized", "error", err.Error())
	return c.JSON(http.StatusUnauthorized, map[string]string{
		"errorCode":    "unauthorized",
		"errorMessage": err.Error(),
	})
}

// 404 Not Found
func NotFoundErrorResponse(c echo.Context, err error) error {
	config.Logger.Warn("Not found", "error", err.Error())
	return c.JSON(http.StatusNotFound, map[string]string{
		"errorCode":    "not_found",
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
