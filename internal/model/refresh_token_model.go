package model

import (
	"github.com/google/uuid"
)

type RefreshToken struct {
	Token    string    `gorm:"primaryKey;size:255;not null"`
	UserID   uuid.UUID `gorm:"type:uuid;not null"`
	ClientID uuid.UUID `gorm:"type:uuid;not null"`
	Scopes   string    `gorm:"type:text;not null"`
}
