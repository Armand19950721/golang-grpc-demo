package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"service/protos/Common"
)

type Program struct {
	Id                  uuid.UUID           `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name                string              `gorm:"type:text;unique;"`
	Click               int64               `gorm:"type:int"`
	State               Common.ProgramState `gorm:"type:int8"`
	Seats               int32               `gorm:"type:int"`
	Categories          string              `gorm:"type:text"`
	Types               string              `gorm:"type:text"`
	Templates           string              `gorm:"type:text"`
	EffectTools         string              `gorm:"type:text"`
	ArInteractModules   string              `gorm:"type:text"`
	ArEditWindowModules string              `gorm:"type:text"`
	CreateAt            time.Time           `gorm:"autoCreateTime"`
	UpdateAt            time.Time           `gorm:"autoUpdateTime"`
	DeleteAt            gorm.DeletedAt
}
