package repositories

import (
	"service/model"
	"service/utils"

	"gorm.io/gorm"
)

func CreateProgram(program model.Program) (*gorm.DB, model.Program) {
	result := utils.DatabaseManager.Create(&program)

	return result, program
}

func GetProgramList(program model.Program) (tx *gorm.DB, models []model.Program, count int64) {
	result := utils.DatabaseManagerSlave.Where(program).Find(&models).Count(&count)

	return result, models, count
}

func QueryProgram(program model.Program) (tx *gorm.DB, data model.Program) {
	result := utils.DatabaseManagerSlave.Where(program).First(&data)

	return result, data
}
