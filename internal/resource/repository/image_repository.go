package repository

import (
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/internal/resource/dto"
	"github.com/ucho456job/lgtmeme/internal/resource/model"
	"github.com/ucho456job/lgtmeme/internal/util/clock"
	"gorm.io/gorm"
)

type ImageRepository interface {
	Create(c echo.Context, id uuid.UUID, url, keyword string) error
	FindImages(c echo.Context, q dto.GetImagesQuery) (*[]model.Image, error)
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
	sqlQ := r.DB.Model(&model.Image{}).Select("id", "url")

	if q.FavoriteImageIDs != "" {
		favoriteImageIDs := strings.Split(q.FavoriteImageIDs, ",")
		sqlQ = sqlQ.Where("id IN ?", favoriteImageIDs)
	}

	if q.Keyword != "" {
		sqlQ = sqlQ.Where("keyword LIKE ?", "%"+q.Keyword+"%")
	}

	if q.AuthCheck {
		sqlQ = sqlQ.Where("confirmed = ?", false)
		sqlQ = sqlQ.Where("reported = ?", true)
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
