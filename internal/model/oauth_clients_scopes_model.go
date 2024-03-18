package model

import "github.com/google/uuid"

type OauthClientsScopes struct {
	ClientID  uuid.UUID `gorm:"primaryKey"`
	ScopeCode string    `gorm:"primaryKey"`
}
