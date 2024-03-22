package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/my_authn_authz/internal/dto"
	"github.com/ucho456job/my_authn_authz/internal/repository"
	"github.com/ucho456job/my_authn_authz/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type LoginHandler interface {
	Login(c echo.Context)
}

type loginHandler struct {
	userRepository repository.UserRepository
	sessionManager repository.SessionManager
}

func NewLoginHandler(
	userRepository repository.UserRepository,
	sessionManager repository.SessionManager,
) *loginHandler {
	return &loginHandler{
		userRepository: userRepository,
		sessionManager: sessionManager,
	}
}

func (h *loginHandler) Login(c echo.Context) error {
	var form dto.LoginForm
	if err := c.Bind(&form); err != nil {
		return util.InternalServerErrorResponse(c, err)
	}
	if err := c.Validate(&form); err != nil {
		return util.BadRequestResponse(c, err)
	}

	user, err := h.userRepository.FindByName(c, form.Username)
	if err != nil {
		return util.UnauthorizedErrorResponse(c, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		return util.UnauthorizedErrorResponse(c, err)
	}

	if err := h.sessionManager.CacheLoginSession(c, user.ID); err != nil {
		return util.InternalServerErrorResponse(c, err)
	}

	query, exists, err := h.sessionManager.LoadPreAuthnSession(c)
	if err != nil {
		return util.InternalServerErrorResponse(c, err)
	}
	if !exists {
		return c.JSON(http.StatusOK, map[string]string{"redirectURL": util.PASSKEY_SCREEN_ENDPOINT})
	}

	queryParams := url.Values{}
	queryParams.Set("response_type", query.ResponseType)
	queryParams.Set("client_id", query.ClientID.String())
	queryParams.Set("redirect_uri", query.RedirectURI)
	if query.Scope != "" {
		queryParams.Set("scope", query.Scope)
	}
	if query.State != "" {
		queryParams.Set("state", query.State)
	}
	if query.Nonce != "" {
		queryParams.Set("nonce", query.Nonce)
	}
	authorizationURL := fmt.Sprintf("%s?%s", "/api/connect/authorize", queryParams.Encode())

	return c.JSON(http.StatusOK, map[string]string{"redirectURL": authorizationURL})
}
