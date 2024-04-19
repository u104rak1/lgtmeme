package repository

// mockgen -source=internal/repository/oauth_client_repository.go -destination=test/mock/repository/mock_oauth_client_repository.go -package=repository_mock

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ucho456job/lgtmeme/internal/dto"
	"github.com/ucho456job/lgtmeme/internal/model"
	"gorm.io/gorm"
)

type OauthClientRepository interface {
	IsValidOAuthClient(c echo.Context, q dto.AuthzQuery) (bool, error)
	FirstByClientID(c echo.Context, clientID uuid.UUID, columns []string) (*model.OauthClient, error)
}

type oauthClientRepository struct {
	DB *gorm.DB
}

func NewOauthClientRepository(db *gorm.DB) OauthClientRepository {
	return &oauthClientRepository{DB: db}
}

func (r *oauthClientRepository) IsValidOAuthClient(c echo.Context, q dto.AuthzQuery) (bool, error) {
	var dbScopes []model.OauthClientsScopes
	if err := r.DB.Raw(`
			SELECT osc.scope_code
			FROM oauth_clients AS oc
			INNER JOIN oauth_clients_scopes AS osc ON oc.id = osc.oauth_client_id
			WHERE oc.client_id = ? AND oc.redirect_uri = ?
	`, q.ClientID, q.RedirectURI).Scan(&dbScopes).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	scopes := strings.Split(q.Scope, " ")
	for _, s := range scopes {
		found := false
		for _, dbS := range dbScopes {
			if s == dbS.ScopeCode {
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

// If columns is empty, return all columns. However, related records are excluded.
func (r *oauthClientRepository) FirstByClientID(c echo.Context, clientID uuid.UUID, columns []string) (*model.OauthClient, error) {
	var oauthClient model.OauthClient
	q := r.DB.Model(&model.OauthClient{}).Where("client_id = ?", clientID)

	if len(columns) > 0 {
		q = q.Select(columns)
	}

	if err := q.First(&oauthClient).Error; err != nil {
		return nil, err
	}
	return &oauthClient, nil
}
