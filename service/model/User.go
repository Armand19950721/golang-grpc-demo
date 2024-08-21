package model

import (
	"service/protos/Common"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id                uuid.UUID           `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email             string              `gorm:"type:text;unique;"`
	Password          string              `gorm:"type:text;"`
	Name              string              `gorm:"type:text"`
	BusinessType      Common.BusinessType `gorm:"type:int8"`
	ProgramId         uuid.NullUUID       `gorm:"type:uuid"`
	ParentId          uuid.NullUUID       `gorm:"type:uuid;index:idx_parent_id;"`
	PermissionGroupId uuid.NullUUID       `gorm:"type:uuid"`
	EmailValid        *bool               `gorm:"type:boolean"`
	CreateAt          time.Time           `gorm:"autoCreateTime"`
	UpdateAt          time.Time           `gorm:"autoUpdateTime"`
	DeleteAt          gorm.DeletedAt
}
