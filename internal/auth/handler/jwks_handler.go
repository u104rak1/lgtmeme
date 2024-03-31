package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/internal/auth/service"
)

type JwksHandler interface {
	Get(c echo.Context) error
}

type jwksHandler struct {
	jwtService service.JwtService
}

func NewJwksHandler(jwtService service.JwtService) JwksHandler {
	return &jwksHandler{
		jwtService: jwtService,
	}
}

func (h *jwksHandler) Get(c echo.Context) error {
	jwks, err := h.jwtService.GetPublicKeys()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, jwks)
}
