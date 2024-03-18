package model

import "github.com/google/uuid"

type OauthClientsApplicationTypes struct {
	ClientID        uuid.UUID `gorm:"primaryKey"`
	ApplicationType string    `gorm:"primaryKey"`
}
