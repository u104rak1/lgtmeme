package middleware

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/util/response"
)

type jwksCache struct {
	keySet     jwk.Set
	lastUpdate time.Time
	mux        sync.Mutex
}

var cache = &jwksCache{}

const cacheDuration = 24 * time.Hour

func VerifyAccessToken(requiredScope string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			accessToken := c.Request().Header.Get("Authorization")
			if accessToken == "" {
				err := errors.New("access token not found")
				return response.Unauthorized(c, err)
			}

			accessToken = strings.TrimPrefix(accessToken, "Bearer ")

			keySet, err := callJWKS()
			if err != nil {
				return response.InternalServerError(c, err)
			}

			parsedToken, err := jwt.Parse([]byte(accessToken), jwt.WithKeySet(keySet))
			if err != nil {
				return response.Unauthorized(c, err)
			}

			if ok := tokenHasScope(parsedToken, requiredScope); !ok {
				return response.Unauthorized(c, errors.New("invalid scope"))
			}

			return next(c)
		}
	}
}

func callJWKS() (keySet jwk.Set, err error) {
	cache.mux.Lock()
	defer cache.mux.Unlock()

	if cache.keySet != nil && time.Since(cache.lastUpdate) < cacheDuration {
		return cache.keySet, nil
	}

	baseURL := os.Getenv("BASE_URL")
	url := baseURL + config.JWKS_ENDPOINT

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/jwk-set+json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get access token")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	keySet, err = jwk.Parse(body)
	if err != nil {
		return nil, err
	}

	cache.keySet = keySet
	cache.lastUpdate = time.Now()

	return keySet, nil
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
