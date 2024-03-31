package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
)

type ErrHandler interface {
	GetView(c echo.Context) error
}

type errHandler struct{}

func NewErrHandler() *errHandler {
	return &errHandler{}
}

func (h *errHandler) GetView(c echo.Context) error {
	return c.File(config.ERROR_VIEW_FILEPATH)
}
