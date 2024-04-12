package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/internal/dto"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/internal/util/response"
)

type LogoutHandler interface {
	Logout(c echo.Context) error
}

type logoutHandler struct {
	sessionManager repository.SessionManagerRepository
}

func NewLogoutHandler(sessionManager repository.SessionManagerRepository) *logoutHandler {
	return &logoutHandler{
		sessionManager: sessionManager,
	}
}

func (h *logoutHandler) Logout(c echo.Context) error {
	if err := h.sessionManager.Logout(c); err != nil {
		return response.InternalServerError(c, err)
	}
	return c.JSON(http.StatusOK, dto.LogoutResp{OK: true})
}
