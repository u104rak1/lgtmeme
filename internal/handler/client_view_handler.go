package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
)

type ViewHandler interface {
	GetErrView(c echo.Context) error
}

type viewHandler struct{}

func NewViewHandler() *viewHandler {
	return &viewHandler{}
}

func (h *viewHandler) GetErrView(c echo.Context) error {
	return c.File(config.ERROR_VIEW_FILEPATH)
}

func (h *viewHandler) GetPrivacyPolicyView(c echo.Context) error {
	return c.File(config.PRIVACY_POLICY_FILEPATH)
}

func (h *viewHandler) GetTermsOfServiceView(c echo.Context) error {
	return c.File(config.TERMS_OF_SERVICE_FILEPATH)
}
