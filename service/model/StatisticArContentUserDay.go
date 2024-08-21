package model

import (
	"time"

	"github.com/google/uuid"
)

// 每日從Redis紀錄歸檔進來
type StatisticArContentUserDay struct {
	Id                     uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserId                 uuid.UUID `gorm:"type:uuid;index:idx_admin_id"`
	ArContentId            uuid.UUID `gorm:"type:uuid;index:idx_ar_content_id"`
	ViewerOpen             int32
	ViewerInited           int32
	ViewerTakePicture      int32
	ClickViewerRightButton int32
	ClickViewerLeftButton  int32
	StatisticDate          time.Time
	CreateAt               time.Time `gorm:"autoCreateTime"`
}
