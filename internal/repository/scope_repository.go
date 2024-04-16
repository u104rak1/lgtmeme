package repository

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/internal/model"
	"gorm.io/gorm"
)

type ScopeRepository interface {
	FindByScopesStr(c echo.Context, scopesStr string) ([]model.MasterScope, error)
}

type scopeRepository struct {
	DB *gorm.DB
}

func NewScopeRepository(db *gorm.DB) ScopeRepository {
	return &scopeRepository{DB: db}
}

func (r *scopeRepository) FindByScopesStr(c echo.Context, scopesStr string) ([]model.MasterScope, error) {
	var scopes []model.MasterScope

	if err := r.DB.Where("code IN ?", strings.Split(scopesStr, " ")).Find(&scopes).Error; err != nil {
		return nil, err
	}

	return scopes, nil
}
