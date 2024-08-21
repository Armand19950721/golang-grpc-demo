package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountSettings struct {
	Id           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserId       uuid.UUID `gorm:"type:uuid;index:idx_user_id"`
	SettingsJson string    `gorm:"type:text"`
	CreateAt     time.Time `gorm:"autoCreateTime"`
	UpdateAt     time.Time `gorm:"autoUpdateTime"`
	DeleteAt     gorm.DeletedAt
}
