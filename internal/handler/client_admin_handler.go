package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/internal/service"
)

type AdminHandler interface {
	GetView(c echo.Context) error
	RedirectAuthz(c echo.Context) error
	Callback(c echo.Context) error
}

type adminHandler struct {
	sessionManagerRepository repository.SessionManager
	adminAccessTokenService  service.AdminAccessTokenService
}

func NewAdminHandler(
	sessionManagerRepository repository.SessionManager,
	adminAccessTokenService service.AdminAccessTokenService,
) *adminHandler {
	return &adminHandler{
		sessionManagerRepository: sessionManagerRepository,
		adminAccessTokenService:  adminAccessTokenService,
	}
}

func (h *adminHandler) GetView(c echo.Context) error {
	return c.File(config.ADMIN_VIEW_FILEPATH)
}

func (h *adminHandler) RedirectAuthz(c echo.Context) error {
	accessToken, err := h.sessionManagerRepository.LoadToken(c, config.ADMIN_ACCESS_TOKEN_SESSION_NAME)
	if err != nil {
		return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
	}
	if accessToken != "" {
		return c.Redirect(http.StatusFound, config.ADMIN_VIEW_ENDPOINT)
	}

	refreshToken, err := h.sessionManagerRepository.LoadToken(c, config.REFRESH_TOKEN_SESSION_NAME)
	if err != nil {
		return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
	}
	if refreshToken != "" {
		respBody, status, err := h.adminAccessTokenService.CallTokenWithRefreshToken(c, &refreshToken)
		if err != nil {
			errURL := fmt.Sprintf("%s?code=%d", config.ERROR_VIEW_ENDPOINT, status)
			return c.Redirect(http.StatusFound, errURL)
		}

		return h.commonSuccessProcess(c, respBody.AccessToken, respBody.RefreshToken)
	}

	baseURL := os.Getenv("BASE_URL")
	url := baseURL + config.AUTHZ_ENDPOINT
	clientID := os.Getenv("ADMIN_CLIENT_ID")
	redirectURI := os.Getenv("ADMIN_REDIRECT_URI")
	scope := os.Getenv("ADMIN_SCOPE")
	state := uuid.New().String()
	nonce := uuid.New().String()

	if err := h.sessionManagerRepository.CacheStateAndNonce(c, state, nonce); err != nil {
		return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
	}

	q := fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s&nonce=%s",
		url, clientID, redirectURI, scope, state, nonce)

	return c.Redirect(http.StatusFound, q)
}

func (h *adminHandler) Callback(c echo.Context) error {
	tokenRespBody, status, err := h.adminAccessTokenService.CallToken(c)
	if err != nil {
		errURL := fmt.Sprintf("%s?code=%d", config.ERROR_VIEW_ENDPOINT, status)
		return c.Redirect(http.StatusFound, errURL)
	}

	keySet, status, err := h.adminAccessTokenService.CallJWKS(c)
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

	return h.commonSuccessProcess(c, tokenRespBody.AccessToken, tokenRespBody.RefreshToken)
}

func (h *adminHandler) commonSuccessProcess(c echo.Context, accessToken, refreshToken string) error {
	if err := h.sessionManagerRepository.CacheToken(c, accessToken, config.ADMIN_ACCESS_TOKEN_SESSION_NAME); err != nil {
		return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
	}

	if err := h.sessionManagerRepository.CacheToken(c, refreshToken, config.REFRESH_TOKEN_SESSION_NAME); err != nil {
		return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
	}

	return c.Redirect(http.StatusFound, config.ADMIN_VIEW_ENDPOINT)
}
