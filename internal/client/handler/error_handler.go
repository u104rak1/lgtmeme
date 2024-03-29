package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
)

type ErrorViewHandler interface {
	GetErrorView(c echo.Context) error
}

type errorViewHandler struct{}

func NewErrorViewHandler() *errorViewHandler {
	return &errorViewHandler{}
}

func (h *errorViewHandler) GetErrorView(c echo.Context) error {
	return c.File(config.ERROR_VIEW_FILEPATH)
}
