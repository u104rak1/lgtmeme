package handler

import (
	"fmt"
	"net/http"

	"github.com/gomodule/redigo/redis"
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
	_, err := h.sessionManagerRepository.LoadGeneralAccessToken(c)
	if err == redis.ErrNil {
		respBody, status, err := h.accessTokenService.CallTokenWithClientCredentials(c)
		if err != nil && status != http.StatusOK {
			errURL := fmt.Sprintf("%s?code=%d", config.ERROR_VIEW_ENDPOINT, status)
			return c.Redirect(http.StatusFound, errURL)
		}

		if err := h.sessionManagerRepository.CacheGeneralAccessToken(c, respBody.AccessToken); err != nil {
			return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
		}
	} else if err != nil {
		return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
	}

	return c.File(config.HOME_VIEW_FILEPATH)
}
