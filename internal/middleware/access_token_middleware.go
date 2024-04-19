package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/internal/service"
	"github.com/ucho456job/lgtmeme/internal/util/response"
)

type AccessTokenMiddleware interface {
	SetGeneralAccessToken() echo.MiddlewareFunc
	SetAdminAccessToken() echo.MiddlewareFunc
	VerifyAccessToken(requiredScope string) echo.MiddlewareFunc
}

type accessTokenMiddleware struct {
	sessionManagerRepository repository.SessionManagerRepository
	accessTokenService       service.AccessTokenService
}

func NewAccessTokenMiddleware(
	sessionManagerRepository repository.SessionManagerRepository,
	accessTokenService service.AccessTokenService,
) AccessTokenMiddleware {
	return &accessTokenMiddleware{
		sessionManagerRepository: sessionManagerRepository,
		accessTokenService:       accessTokenService,
	}
}

func (m *accessTokenMiddleware) SetGeneralAccessToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, err := m.sessionManagerRepository.LoadGeneralAccessToken(c)
			if err != nil {
				if err == redis.ErrNil {
					respBody, status, err := m.accessTokenService.CallTokenWithClientCredentials(c)
					if err != nil && status != http.StatusOK {
						return response.HandleErrResp(c, status, err)
					}

					if err := m.sessionManagerRepository.CacheGeneralAccessToken(c, respBody.AccessToken); err != nil {
						return response.InternalServerError(c, err)
					}
				} else {
					return response.InternalServerError(c, err)
				}
			}
			return next(c)
		}
	}
}

func (m *accessTokenMiddleware) SetAdminAccessToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			accessToken, err := m.sessionManagerRepository.LoadToken(c, config.ADMIN_ACCESS_TOKEN_SESSION_NAME)
			if err != nil && err != redis.ErrNil {
				return response.InternalServerError(c, err)
			}
			if accessToken != "" {
				return next(c)
			}

			refreshToken, err := m.sessionManagerRepository.LoadToken(c, config.REFRESH_TOKEN_SESSION_NAME)
			if err != nil && err != redis.ErrNil {
				return response.InternalServerError(c, err)
			}
			if refreshToken != "" {
				respBody, status, err := m.accessTokenService.CallTokenWithRefreshToken(c, &refreshToken)
				if err != nil {
					return response.HandleErrResp(c, status, err)
				}

				if err := m.sessionManagerRepository.CacheToken(c, respBody.AccessToken, config.ADMIN_ACCESS_TOKEN_SESSION_NAME); err != nil {
					return response.InternalServerError(c, err)
				}

				if err := m.sessionManagerRepository.CacheToken(c, respBody.RefreshToken, config.REFRESH_TOKEN_SESSION_NAME); err != nil {
					return response.InternalServerError(c, err)
				}

				return next(c)
			}

			baseURL := os.Getenv("BASE_URL")
			url := baseURL + config.AUTHZ_ENDPOINT
			clientID := os.Getenv("ADMIN_CLIENT_ID")
			redirectURI := os.Getenv("ADMIN_REDIRECT_URI")
			scope := os.Getenv("ADMIN_SCOPE")
			state := uuid.New().String()
			nonce := uuid.New().String()

			if err := m.sessionManagerRepository.CacheStateAndNonce(c, state, nonce); err != nil {
				return response.InternalServerError(c, err)
			}

			q := fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s&nonce=%s",
				url, clientID, redirectURI, scope, state, nonce)

			return c.Redirect(http.StatusFound, q)
		}
	}
}

func (m *accessTokenMiddleware) VerifyAccessToken(requiredScope string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			accessToken := c.Request().Header.Get("Authorization")
			if accessToken == "" {
				err := errors.New("access token not found")
				return response.Unauthorized(c, err)
			}

			accessToken = strings.TrimPrefix(accessToken, "Bearer ")

			keySet, err := m.sessionManagerRepository.LoadPublicKey(c)
			if err == redis.ErrNil {
				var status int
				keySet, status, err = m.accessTokenService.CallJWKS(c)
				if err != nil {
					return response.HandleErrResp(c, status, err)
				}

				if err := m.sessionManagerRepository.CachePublicKey(c, keySet); err != nil {
					return response.InternalServerError(c, err)
				}
			} else if err != nil {
				return response.InternalServerError(c, err)
			}

			parsedToken, err := jwt.Parse([]byte(accessToken), jwt.WithKeySet(keySet))
			if err != nil {
				return response.Unauthorized(c, err)
			}

			if err := jwt.Validate(parsedToken); err != nil {
				return response.Unauthorized(c, err)
			}

			if ok := tokenHasScope(parsedToken, requiredScope); !ok {
				return response.Unauthorized(c, errors.New("invalid scope"))
			}

			if sub, ok := parsedToken.Get("sub"); ok {
				c.Set("userID", sub)
			}

			return next(c)
		}
	}
}

func tokenHasScope(token jwt.Token, scope string) bool {
	scopesIf, ok := token.Get("scope")
	if !ok {
		return false
	}

	scopes, ok := scopesIf.(string)
	if !ok {
		return false
	}

	for _, s := range strings.Split(scopes, " ") {
		if s == scope {
			return true
		}
	}
	return false
}
