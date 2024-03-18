package model

import (
	"github.com/google/uuid"
)

type OauthClient struct {
	ID               uuid.UUID         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name             string            `gorm:"type:varchar(30);not null"`
	ClientID         uuid.UUID         `gorm:"type:uuid;unique"`
	ClientSecret     string            `gorm:"type:varchar(255);unique"`
	RedirectURI      string            `gorm:"type:text"`
	ClientType       string            `gorm:"type:varchar(50)"`
	Scopes           []Scope           `gorm:"many2many:oauth_clients_scopes;"`
	ApplicationTypes []ApplicationType `gorm:"many2many:oauth_clients_application_types;"`
}
