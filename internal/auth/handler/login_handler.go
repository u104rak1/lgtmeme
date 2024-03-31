package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/auth/dto"
	"github.com/ucho456job/lgtmeme/internal/auth/repository"
	"github.com/ucho456job/lgtmeme/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type LoginHandler interface {
	Login(c echo.Context) error
	GetView(c echo.Context) error
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
		// TODO: Think about redirect destinations
		return c.JSON(http.StatusOK, dto.LoginResp{RedirectURL: config.PASSKEY_VIEW_ENDPOINT})
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
	authorizationURL := fmt.Sprintf("%s?%s", config.AUTHZ_ENDPOINT, queryParams.Encode())

	return c.JSON(http.StatusOK, dto.LoginResp{RedirectURL: authorizationURL})
}

func (h *loginHandler) GetView(c echo.Context) error {
	return c.File(config.LOGIN_VIEW_FILEPATH)
}
