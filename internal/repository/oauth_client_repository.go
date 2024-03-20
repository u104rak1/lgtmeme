package repository

import (
	"errors"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ucho456job/my_authn_authz/internal/dto"
	"github.com/ucho456job/my_authn_authz/internal/model"
	"gorm.io/gorm"
)

type OauthClientRepository interface {
	ExistsForAuthz(c echo.Context, q dto.AuthoraizationQuery) (bool, error)
}

type oauthClientRepository struct {
	DB *gorm.DB
}

func NewOauthClientRepository(db *gorm.DB) OauthClientRepository {
	return &oauthClientRepository{DB: db}
}

func (r *oauthClientRepository) ExistsForAuthz(c echo.Context, q dto.AuthoraizationQuery) (bool, error) {
	var oauthClient model.OauthClient
	if err := r.DB.Model(&model.OauthClient{}).Preload("Scopes").Where("client_id = ? AND redirect_uri = ?", q.ClientID, q.RedirectURI).First(&oauthClient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	scopes := strings.Split(q.Scope, " ")
	for _, s := range scopes {
		found := false
		for _, cs := range oauthClient.Scopes {
			if s == cs.Code {
				found = true
				break
			}
		}
		if !found {
			return false, nil
		}
	}

	return true, nil
}
