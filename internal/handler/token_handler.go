package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/my_authn_authz/internal/dto"
	"github.com/ucho456job/my_authn_authz/internal/repository"
	"github.com/ucho456job/my_authn_authz/internal/util"
)

type TokenHandler interface {
	GenerateToken(c echo.Context) error
}

type tokenHandler struct {
	oauthClientRepository    repository.OauthClientRepository
	refreshTokenRepository   repository.RefreshTokenRepository
	userRepository           repository.UserRepository
	sessionManagerRepository repository.SessionManager
	jwtService               util.JwtService
}

func NewTokenHandler(
	oauthClientRepository repository.OauthClientRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
	userRepository repository.UserRepository,
	sessionManagerRepository repository.SessionManager,
	jwtService util.JwtService,
) *tokenHandler {
	return &tokenHandler{
		oauthClientRepository:    oauthClientRepository,
		refreshTokenRepository:   refreshTokenRepository,
		userRepository:           userRepository,
		sessionManagerRepository: sessionManagerRepository,
		jwtService:               jwtService,
	}
}

func (h *tokenHandler) GenerateToken(c echo.Context) error {
	var form dto.TokenForm
	if err := c.Bind(&form); err != nil {
		return util.InternalServerErrorResponse(c, err)
	}
	if err := c.Validate(&form); err != nil {
		return util.BadRequestResponse(c, err)
	}

	oauthClient, err := h.oauthClientRepository.FindByClientID(c, form.ClientID)
	if err != nil {
		return util.NotFoundErrorResponse(c, err)
	}
	if oauthClient.ClientSecret != form.ClientSecret {
		return util.BadRequestResponse(c, errors.New("invalid client_secret"))
	}

	if form.GrantType == "authorization_code" {
		authzCodeCtx, err := h.sessionManagerRepository.LoadAuthzCodeWithCtx(c, form.Code)
		if err != nil {
			return util.InternalServerErrorResponse(c, err)
		}
		if authzCodeCtx.RedirectURI != form.RedirectURI || authzCodeCtx.ClientID != form.ClientID {
			return util.BadRequestResponse(c, errors.New("invalid redirect_uri or client_id"))
		}

		user, err := h.userRepository.FindByID(c, authzCodeCtx.UserID)
		if err != nil {
			return util.NotFoundErrorResponse(c, err)
		}

		expiresIn := util.ACCESS_TOKEN_EXPIRES_IN

		accessToken, err := h.jwtService.GenerateAccessToken(user.ID, oauthClient, expiresIn)
		if err != nil {
			return util.InternalServerErrorResponse(c, err)
		}

		idToken, err := h.jwtService.GenerateIDToken(oauthClient, user, authzCodeCtx.Nonce)
		if err != nil {
			return util.InternalServerErrorResponse(c, err)
		}

		refreshToken, err := h.jwtService.GenerateRefreshToken(authzCodeCtx.UserID)
		if err != nil {
			return util.InternalServerErrorResponse(c, err)
		}

		if err := h.refreshTokenRepository.CreateRefreshToken(c, user.ID, form.ClientID, refreshToken, authzCodeCtx.Scope); err != nil {
			return util.InternalServerErrorResponse(c, err)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"accessToken":  accessToken,
			"tokenType":    "Bearer",
			"expiresIn":    expiresIn.Seconds(),
			"refreshToken": refreshToken,
			"idToken":      idToken,
		})
	}
	return nil
}
