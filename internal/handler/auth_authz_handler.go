package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/dto"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/internal/util/response"
)

type AuthzHandler interface {
	Authorize(c echo.Context) error
}

type authzHandler struct {
	oauthClientRepository    repository.OauthClientRepository
	userRepository           repository.UserRepository
	scopeRepository          repository.ScopeRepository
	sessionManagerRepository repository.SessionManagerRepository
}

func NewAuthzHandler(
	oauthClientRepository repository.OauthClientRepository,
	userRepository repository.UserRepository,
	scopeRepository repository.ScopeRepository,
	sessionManagerRepository repository.SessionManagerRepository,
) *authzHandler {
	return &authzHandler{
		oauthClientRepository:    oauthClientRepository,
		userRepository:           userRepository,
		scopeRepository:          scopeRepository,
		sessionManagerRepository: sessionManagerRepository,
	}
}

func (h *authzHandler) Authorize(c echo.Context) error {
	clientID, err := uuid.Parse(c.QueryParam("client_id"))
	if err != nil {
		return response.BadRequest(c, err)
	}

	q := &dto.AuthzQuery{
		ResponseType: c.QueryParam("response_type"),
		ClientID:     clientID,
		RedirectURI:  c.QueryParam("redirect_uri"),
		Scope:        c.QueryParam("scope"),
		State:        c.QueryParam("state"),
		Nonce:        c.QueryParam("nonce"),
	}

	if err := c.Validate(q); err != nil {
		return response.BadRequest(c, err)
	}

	valid, err := h.oauthClientRepository.IsValidOAuthClient(c, *q)
	if err != nil {
		return response.InternalServerError(c, err)
	}
	if !valid {
		err = errors.New("client ID or redirect URI or scope are incorrect")
		return response.BadRequest(c, err)
	}

	userID, isLogin, err := h.sessionManagerRepository.LoadLoginSession(c)
	if err != nil {
		return response.InternalServerError(c, err)
	}
	if !isLogin {
		if err := h.sessionManagerRepository.CachePreAuthnSession(c, *q); err != nil {
			return response.InternalServerError(c, err)
		}

		scopes, err := h.scopeRepository.FindByScopesStr(c, q.Scope)
		if err != nil {
			return response.InternalServerError(c, err)
		}

		scopesStr := ""
		descriptionsStr := ""
		for i, s := range *scopes {
			if i != 0 {
				scopesStr += ","
				descriptionsStr += ","
			}
			scopesStr += s.Code
			descriptionsStr += s.Description
		}

		redirectURL := fmt.Sprintf("%s?scopes=%s&descriptions=%s", config.LOGIN_VIEW_ENDPOINT, scopesStr, descriptionsStr)
		return c.Redirect(http.StatusFound, redirectURL)
	}

	exists, err := h.userRepository.ExistsByID(c, userID)
	if err != nil {
		return response.InternalServerError(c, err)
	}
	if !exists {
		err = errors.New("user not found")
		return response.BadRequest(c, err)
	}

	authzCode := uuid.New().String()
	if err := h.sessionManagerRepository.CacheAuthzCodeCtx(c, *q, authzCode, userID); err != nil {
		return response.InternalServerError(c, err)
	}

	redirectURL, err := url.Parse(q.RedirectURI)
	if err != nil {
		return response.InternalServerError(c, err)
	}
	query := redirectURL.Query()
	query.Set("code", authzCode)
	if q.State != "" {
		query.Set("state", q.State)
	}
	redirectURL.RawQuery = query.Encode()

	return c.Redirect(http.StatusFound, redirectURL.String())
}
