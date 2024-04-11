package repository

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(c echo.Context, userID uuid.UUID) (*model.User, error)
	FindByName(c echo.Context, name string) (*model.User, error)
	ExistsByID(c echo.Context, userID uuid.UUID) (bool, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) FindByID(c echo.Context, userID uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.DB.Model(&model.User{}).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByName(c echo.Context, name string) (*model.User, error) {
	var user model.User
	if err := r.DB.Model(&model.User{}).Where("name = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) ExistsByID(c echo.Context, userID uuid.UUID) (bool, error) {
	var count int64
	if err := r.DB.Model(&model.User{}).Where("id = ?", userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
