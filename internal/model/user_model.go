package model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name     string    `gorm:"type:varchar(20);not null;unique"`
	Password string    `gorm:"size:255;not null"`
}
