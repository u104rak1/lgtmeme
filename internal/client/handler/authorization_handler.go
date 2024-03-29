package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
)

type AuthorizationHandler interface {
	RedirectAuthz(c echo.Context) error
}

type authorizationHandler struct{}

func NewAuthorizationHandler() *authorizationHandler {
	return &authorizationHandler{}
}

func (h *authorizationHandler) RedirectAuthz(c echo.Context) error {
	baseURL := os.Getenv("BASE_URL")
	url := baseURL + config.AUTHORAIZETION_ENDPOINT
	clientID := os.Getenv("OWNER_CLIENT_ID")
	redirectURI := os.Getenv("OWNER_REDIRECT_URI")
	scope := os.Getenv("OWNER_SCOPE")
	state := uuid.New().String()
	nonce := uuid.New().String()
	q := fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s&nonce=%s",
		url, clientID, redirectURI, scope, state, nonce)

	return c.Redirect(http.StatusFound, q)
}
