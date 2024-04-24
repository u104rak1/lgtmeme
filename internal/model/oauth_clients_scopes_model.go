package model

import "github.com/google/uuid"

type OauthClientsScopes struct {
	OAuthClientID uuid.UUID `gorm:"column:oauth_client_id;primaryKey"`
	ScopeCode     string    `gorm:"column:scope_code;primaryKey"`
}
