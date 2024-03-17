package repository

import (
	"github.com/ucho456job/my_authn_authz/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByName(name string) (*model.User, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) FindByName(name string) (*model.User, error) {
	var user model.User
	if err := r.DB.Model(&model.User{}).Where("name = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
