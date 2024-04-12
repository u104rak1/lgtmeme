package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/dto"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/internal/service"
	"github.com/ucho456job/lgtmeme/internal/util/response"
)

type ClientImageHandler interface {
	Post(c echo.Context) error
	BulkGet(c echo.Context) error
	Patch(c echo.Context) error
	Delete(c echo.Context) error
}

type clientImageHandler struct {
	sessionManagerRepository repository.SessionManager
	accessTokenService       service.AccessTokenService
	imageService             service.ImageService
}

func NewClientImageHandler(
	sessionManagerRepository repository.SessionManager,
	accessTokenService service.AccessTokenService,
	imageService service.ImageService,
) *clientImageHandler {
	return &clientImageHandler{
		sessionManagerRepository: sessionManagerRepository,
		accessTokenService:       accessTokenService,
		imageService:             imageService,
	}
}

func (h *clientImageHandler) Post(c echo.Context) error {
	var body dto.PostImageReqBody
	if err := c.Bind(&body); err != nil {
		return response.BadRequest(c, err)
	}

	if err := c.Validate(body); err != nil {
		return response.BadRequest(c, err)
	}

	accessToken, err := h.sessionManagerRepository.LoadGeneralAccessToken(c)
	if err != nil {
		return response.InternalServerError(c, err)
	}

	respBody, status, err := h.imageService.Post(c, body, accessToken)
	if err != nil {
		return response.HandleErrResp(c, status, err)
	}

	return c.JSON(http.StatusCreated, respBody)
}

func (h *clientImageHandler) BulkGet(c echo.Context) error {
	q := new(dto.GetImagesQuery)
	if err := c.Bind(q); err != nil {
		return response.BadRequest(c, err)
	}

	if err := c.Validate(q); err != nil {
		return response.BadRequest(c, err)
	}

	accessToken, err := h.sessionManagerRepository.LoadGeneralAccessToken(c)
	if err != nil {
		return response.InternalServerError(c, err)
	}

	respBody, status, err := h.imageService.BulkGet(c, *q, accessToken)
	if err != nil {
		return response.HandleErrResp(c, status, err)
	}

	return c.JSON(http.StatusOK, respBody)
}

func (h *clientImageHandler) Patch(c echo.Context) error {
	var body dto.PatchImageReqBody
	if err := c.Bind(&body); err != nil {
		return response.BadRequest(c, err)
	}

	if err := c.Validate(body); err != nil {
		return response.BadRequest(c, err)
	}

	accessToken, err := h.sessionManagerRepository.LoadGeneralAccessToken(c)
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

func (h *clientImageHandler) Delete(c echo.Context) error {
	accessToken, err := h.sessionManagerRepository.LoadToken(c, config.ADMIN_ACCESS_TOKEN_SESSION_NAME)
	if err != nil {
		return response.InternalServerError(c, err)
	}

	imageID := c.Param("image_id")

	status, err := h.imageService.Delete(c, imageID, accessToken)
	if err != nil {
		return response.HandleErrResp(c, status, err)
	}

	return c.JSON(http.StatusNoContent, nil)
}
