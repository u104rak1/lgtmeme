package handler

import (
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/my_authn_authz/internal/dto"
	"github.com/ucho456job/my_authn_authz/internal/model"
	"github.com/ucho456job/my_authn_authz/internal/repository"
	"github.com/ucho456job/my_authn_authz/internal/util"
)

type AuthorizationHandler interface {
	AuthorizationHandle(c echo.Context) error
}

type authorizationHandler struct {
	oauthClientRepository    repository.OauthClientRepository
	userRepository           repository.UserRepository
	sessionManagerRepository repository.SessionManager
}

func NewAuthorizationHandler(
	oauthClientRepository repository.OauthClientRepository,
	userRepository repository.UserRepository,
	sessionManagerRepository repository.SessionManager,
) *authorizationHandler {
	return &authorizationHandler{
		oauthClientRepository:    oauthClientRepository,
		userRepository:           userRepository,
		sessionManagerRepository: sessionManagerRepository,
	}
}

func (h *authorizationHandler) AuthorizationHandle(c echo.Context) error {
	clientID, err := model.ParseClientID(c.QueryParam("client_id"))
	if err != nil {
		return util.BadRequestResponse(c, err)
	}

	q := &dto.AuthorizationQuery{
		ResponseType: c.QueryParam("response_type"),
		ClientID:     clientID,
		RedirectURI:  c.QueryParam("redirect_uri"),
		Scope:        c.QueryParam("scope"),
		State:        c.QueryParam("state"),
		Nonce:        c.QueryParam("nonce"),
	}

	if err := c.Validate(q); err != nil {
		return util.BadRequestResponse(c, err)
	}

	exists, err := h.oauthClientRepository.ExistsForAuthz(c, *q)
	if err != nil {
		return util.RedirectWithErrorForAuthz(c, *q, "server_error", "An internal server error occurred")
	}
	if !exists {
		return util.RedirectWithErrorForAuthz(c, *q, "invalid_request", "Client ID or Redirect URI or scope are incorrect")
	}

	userID, isLogin, err := h.sessionManagerRepository.LoadLoginSession(c)
	if err != nil {
		return util.RedirectWithErrorForAuthz(c, *q, "server_error", "Failed to get login session")
	}
	if !isLogin {
		if err := h.sessionManagerRepository.CachePreAuthnSession(c, *q); err != nil {
			return util.RedirectWithErrorForAuthz(c, *q, "server_error", "Failed to save pre authentication session")
		}
		return c.Redirect(http.StatusFound, util.LOGIN_SCREEN_ENDPOINT)
	}

	exists, err = h.userRepository.ExistsByID(c, userID)
	if err != nil || !exists {
		return util.RedirectWithErrorForAuthz(c, *q, "access_denied", "User does not exist")
	}

	authzCode := uuid.New().String()
	if err := h.sessionManagerRepository.CacheAuthzCodeWithCtx(c, *q, userID, authzCode); err != nil {
		return util.RedirectWithErrorForAuthz(c, *q, "server_error", "Failed to save authorization code")
	}

	redirectURL, err := url.Parse(q.RedirectURI)
	if err != nil {
		return util.RedirectWithErrorForAuthz(c, *q, "server_error", "Failed to parse redirect URI")
	}
	query := redirectURL.Query()
	query.Set("code", authzCode)
	if q.State != "" {
		query.Set("state", q.State)
	}
	redirectURL.RawQuery = query.Encode()

	return c.Redirect(http.StatusFound, redirectURL.String())
}
