package model

type RefreshToken struct {
	Token    string   `gorm:"primaryKey;size:255;not null"`
	UserID   UserID   `gorm:"type:uuid;not null"`
	ClientID ClientID `gorm:"type:uuid;not null"`
	Scopes   string   `gorm:"type:text;not null"`
}
