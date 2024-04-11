package handler

import (
	"fmt"
	"net/http"

	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/dto"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/internal/service"
	"github.com/ucho456job/lgtmeme/internal/util/response"
)

type ClientImageHandler interface {
	GetView(c echo.Context) error
	Post(c echo.Context) error
	BulkGet(c echo.Context) error
	Patch(c echo.Context) error
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

func (h *clientImageHandler) GetView(c echo.Context) error {
	_, err := h.sessionManagerRepository.LoadGeneralAccessToken(c)
	if err == redis.ErrNil {
		respBody, status, err := h.accessTokenService.CallTokenWithClientCredentials(c)
		if err != nil && status != http.StatusOK {
			errURL := fmt.Sprintf("%s?code=%d", config.ERROR_VIEW_ENDPOINT, status)
			return c.Redirect(http.StatusFound, errURL)
		}

		if err := h.sessionManagerRepository.CacheGeneralAccessToken(c, respBody.AccessToken); err != nil {
			return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
		}
	} else if err != nil {
		return c.Redirect(http.StatusFound, config.ERROR_VIEW_ENDPOINT)
	}

	return c.File(config.IMAGE_NEW_VIEW_FILEPATH)
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
