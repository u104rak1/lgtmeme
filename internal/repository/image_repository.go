package repository

// mockgen -source=internal/repository/image_repository.go -destination=test/mock/repository/mock_image_repository.go -package=repository_mock

import (
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/config"
	"github.com/ucho456job/lgtmeme/internal/dto"
	"github.com/ucho456job/lgtmeme/internal/model"
	"github.com/ucho456job/lgtmeme/internal/util/timer"
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
	DB    *gorm.DB
	Timer timer.Timer
}

func NewImageRepository(db *gorm.DB, timer timer.Timer) ImageRepository {
	return &imageRepository{
		DB:    db,
		Timer: timer,
	}
}

func (r *imageRepository) Create(c echo.Context, id uuid.UUID, url, keyword string) error {
	newImage := &model.Image{
		ID:        id,
		URL:       url,
		Keyword:   keyword,
		UsedCount: 0,
		Reported:  false,
		Confirmed: false,
		CreatedAt: r.Timer.Now(),
	}

	return r.DB.Debug().Create(newImage).Error
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
	if err := sqlQ.Offset(q.Page * config.GET_IMAGES_LIMIT).Limit(config.GET_IMAGES_LIMIT).Find(&images).Error; err != nil {
		return nil, err
	}

	return &images, nil
}

func (r *imageRepository) FindURLByID(c echo.Context, id uuid.UUID) (*string, error) {
	var url string
	if err := r.DB.Model(&model.Image{}).Select("url").Where("id = ?", id).First(&url).Error; err != nil {
		return nil, err
	}

	return &url, nil
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

	return r.DB.Model(&model.Image{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *imageRepository) Delete(c echo.Context, id uuid.UUID) error {
	return r.DB.Where("id = ?", id).Delete(&model.Image{}).Error
}
