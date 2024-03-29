package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/auth/dto"
	"github.com/ucho456job/lgtmeme/internal/auth/repository"
	"github.com/ucho456job/lgtmeme/internal/auth/service"
	"github.com/ucho456job/lgtmeme/internal/util"
)

type TokenHandler interface {
	GenerateToken(c echo.Context) error
}

type tokenHandler struct {
	oauthClientRepository    repository.OauthClientRepository
	refreshTokenRepository   repository.RefreshTokenRepository
	userRepository           repository.UserRepository
	sessionManagerRepository repository.SessionManager
	jwtService               service.JwtService
}

func NewTokenHandler(
	oauthClientRepository repository.OauthClientRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
	userRepository repository.UserRepository,
	sessionManagerRepository repository.SessionManager,
	jwtService service.JwtService,
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

	expiresIn := config.ACCESS_TOKEN_EXPIRES_IN

	switch form.GrantType {
	case "authorization_code":
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

		accessToken, err := h.jwtService.GenerateAccessToken(&user.ID, oauthClient, expiresIn)
		if err != nil {
			return util.InternalServerErrorResponse(c, err)
		}

		idToken, err := h.jwtService.GenerateIDToken(oauthClient, user, authzCodeCtx.Nonce)
		if err != nil {
			return util.InternalServerErrorResponse(c, err)
		}

		refreshToken, err := h.jwtService.GenerateRefreshToken()
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
	case "refresh_token":
		refreshTokenData, err := h.refreshTokenRepository.FindByToken(c, form.RefreshToken)
		if err != nil {
			return util.NotFoundErrorResponse(c, err)
		}

		exists, err := h.userRepository.ExistsByID(c, refreshTokenData.UserID)
		if err != nil {
			return util.InternalServerErrorResponse(c, err)
		}
		if !exists {
			return util.NotFoundErrorResponse(c, errors.New("user not found"))
		}

		accessToken, err := h.jwtService.GenerateAccessToken(&refreshTokenData.UserID, oauthClient, expiresIn)
		if err != nil {
			return util.InternalServerErrorResponse(c, err)
		}

		newRefreshToken, err := h.jwtService.GenerateRefreshToken()
		if err != nil {
			return util.InternalServerErrorResponse(c, err)
		}

		if err := h.refreshTokenRepository.UpdateRefreshToken(c, refreshTokenData.UserID, form.ClientID, newRefreshToken, refreshTokenData.Scopes); err != nil {
			return util.InternalServerErrorResponse(c, err)
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"accessToken":  accessToken,
			"tokenType":    "Bearer",
			"expiresIn":    expiresIn.Seconds(),
			"refreshToken": newRefreshToken,
		})
	case "client_credentials":
		accessToken, err := h.jwtService.GenerateAccessToken(nil, oauthClient, expiresIn)
		if err != nil {
			return util.InternalServerErrorResponse(c, err)
		}

		return c.JSON(http.StatusOK, dto.ClientCredentialsResponse{
			AccessToken: accessToken,
			TokenType:   "Bearer",
			ExpiresIn:   int(expiresIn.Seconds()),
		})
	default:
		return util.BadRequestResponse(c, errors.New("unsupported grant_type"))
	}
}
