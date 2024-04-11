package middleware

import (
	"errors"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/internal/service"
	"github.com/ucho456job/lgtmeme/internal/util/response"
)

type VerifyAccessTokenMiddleware interface {
	Verify(requiredScope string) echo.MiddlewareFunc
}

type verifyAccessTokenMiddleware struct {
	sessionManagerRepository repository.SessionManager
	accessTokenService       service.AccessTokenService
}

func NewVerifyAccessTokenMiddleware(
	sessionManagerRepository repository.SessionManager,
	accessTokenService service.AccessTokenService,
) VerifyAccessTokenMiddleware {
	return &verifyAccessTokenMiddleware{
		sessionManagerRepository: sessionManagerRepository,
		accessTokenService:       accessTokenService,
	}
}

func (m *verifyAccessTokenMiddleware) Verify(requiredScope string) echo.MiddlewareFunc {
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
