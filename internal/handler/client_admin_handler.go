package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/internal/service"
	"github.com/ucho456job/lgtmeme/internal/util/response"
)

type AdminHandler interface {
	Callback(c echo.Context) error
}

type adminHandler struct {
	sessionManagerRepository repository.SessionManagerRepository
	accessTokenService       service.AccessTokenService
}

func NewAdminHandler(
	sessionManagerRepository repository.SessionManagerRepository,
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
		return response.HandleErrResp(c, status, err)
	}

	keySet, status, err := h.accessTokenService.CallJWKS(c)
	if err != nil {
		return response.HandleErrResp(c, status, err)
	}

	state, nonce, err := h.sessionManagerRepository.LoadStateAndNonce(c)
	if err != nil {
		return response.InternalServerError(c, err)
	}

	if state != c.QueryParam("state") {
		return response.BadRequest(c, errors.New("state does not match"))
	}

	rawIDToken := tokenRespBody.IDToken

	parsedIDToken, err := jwt.Parse([]byte(rawIDToken), jwt.WithKeySet(keySet))
	if err != nil {
		return response.InternalServerError(c, err)
	}

	nonceClaim, ok := parsedIDToken.Get("nonce")
	if !ok {
		return response.InternalServerError(c, errors.New("nonce claim not found"))
	}

	if nonceClaim != nonce {
		return response.BadRequest(c, errors.New("nonce does not match"))
	}

	if err := h.sessionManagerRepository.CacheToken(c, tokenRespBody.AccessToken, config.ADMIN_ACCESS_TOKEN_SESSION_NAME); err != nil {
		return response.InternalServerError(c, err)
	}

	if err := h.sessionManagerRepository.CacheToken(c, tokenRespBody.RefreshToken, config.REFRESH_TOKEN_SESSION_NAME); err != nil {
		return response.InternalServerError(c, err)
	}

	return c.Redirect(http.StatusFound, config.ADMIN_VIEW_ENDPOINT)
}
