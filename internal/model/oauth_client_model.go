package model

import (
	"github.com/google/uuid"
)

type OauthClient struct {
	ID               uuid.UUID               `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name             string                  `gorm:"column:name;type:varchar(30);not null"`
	ClientID         uuid.UUID               `gorm:"column:client_id;type:uuid;unique"`
	ClientSecret     string                  `gorm:"column:client_secret;type:varchar(255);unique"`
	RedirectURI      string                  `gorm:"column:redirect_uri;type:text"`
	ApplicationURL   string                  `gorm:"column:application_url;type:text;not null"`
	ClientType       string                  `gorm:"column:client_type;type:varchar(50)"`
	Scopes           []MasterScope           `gorm:"many2many:oauth_clients_scopes;foreignKey:ID;joinForeignKey:oauth_client_id;references:Code;joinReferences:scope_code"`
	ApplicationTypes []MasterApplicationType `gorm:"many2many:oauth_clients_application_types;foreignKey:ID;joinForeignKey:client_id;references:Type;joinReferences:application_type"`
}
