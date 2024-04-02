package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/internal/resource/dto"
	"github.com/ucho456job/lgtmeme/internal/resource/repository"
	"github.com/ucho456job/lgtmeme/internal/util/response"
)

type ImageHandler interface {
	GetImages(c echo.Context) error
}

type imageHandler struct {
	imageRepository repository.ImageRepository
}

func NewImageHandler(imageRepository repository.ImageRepository) ImageHandler {
	return &imageHandler{
		imageRepository: imageRepository,
	}
}

func (h *imageHandler) GetImages(c echo.Context) error {
	q := new(dto.GetImagesQuery)
	if err := c.Bind(q); err != nil {
		return response.BadRequest(c, err)
	}

	if err := c.Validate(q); err != nil {
		return response.BadRequest(c, err)
	}

	imgs, err := h.imageRepository.FindImages(c, *q)
	if err != nil {
		return response.InternalServerError(c, err)
	}

	total := len(*imgs)
	respImgs := make([]dto.GetImagesImages, total)
	for i, img := range *imgs {
		respImgs[i] = dto.GetImagesImages{
			ID:  img.ID,
			URL: img.URL,
		}
	}

	return c.JSON(http.StatusOK, dto.GetImagesResp{
		Total:  total,
		Images: respImgs,
	})
}
