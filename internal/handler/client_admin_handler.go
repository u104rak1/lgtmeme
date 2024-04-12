package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/internal/service"
)

type AdminHandler interface {
	Callback(c echo.Context) error
}

type adminHandler struct {
	sessionManagerRepository repository.SessionManager
	accessTokenService       service.AccessTokenService
}

func NewAdminHandler(
	sessionManagerRepository repository.SessionManager,
	accessTokenService service.AccessTokenService,
) *adminHandler {
	return &adminHandler{
		sessionManagerRepository: sessionManagerRepository,
		accessTokenService:       accessTokenService,
	}
}

func (h *adminHandler) Callback(c echo.Context) error {
	tokenRespBody, status, err := h.accessTokenService.CallTokenWithAuthzCode(c)
	if err != nil {
		errURL := fmt.Sprintf("%s?code=%d", config.ERROR_VIEW_ENDPOINT, status)
		return c.Redirect(http.StatusFound, errURL)
	}

	keySet, status, err := h.accessTokenService.CallJWKS(c)
	if err != nil {
		errURL := fmt.Sprintf("%s?code=%d", config.ERROR_VIEW_ENDPOINT, status)
		return c.Redirect(http.StatusFound, errURL)
	}

	state, nonce, err := h.sessionManagerRepository.LoadStateAndNonce(c)
	if err != nil {
		return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
	}

	if state != c.QueryParam("state") {
		return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
	}

	rawIDToken := tokenRespBody.IDToken

	parsedIDToken, err := jwt.Parse([]byte(rawIDToken), jwt.WithKeySet(keySet))
	if err != nil {
		log.Printf("Error parsing JWT: %v", err)
		return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
	}

	nonceClaim, ok := parsedIDToken.Get("nonce")
	if !ok {
		log.Printf("ID token does not have a nonce claim")
		return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
	}

	if nonceClaim != nonce {
		log.Printf("Nonce does not match")
		return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
	}

	if err := h.sessionManagerRepository.CacheToken(c, tokenRespBody.AccessToken, config.ADMIN_ACCESS_TOKEN_SESSION_NAME); err != nil {
		return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
	}

	if err := h.sessionManagerRepository.CacheToken(c, tokenRespBody.RefreshToken, config.REFRESH_TOKEN_SESSION_NAME); err != nil {
		return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
	}

	return c.Redirect(http.StatusFound, config.ADMIN_VIEW_ENDPOINT)
}
