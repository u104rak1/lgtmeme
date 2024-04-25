package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/internal/dto"
	"github.com/ucho456job/lgtmeme/internal/repository"
	"github.com/ucho456job/lgtmeme/internal/service"
	"github.com/ucho456job/lgtmeme/internal/util/response"
	"github.com/ucho456job/lgtmeme/internal/util/uuidgen"
)

type ResourceImageHandler interface {
	Post(c echo.Context) error
	BulkGet(c echo.Context) error
	Patch(c echo.Context) error
	Delete(c echo.Context) error
}

type resourceImageHandler struct {
	imageRepository repository.ImageRepository
	userRepository  repository.UserRepository
	storageService  service.StorageService
	uuidGenerator   uuidgen.UUIDGenerator
}

func NewResourceImageHandler(
	imageRepository repository.ImageRepository,
	userRepository repository.UserRepository,
	storageService service.StorageService,
	uuidGenerator uuidgen.UUIDGenerator,
) ResourceImageHandler {
	return &resourceImageHandler{
		imageRepository: imageRepository,
		userRepository:  userRepository,
		storageService:  storageService,
		uuidGenerator:   uuidGenerator,
	}
}

func (h *resourceImageHandler) Post(c echo.Context) error {
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

func (h *resourceImageHandler) BulkGet(c echo.Context) error {
	q := new(dto.GetImagesQuery)
	if err := c.Bind(q); err != nil {
		return response.BadRequest(c, err)
	}

	if err := c.Validate(q); err != nil {
		return response.BadRequest(c, err)
	}

	imgs, err := h.imageRepository.FindByGetImagesQuery(c, *q)
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

func (h *resourceImageHandler) Patch(c echo.Context) error {
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

func (h *resourceImageHandler) Delete(c echo.Context) error {
	imageIDStr := c.Param("image_id")
	imageID, err := uuid.Parse(imageIDStr)
	if err != nil {
		return response.BadRequest(c, err)
	}

	tmpUserID := c.Get("userID").(string)
	userID, err := uuid.Parse(tmpUserID)
	if err != nil {
		return response.InternalServerError(c, err)
	}

	columns := []string{"role"}
	user, err := h.userRepository.FirstByID(c, userID, columns)
	if err != nil {
		return response.InternalServerError(c, err)
	}

	if user.Role != "admin" {
		return response.Forbidden(c, errors.New("permission denied"))
	}

	columns = []string{"url"}
	img, err := h.imageRepository.FirstByID(c, imageID, columns)
	if img == nil {
		return response.NotFound(c, errors.New("image not found"))
	}
	if err != nil {
		return response.InternalServerError(c, err)
	}

	if err := h.imageRepository.Delete(c, imageID); err != nil {
		return response.InternalServerError(c, err)
	}

	if err := h.storageService.Delete(c, img.URL); err != nil {
		return response.InternalServerError(c, err)
	}

	return c.JSON(http.StatusNoContent, nil)
}
