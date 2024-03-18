package model

type HealthCheck struct {
	Key   string `gorm:"type:varchar(20);not null;primary_key"`
	Value string `gorm:"type:varchar(20);not null"`
}
