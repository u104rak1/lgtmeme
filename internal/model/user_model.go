package model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name     string    `gorm:"type:varchar(10);not null;unique"`
	Password string    `gorm:"size:255;not null"`
}
