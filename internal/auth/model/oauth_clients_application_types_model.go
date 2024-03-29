package model

import "github.com/google/uuid"

type OauthClientsApplicationTypes struct {
	ClientID        uuid.UUID `gorm:"column:client_id;primaryKey"`
	ApplicationType string    `gorm:"column:application_type;primaryKey"`
}
