package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/dto"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/internal/service"
	"github.com/ucho456job/lgtmeme/internal/util/response"
)

type TokenHandler interface {
	Generate(c echo.Context) error
}

type tokenHandler struct {
	oauthClientRepository    repository.OauthClientRepository
	refreshTokenRepository   repository.RefreshTokenRepository
	userRepository           repository.UserRepository
	sessionManagerRepository repository.SessionManagerRepository
	jwtService               service.JWTService
}

func NewTokenHandler(
	oauthClientRepository repository.OauthClientRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
	userRepository repository.UserRepository,
	sessionManagerRepository repository.SessionManagerRepository,
	jwtService service.JWTService,
) *tokenHandler {
	return &tokenHandler{
		oauthClientRepository:    oauthClientRepository,
		refreshTokenRepository:   refreshTokenRepository,
		userRepository:           userRepository,
		sessionManagerRepository: sessionManagerRepository,
		jwtService:               jwtService,
	}
}

func (h *tokenHandler) Generate(c echo.Context) error {
	var form dto.TokenForm
	if err := c.Bind(&form); err != nil {
		return response.InternalServerError(c, err)
	}
	if err := c.Validate(&form); err != nil {
		return response.BadRequest(c, err)
	}

	oauthClient, err := h.oauthClientRepository.FindByClientID(c, form.ClientID)
	if err != nil {
		return response.NotFound(c, err)
	}
	if oauthClient.ClientSecret != form.ClientSecret {
		err = errors.New("invalid client_secret")
		return response.BadRequest(c, err)
	}

	expiresIn := config.ACCESS_TOKEN_EXPIRES_IN

	switch form.GrantType {
	case "authorization_code":
		authzCodeCtx, err := h.sessionManagerRepository.LoadAuthzCodeCtx(c, form.Code)
		if err != nil {
			return response.InternalServerError(c, err)
		}
		if authzCodeCtx.RedirectURI != form.RedirectURI || authzCodeCtx.ClientID != form.ClientID {
			err = errors.New("invalid redirect_uri or client_id")
			return response.BadRequest(c, err)
		}

		user, err := h.userRepository.FindByID(c, authzCodeCtx.UserID)
		if err != nil {
			return response.InternalServerError(c, err)
		}
		if user == nil {
			err = errors.New("user not found")
			return response.NotFound(c, err)
		}

		accessToken, err := h.jwtService.GenerateAccessToken(&user.ID, oauthClient, expiresIn)
		if err != nil {
			return response.InternalServerError(c, err)
		}

		idToken, err := h.jwtService.GenerateIDToken(oauthClient, user, authzCodeCtx.Nonce)
		if err != nil {
			return response.InternalServerError(c, err)
		}

		refreshToken, err := h.jwtService.GenerateRefreshToken()
		if err != nil {
			return response.InternalServerError(c, err)
		}

		if err := h.refreshTokenRepository.Create(c, user.ID, form.ClientID, refreshToken, authzCodeCtx.Scope); err != nil {
			return response.InternalServerError(c, err)
		}

		return c.JSON(http.StatusOK, dto.AuthzCodeResp{
			AccessToken:  accessToken,
			TokenType:    "Bearer",
			ExpiresIn:    int(expiresIn.Seconds()),
			RefreshToken: refreshToken,
			IDToken:      idToken,
		})

	case "refresh_token":
		refreshTokenData, err := h.refreshTokenRepository.FindByToken(c, form.RefreshToken)
		if err != nil {
			return response.NotFound(c, err)
		}

		exists, err := h.userRepository.ExistsByID(c, refreshTokenData.UserID)
		if err != nil {
			return response.InternalServerError(c, err)
		}
		if !exists {
			err = errors.New("user not found")
			return response.NotFound(c, err)
		}

		accessToken, err := h.jwtService.GenerateAccessToken(&refreshTokenData.UserID, oauthClient, expiresIn)
		if err != nil {
			return response.InternalServerError(c, err)
		}

		newRefreshToken, err := h.jwtService.GenerateRefreshToken()
		if err != nil {
			return response.InternalServerError(c, err)
		}

		if err := h.refreshTokenRepository.Update(c, refreshTokenData.UserID, form.ClientID, newRefreshToken, refreshTokenData.Scopes); err != nil {
			return response.InternalServerError(c, err)
		}

		return c.JSON(http.StatusOK, dto.RefreshTokenResp{
			AccessToken:  accessToken,
			TokenType:    "Bearer",
			ExpiresIn:    int(expiresIn.Seconds()),
			RefreshToken: newRefreshToken,
		})

	case "client_credentials":
		accessToken, err := h.jwtService.GenerateAccessToken(nil, oauthClient, expiresIn)
		if err != nil {
			return response.InternalServerError(c, err)
		}

		return c.JSON(http.StatusOK, dto.ClientCredentialsResponse{
			AccessToken: accessToken,
			TokenType:   "Bearer",
			ExpiresIn:   int(expiresIn.Seconds()),
		})

	default:
		err = errors.New("unsupported grant_type")
		return response.BadRequest(c, err)
	}
}
