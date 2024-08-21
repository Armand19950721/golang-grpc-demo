package model

import (
	"database/sql"
	ArContentPb "service/protos/ArContent"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ArContent struct {
	Id              uuid.UUID                         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserId          uuid.UUID                         `gorm:"type:uuid;index:idx_user_id"`
	Name            string                            `gorm:"type:text"`
	Tag             string                            `gorm:"type:text"`
	IsOn            *bool                             `gorm:"type:boolean"`
	Category        ArContentPb.ArContentCategoryEnum `gorm:"type:int8"`
	Type            ArContentPb.ArContentTypeEnum     `gorm:"type:int8"`
	Template        ArContentPb.ArContentTemplateEnum `gorm:"type:int8"`
	ThumbnailName   string                            `gorm:"type:text"`
	TemplateSetting string                            `gorm:"type:text"`
	ViewerSetting   string                            `gorm:"type:text"`
	ViewerUrlId     sql.NullString                    `gorm:"type:text;index:idx_viewer_url_id"`
	CreateAt        time.Time                         `gorm:"autoCreateTime"`
	UpdateAt        time.Time                         `gorm:"autoUpdateTime"`
	DeleteAt        gorm.DeletedAt
}
