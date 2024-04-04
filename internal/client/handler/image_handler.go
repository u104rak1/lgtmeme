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

type ImageHandler interface {
	GetCreateImageView(c echo.Context) error
	Post(c echo.Context) error
	BulkGet(c echo.Context) error
	Patch(c echo.Context) error
}

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

func (h *imageHandler) GetCreateImageView(c echo.Context) error {
	return c.File(config.CREATE_IMAGE_VIEW_FILEPATH)
}

func (h *imageHandler) Post(c echo.Context) error {
	var body resourceDto.PostImageReqBody
	if err := c.Bind(&body); err != nil {
		return response.BadRequest(c, err)
	}

	if err := c.Validate(body); err != nil {
		return response.BadRequest(c, err)
	}

	accessToken, err := h.sessionManagerRepository.LoadToken(c, config.GENERAL_ACCESS_TOKEN_SESSION_NAME)
	if err != nil {
		return response.InternalServerError(c, err)
	}

	respBody, status, err := h.imageService.Post(c, body, accessToken)
	if err != nil {
		return response.HandleErrResp(c, status, err)
	}

	return c.JSON(http.StatusCreated, respBody)
}

func (h *imageHandler) BulkGet(c echo.Context) error {
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

	respBody, status, err := h.imageService.BulkGet(c, *q, accessToken)
	if err != nil {
		return response.HandleErrResp(c, status, err)
	}

	return c.JSON(http.StatusOK, respBody)
}

func (h *imageHandler) Patch(c echo.Context) error {
	var body resourceDto.PatchImageReqBody
	if err := c.Bind(&body); err != nil {
		return response.BadRequest(c, err)
	}

	if err := c.Validate(body); err != nil {
		return response.BadRequest(c, err)
	}

	accessToken, err := h.sessionManagerRepository.LoadToken(c, config.GENERAL_ACCESS_TOKEN_SESSION_NAME)
	if err != nil {
		return response.InternalServerError(c, err)
	}

	imageID := c.Param("image_id")

	status, err := h.imageService.Patch(c, body, imageID, accessToken)
	if err != nil {
		return response.HandleErrResp(c, status, err)
	}

	return c.JSON(http.StatusNoContent, nil)
}
