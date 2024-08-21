package model

import (
	"time"

	"github.com/google/uuid"
)

// 每日總結一次 Group by user
type StatisticArContentUserSum struct {
	Id                     uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserId                 uuid.UUID `gorm:"type:uuid;index:idx_user_id"`
	ViewerOpen             int32
	ViewerInited           int32
	ViewerTakePicture      int32
	ClickViewerRightButton int32
	ClickViewerLeftButton  int32
	CreateAt               time.Time `gorm:"autoCreateTime"`
	UpdateAt               time.Time
}
