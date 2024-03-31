package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/client/service"
)

type AuthzHandler interface {
	RedirectAuthz(c echo.Context) error
	Callback(c echo.Context) error
}

type authzHandler struct {
	ownerAccessTokenService service.OwnerAccessTokenService
}

func NewAuthzHandler(ownerAccessTokenService service.OwnerAccessTokenService) *authzHandler {
	return &authzHandler{}
}

func (h *authzHandler) RedirectAuthz(c echo.Context) error {
	baseURL := os.Getenv("BASE_URL")
	url := baseURL + config.AUTHZ_ENDPOINT
	clientID := os.Getenv("OWNER_CLIENT_ID")
	redirectURI := os.Getenv("OWNER_REDIRECT_URI")
	scope := os.Getenv("OWNER_SCOPE")
	state := uuid.New().String()
	nonce := uuid.New().String()
	q := fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s&nonce=%s",
		url, clientID, redirectURI, scope, state, nonce)

	return c.Redirect(http.StatusFound, q)
}

func (h *authzHandler) Callback(c echo.Context) error {
	// tokenRespBody, status, err := h.ownerAccessTokenService.CallToken(c)
	// if err != nil && status != http.StatusOK {
	// 	errURL := fmt.Sprintf("%s?code=%d", config.ERROR_VIEW_ENDPOINT, status)
	// 	return c.Redirect(http.StatusFound, errURL)
	// }
	return nil
}
