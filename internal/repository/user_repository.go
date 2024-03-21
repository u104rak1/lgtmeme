package repository

import (
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/my_authn_authz/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByName(c echo.Context, name string) (*model.User, error)
	ExistsByID(c echo.Context, userID string) (bool, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) FindByName(c echo.Context, name string) (*model.User, error) {
	var user model.User
	if err := r.DB.Model(&model.User{}).Where("name = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) ExistsByID(c echo.Context, userID string) (bool, error) {
	var count int64
	if err := r.DB.Model(&model.User{}).Where("id = ?", userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
