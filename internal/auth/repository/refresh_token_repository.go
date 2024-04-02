package repository

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/internal/auth/model"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(c echo.Context, userID uuid.UUID, clientID uuid.UUID, token, scope string) error
	FindByToken(c echo.Context, token string) (model.RefreshToken, error)
	Update(c echo.Context, userID uuid.UUID, clientID uuid.UUID, newToken, scope string) error
}

type refreshTokenRepository struct {
	DB *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{DB: db}
}

func (r *refreshTokenRepository) Create(c echo.Context, userID uuid.UUID, clientID uuid.UUID, token, scope string) error {
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

func (r *refreshTokenRepository) FindByToken(c echo.Context, token string) (model.RefreshToken, error) {
	var refreshToken model.RefreshToken
	if err := r.DB.Where("token = ?", token).First(&refreshToken).Error; err != nil {
		return refreshToken, err
	}

	return refreshToken, nil
}

func (r *refreshTokenRepository) Update(c echo.Context, userID uuid.UUID, clientID uuid.UUID, newToken, scope string) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND client_id = ?", userID, clientID).Delete(&model.RefreshToken{}).Error; err != nil {
			return err
		}

		newRefreshToken := model.RefreshToken{
			Token:    newToken,
			UserID:   userID,
			ClientID: clientID,
			Scopes:   scope,
		}
		if err := tx.Create(&newRefreshToken).Error; err != nil {
			return err
		}

		return nil
	})
}
