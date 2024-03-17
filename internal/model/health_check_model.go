package model

import "github.com/google/uuid"

type HealthCheck struct {
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Key   string    `gorm:"type:varchar(20);not null"`
	Value string    `gorm:"type:varchar(20);not null"`
}
