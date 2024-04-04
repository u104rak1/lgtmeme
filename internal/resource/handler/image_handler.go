package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/internal/resource/dto"
	"github.com/ucho456job/lgtmeme/internal/resource/repository"
	"github.com/ucho456job/lgtmeme/internal/resource/service"
	"github.com/ucho456job/lgtmeme/internal/util/response"
	"github.com/ucho456job/lgtmeme/internal/util/uuidgen"
)

type ImageHandler interface {
	Post(c echo.Context) error
	BulkGet(c echo.Context) error
	Patch(c echo.Context) error
}

type imageHandler struct {
	imageRepository repository.ImageRepository
	storageService  service.StorageService
	uuidGenerator   uuidgen.UUIDGenerator
}

func NewImageHandler(
	imageRepository repository.ImageRepository,
	storageService service.StorageService,
	uuidGenerator uuidgen.UUIDGenerator,
) ImageHandler {
	return &imageHandler{
		imageRepository: imageRepository,
		storageService:  storageService,
		uuidGenerator:   uuidGenerator,
	}
}

func (h *imageHandler) Post(c echo.Context) error {
	var body dto.PostImageReqBody
	if err := c.Bind(&body); err != nil {
		return response.BadRequest(c, err)
	}

	if err := c.Validate(body); err != nil {
		return response.BadRequest(c, err)
	}

	imgID := h.uuidGenerator.New()

	imgURL, err := h.storageService.Upload(c, imgID.String(), body.Base64image)
	if err != nil {
		return response.InternalServerError(c, err)
	}

	if err := h.imageRepository.Create(c, imgID, imgURL, body.Keyword); err != nil {
		if err := h.storageService.Delete(c, imgURL); err != nil {
			return response.InternalServerError(c, err)
		}
		return response.InternalServerError(c, err)
	}

	return c.JSON(http.StatusCreated, dto.PostImageResp{
		ImageURL: imgURL,
	})
}

func (h *imageHandler) BulkGet(c echo.Context) error {
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

func (h *imageHandler) Patch(c echo.Context) error {
	var body dto.PatchImageReqBody
	if err := c.Bind(&body); err != nil {
		return response.BadRequest(c, err)
	}

	if err := c.Validate(body); err != nil {
		return response.BadRequest(c, err)
	}

	imageIDStr := c.Param("image_id")
	imageID, err := uuid.Parse(imageIDStr)
	if err != nil {
		return response.BadRequest(c, err)
	}

	exists, err := h.imageRepository.ExistsByID(c, imageID)
	if err != nil {
		return response.InternalServerError(c, err)
	}
	if !exists {
		return response.NotFound(c, nil)
	}

	if err := h.imageRepository.Update(c, imageID, body.Type); err != nil {
		return response.InternalServerError(c, err)
	}

	return c.JSON(http.StatusNoContent, nil)
}
