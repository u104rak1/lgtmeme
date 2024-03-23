package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/my_authn_authz/internal/util"
)

type JwksHandler interface {
	GetJwks(c echo.Context) error
}

type jwksHandler struct {
	jwtService util.JwtService
}

func NewJwksHandler(jwtService util.JwtService) JwksHandler {
	return &jwksHandler{
		jwtService: jwtService,
	}
}

func (h *jwksHandler) GetJwks(c echo.Context) error {
	jwks, err := h.jwtService.GetPublicKeys()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, jwks)
}
