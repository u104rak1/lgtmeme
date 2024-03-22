package model

type OauthClientsScopes struct {
	ClientID  ClientID `gorm:"column:client_id;primaryKey"`
	ScopeCode string   `gorm:"column:scope_code;primaryKey"`
}
