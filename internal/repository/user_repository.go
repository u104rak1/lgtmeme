package repository

import (
	"github.com/ucho456job/my_authn_authz/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByName(name string) (*model.User, error)
}

type GormUserRepository struct {
	DB *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{DB: db}
}

func (r *GormUserRepository) FindByName(name string) (*model.User, error) {
	var user model.User
	if err := r.DB.Where("name = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
