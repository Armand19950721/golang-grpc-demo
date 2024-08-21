package model

import (
	"service/protos/Common"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MailState struct {
	Id        uuid.UUID        `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	AdminId   uuid.UUID        `gorm:"type:uuid"`
	UserId    uuid.UUID        `gorm:"type:uuid;index:idx_user_id;"`
	MailType  Common.MailType  `gorm:"type:int8"`
	MailState Common.MailState `gorm:"type:int8"`
	CreateAt  time.Time        `gorm:"autoCreateTime"`
	UpdateAt  time.Time
	DeleteAt  gorm.DeletedAt
}
