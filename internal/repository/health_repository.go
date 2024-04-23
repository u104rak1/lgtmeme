package repository

// mockgen -source=internal/repository/health_repository.go -destination=test/mock/repository/mock_health_repository.go -package=repository_mock

import (
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/internal/model"
	"gorm.io/gorm"
)

type HealthRepository interface {
	CheckPostgres(c echo.Context, key string) (value string, err error)
}

type healthRepository struct {
	DB *gorm.DB
}

func NewHealthRepository(db *gorm.DB) HealthRepository {
	return &healthRepository{DB: db}
}

func (r *healthRepository) CheckPostgres(c echo.Context, key string) (value string, err error) {
	if err := r.DB.Model(&model.HealthCheck{}).Select("value").Where("key = ?", key).First(&value).Error; err != nil {
		return "", err
	}
	return value, nil
}
