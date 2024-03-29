package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/client/repository"
	"github.com/ucho456job/lgtmeme/internal/client/service"
)

type HomeViewHandler interface {
	GetHomeView(c echo.Context) error
}

type homeViewHandler struct {
	sessionManagerRepository repository.SessionManager
	clientCredentialsService service.ClientCredentialsService
}

func NewHomeViewHandler(
	sessionManagerRepository repository.SessionManager,
	clientCredentialsService service.ClientCredentialsService,
) *homeViewHandler {
	return &homeViewHandler{
		sessionManagerRepository: sessionManagerRepository,
		clientCredentialsService: clientCredentialsService,
	}
}

func (h *homeViewHandler) GetHomeView(c echo.Context) error {
	accessToken, err := h.sessionManagerRepository.LoadClientCredentialsAccessToken(c)
	if err != nil {
		return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
	}

	if accessToken == "" {
		accessToken, statusCode, err := h.clientCredentialsService.GetAccessToken(c)
		if err != nil && statusCode != http.StatusOK {
			errURL := fmt.Sprintf("%s?code=%d", config.ERROR_VIEW_ENDPOINT, statusCode)
			return c.Redirect(http.StatusFound, errURL)
		}

		if err := h.sessionManagerRepository.CacheClientCredentialsAccessToken(c, accessToken); err != nil {
			return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
		}
	}

	return c.File(config.HOME_VIEW_FILEPATH)
}
