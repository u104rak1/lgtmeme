package model

type OauthClientsApplicationTypes struct {
	ClientID        ClientID `gorm:"column:client_id;primaryKey"`
	ApplicationType string   `gorm:"column:application_type;primaryKey"`
}
