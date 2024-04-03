package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/client/repository"
	"github.com/ucho456job/lgtmeme/internal/client/service"
	resourceDto "github.com/ucho456job/lgtmeme/internal/resource/dto"
	"github.com/ucho456job/lgtmeme/internal/util/response"
)

type ImageHandler interface{}

type imageHandler struct {
	sessionManagerRepository repository.SessionManager
	imageService             service.ImageService
}

func NewImageHandler(
	sessionManagerRepository repository.SessionManager,
	imageService service.ImageService,
) *imageHandler {
	return &imageHandler{
		sessionManagerRepository: sessionManagerRepository,
		imageService:             imageService,
	}
}

func (h *imageHandler) Get(c echo.Context) error {
	q := new(resourceDto.GetImagesQuery)
	if err := c.Bind(q); err != nil {
		return response.BadRequest(c, err)
	}

	if err := c.Validate(q); err != nil {
		return response.BadRequest(c, err)
	}

	accessToken, err := h.sessionManagerRepository.LoadToken(c, config.GENERAL_ACCESS_TOKEN_SESSION_NAME)
	if err != nil {
		return response.InternalServerError(c, err)
	}

	respBody, status, err := h.imageService.GetImages(c, *q, accessToken)
	if err != nil {
		return response.HandleErrResp(c, status, err)
	}

	return c.JSON(http.StatusOK, respBody)
}
