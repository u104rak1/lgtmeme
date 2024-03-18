package model

type Scope struct {
	Code        string        `gorm:"type:varchar(20);primaryKey"`
	Description string        `gorm:"type:text"`
	Clients     []OauthClient `gorm:"many2many:oauth_clients_scopes;"`
}
