package model

type ApplicationType struct {
	Type    string        `gorm:"column:type;type:varchar(20);primaryKey"`
	Clients []OauthClient `gorm:"many2many:oauth_clients_application_types;"`
}
