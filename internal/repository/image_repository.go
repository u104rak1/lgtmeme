package repository

import (
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/internal/dto"
	"github.com/ucho456job/lgtmeme/internal/model"
	"github.com/ucho456job/lgtmeme/internal/util/clock"
	"gorm.io/gorm"
)

type ImageRepository interface {
	Create(c echo.Context, id uuid.UUID, url, keyword string) error
	FindImages(c echo.Context, q dto.GetImagesQuery) (*[]model.Image, error)
	FindURLByID(c echo.Context, id uuid.UUID) (*string, error)
	ExistsByID(c echo.Context, id uuid.UUID) (bool, error)
	Update(c echo.Context, id uuid.UUID, reqType dto.PatchImageReqType) error
	Delete(c echo.Context, id uuid.UUID) error
}

type imageRepository struct {
	DB      *gorm.DB
	Clocker clock.Clocker
}

func NewImageRepository(db *gorm.DB, clocker clock.Clocker) ImageRepository {
	return &imageRepository{
		DB:      db,
		Clocker: clocker,
	}
}

func (r *imageRepository) Create(c echo.Context, id uuid.UUID, url, keyword string) error {
	newImage := &model.Image{
		ID:        id,
		URL:       url,
		Keyword:   keyword,
		CreatedAt: r.Clocker.Now(),
	}

	if err := r.DB.Create(newImage).Error; err != nil {
		return err
	}

	return nil
}

func (r *imageRepository) FindImages(c echo.Context, q dto.GetImagesQuery) (*[]model.Image, error) {
	sqlQ := r.DB.Debug().Model(&model.Image{})

	if q.FavoriteImageIDs != "" {
		favoriteImageIDs := strings.Split(q.FavoriteImageIDs, ",")
		sqlQ = sqlQ.Where("id IN ?", favoriteImageIDs)
	}

	if q.Keyword != "" {
		sqlQ = sqlQ.Where("keyword LIKE ?", "%"+q.Keyword+"%")
	}

	if q.AuthCheck {
		sqlQ = sqlQ.Where("confirmed = ?", false).Where("reported = ?", true)
	} else {
		sqlQ = sqlQ.Where("confirmed = ? OR reported = ?", true, false)
	}

	if q.Sort == "latest" {
		sqlQ = sqlQ.Order("created_at DESC")
	} else {
		sqlQ = sqlQ.Order("used_count DESC, created_at DESC")
	}

	var images []model.Image
	if err := sqlQ.Offset(q.Page * 9).Limit(9).Find(&images).Error; err != nil {
		return nil, err
	}

	return &images, nil
}

func (r *imageRepository) FindURLByID(c echo.Context, id uuid.UUID) (*string, error) {
	var image model.Image
	if err := r.DB.Model(&model.Image{}).Where("id = ?", id).First(&image).Error; err != nil {
		return nil, err
	}

	return &image.URL, nil
}

func (r *imageRepository) ExistsByID(c echo.Context, id uuid.UUID) (bool, error) {
	var count int64
	if err := r.DB.Model(&model.Image{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *imageRepository) Update(c echo.Context, id uuid.UUID, reqType dto.PatchImageReqType) error {
	var updateData map[string]interface{}
	switch reqType {
	case dto.PatchImageReqTypeUsed:
		updateData = map[string]interface{}{"used_count": gorm.Expr("used_count + ?", 1)}
	case dto.PatchImageReqTypeReport:
		updateData = map[string]interface{}{"reported": true}
	case dto.PatchImageReqTypeConfirm:
		updateData = map[string]interface{}{"confirmed": true}
	}

	if err := r.DB.Model(&model.Image{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
		return err
	}

	return nil
}

func (r *imageRepository) Delete(c echo.Context, id uuid.UUID) error {
	if err := r.DB.Where("id = ?", id).Delete(&model.Image{}).Error; err != nil {
		return err
	}

	return nil
}
