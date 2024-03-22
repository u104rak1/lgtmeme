package model

import "github.com/google/uuid"

type UserID uuid.UUID

type User struct {
	ID       UserID `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name     string `gorm:"column:name;type:varchar(20);not null;unique"`
	Password string `gorm:"column:password;size:255;not null"`
	Role     string `gorm:"column:role;type:varchar(20);"`
}

func ParseUserID(id string) (UserID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return UserID{}, err
	}
	return UserID(parsedID), nil
}
