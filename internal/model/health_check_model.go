package model

type HealthCheck struct {
	Key   string `gorm:"column:key;type:varchar(20);not null;primary_key"`
	Value string `gorm:"column:value;type:varchar(20);not null"`
}
