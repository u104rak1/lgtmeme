package middleware

import (
	"fmt"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/internal/service"
)

type GeneralAccessTokenMiddleware interface {
	Set() echo.MiddlewareFunc
}

type generalAccessTokenMiddleware struct {
	sessionManagerRepository repository.SessionManager
	accessTokenService       service.AccessTokenService
}

func NewGeneralAccessTokenMiddleware(
	sessionManagerRepository repository.SessionManager,
	accessTokenService service.AccessTokenService,
) GeneralAccessTokenMiddleware {
	return &generalAccessTokenMiddleware{
		sessionManagerRepository: sessionManagerRepository,
		accessTokenService:       accessTokenService,
	}
}

func (m *generalAccessTokenMiddleware) Set() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, err := m.sessionManagerRepository.LoadGeneralAccessToken(c)
			if err == redis.ErrNil {
				respBody, status, err := m.accessTokenService.CallTokenWithClientCredentials(c)
				if err != nil && status != http.StatusOK {
					errURL := fmt.Sprintf("%s?code=%d", config.ERROR_VIEW_ENDPOINT, status)
					return c.Redirect(http.StatusFound, errURL)
				}

				if err := m.sessionManagerRepository.CacheGeneralAccessToken(c, respBody.AccessToken); err != nil {
					return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
				}
			} else if err != nil {
				return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
			}
			return next(c)
		}
	}
}
