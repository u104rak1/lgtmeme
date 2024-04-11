package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/internal/service"
)

type HomeHandler interface {
	GetView(c echo.Context) error
}

type homeHandler struct {
	sessionManagerRepository repository.SessionManager
	accessTokenService       service.AccessTokenService
}

func NewHomeHandler(
	sessionManagerRepository repository.SessionManager,
	accessTokenService service.AccessTokenService,
) *homeHandler {
	return &homeHandler{
		sessionManagerRepository: sessionManagerRepository,
		accessTokenService:       accessTokenService,
	}
}

func (h *homeHandler) GetView(c echo.Context) error {
	sessName := config.GENERAL_ACCESS_TOKEN_SESSION_NAME

	accessToken, err := h.sessionManagerRepository.LoadToken(c, sessName)
	if err != nil {
		return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
	}

	if accessToken == "" {
		respBody, status, err := h.accessTokenService.CallTokenWithClientCredentials(c)
		if err != nil && status != http.StatusOK {
			errURL := fmt.Sprintf("%s?code=%d", config.ERROR_VIEW_ENDPOINT, status)
			return c.Redirect(http.StatusFound, errURL)
		}

		if err := h.sessionManagerRepository.CacheToken(c, respBody.AccessToken, sessName); err != nil {
			return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
		}
	}

	return c.File(config.HOME_VIEW_FILEPATH)
}
