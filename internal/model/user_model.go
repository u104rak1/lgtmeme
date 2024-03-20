package model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name     string    `gorm:"column:name;type:varchar(20);not null;unique"`
	Password string    `gorm:"column:password;size:255;not null"`
	Role     string    `gorm:"column:role;type:varchar(20);"`
}
