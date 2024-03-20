package model

import "github.com/google/uuid"

type OauthClientsScopes struct {
	ClientID  uuid.UUID `gorm:"column:client_id;primaryKey"`
	ScopeCode string    `gorm:"column:scope_code;primaryKey"`
}
