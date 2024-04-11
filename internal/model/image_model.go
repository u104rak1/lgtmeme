package model

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4()"`
	URL       string    `gorm:"column:url;not null"`
	Keyword   string    `gorm:"column:keyword;type:varchar(50);not null;default:''"`
	UsedCount int       `gorm:"column:used_count;not null;default:0"`
	Reported  bool      `gorm:"column:reported;not null;default:false"`
	Confirmed bool      `gorm:"column:confirmed;not null;default:false"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()"`
}
