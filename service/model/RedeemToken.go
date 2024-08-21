package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RedeemToken struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserId      string    `gorm:"type:text;index:idx_user_id;"`
	MailStateId string    `gorm:"type:text;index:idx_mail_state_id;"`
	CreateAt    time.Time `gorm:"autoCreateTime"`
	UpdateAt    time.Time `gorm:"autoUpdateTime"`
	DeleteAt    gorm.DeletedAt
}
