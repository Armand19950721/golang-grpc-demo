package repositories

import (
	"service/model"
	"service/utils"

	"gorm.io/gorm"
)

func CreatePermissionGroup(permissionGroup model.PermissionGroup) (*gorm.DB, model.PermissionGroup) {
	result := utils.DatabaseManager.Create(&permissionGroup)

	return result, permissionGroup
}

func GetPermissionGroupList(permissionGroupWhere model.PermissionGroup) (tx *gorm.DB, models []model.PermissionGroup, count int64) {
	result := utils.DatabaseManagerSlave.Where(permissionGroupWhere).Find(&models).Count(&count)

	return result, models, count
}

func QueryPermissionGroup(permissionGroup model.PermissionGroup) (tx *gorm.DB, data model.PermissionGroup) {
	result := utils.DatabaseManagerSlave.Where(permissionGroup).First(&data)

	return result, data
}

func UpdatePermissionGroup(permissionGroupWhere, permissionGroup model.PermissionGroup) (tx *gorm.DB, data model.PermissionGroup) {
	result := utils.DatabaseManager.Where(permissionGroupWhere).Updates(&permissionGroup)

	return result, permissionGroup
}
