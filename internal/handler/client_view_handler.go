package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
)

type ViewHandler interface {
	GetHomeView(c echo.Context, filePath string) error
	GetImageView(c echo.Context) error
	GetPrivacyPolicyView(c echo.Context) error
	GetTermsOfServiceView(c echo.Context) error
	GetAdminView(c echo.Context) error
}

type viewHandler struct{}

func NewViewHandler() *viewHandler {
	return &viewHandler{}
}

func (h *viewHandler) GetHomeView(c echo.Context) error {
	return c.File(config.HOME_VIEW_FILEPATH)
}

func (h *viewHandler) GetImageView(c echo.Context) error {
	return c.File(config.IMAGE_NEW_VIEW_FILEPATH)
}

func (h *viewHandler) GetPrivacyPolicyView(c echo.Context) error {
	return c.File(config.PRIVACY_POLICY_FILEPATH)
}

func (h *viewHandler) GetTermsOfServiceView(c echo.Context) error {
	return c.File(config.TERMS_OF_SERVICE_FILEPATH)
}

func (h *viewHandler) GetAdminView(c echo.Context) error {
	return c.File(config.ADMIN_VIEW_FILEPATH)
}
