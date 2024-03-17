package repository

import (
	"github.com/ucho456job/my_authn_authz/internal/model"
	"gorm.io/gorm"
)

type HealthCheckRepository interface {
	CheckHealthForPostgres(key string) (value string, err error)
}

type healthCheckRepository struct {
	DB *gorm.DB
}

func NewHealthCheckRepository(db *gorm.DB) HealthCheckRepository {
	return &healthCheckRepository{DB: db}
}

func (r *healthCheckRepository) CheckHealthForPostgres(key string) (value string, err error) {
	if err := r.DB.Model(&model.HealthCheck{}).Select("value").Where("key = ?", key).First(&value).Error; err != nil {
		return "", err
	}
	return value, nil
}
