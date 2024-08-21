package repositories

import (
	"service/model"
	"service/utils"

	"gorm.io/gorm"
)

func CreateAccountSettings(AccountSettings model.AccountSettings) (*gorm.DB, model.AccountSettings) {
	result := utils.DatabaseManager.Create(&AccountSettings)

	return result, AccountSettings
}

func UpdateAccountSettings(accountSettingsWhere, accountSettings model.AccountSettings) (tx *gorm.DB, data model.AccountSettings) {
	result := utils.DatabaseManager.Where(accountSettingsWhere).Updates(&accountSettings)

	return result, accountSettings
}

func QueryAccountSettings(accountSettingsWhere model.AccountSettings) (tx *gorm.DB, accountSettings model.AccountSettings) {
	result := utils.DatabaseManagerSlave.Where(accountSettingsWhere).First(&accountSettings)

	return result, accountSettings
}
