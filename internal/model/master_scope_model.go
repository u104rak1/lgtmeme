package model

type MasterScope struct {
	Code        string        `gorm:"column:code;type:varchar(20);primaryKey"`
	Description string        `gorm:"column:description;type:text"`
	Clients     []OauthClient `gorm:"many2many:oauth_clients_scopes;"`
}
