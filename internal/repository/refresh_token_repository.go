package repository

import (
	"github.com/google/uuid"
	"github.com/ucho456job/my_authn_authz/internal/model"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	CreateRefreshToken(userID uuid.UUID, clientID uuid.UUID, token, scope string) error
}

type refreshTokenRepository struct {
	DB *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{DB: db}
}

func (r *refreshTokenRepository) CreateRefreshToken(userID uuid.UUID, clientID uuid.UUID, token, scope string) error {
	refreshToken := model.RefreshToken{
		Token:    token,
		UserID:   userID,
		ClientID: clientID,
		Scopes:   scope,
	}

	if err := r.DB.Create(&refreshToken).Error; err != nil {
		return err
	}

	return nil
}
