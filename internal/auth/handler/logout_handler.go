package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/internal/auth/dto"
	"github.com/ucho456job/lgtmeme/internal/auth/repository"
	"github.com/ucho456job/lgtmeme/internal/util"
)

type LogoutHandler interface {
	Logout(c echo.Context) error
}

type logoutHandler struct {
	sessionManager repository.SessionManager
}

func NewLogoutHandler(sessionManager repository.SessionManager) *logoutHandler {
	return &logoutHandler{
		sessionManager: sessionManager,
	}
}

func (h *logoutHandler) Logout(c echo.Context) error {
	if err := h.sessionManager.Logout(c); err != nil {
		return util.InternalServerErrorResponse(c, err)
	}
	return c.JSON(http.StatusOK, dto.LogoutResp{OK: true})
}
