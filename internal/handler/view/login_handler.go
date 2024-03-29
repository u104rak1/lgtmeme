package view_handler

import (
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
)

type LoginViewHandler interface {
	GetLoginView(c echo.Context) error
}

type loginViewHandler struct{}

func NewLoginViewHandler() *loginViewHandler {
	return &loginViewHandler{}
}

func (h *loginViewHandler) GetLoginView(c echo.Context) error {
	return c.File(config.LOGIN_VIEW_FILEPATH)
}
