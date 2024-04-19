package model

import "github.com/google/uuid"

type OauthClientsApplicationTypes struct {
	OAuthClientID   uuid.UUID `gorm:"column:oauth_client_id;primaryKey"`
	ApplicationType string    `gorm:"column:application_type;primaryKey"`
}
