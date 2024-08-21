package repositories

import (
	"service/model"
	"service/utils"

	"gorm.io/gorm"
)

func CreateMailState(mailState model.MailState) (*gorm.DB, model.MailState) {
	utils.PrintObj(mailState, "CreateMailState")

	result := utils.DatabaseManager.Create(&mailState)

	return result, mailState
}

func QueryMailState(mailStateQuery model.MailState) (tx *gorm.DB, data model.MailState) {
	utils.PrintObj(mailStateQuery, "QueryMailState")

	result := utils.DatabaseManagerSlave.Where(mailStateQuery).First(&data)

	return result, data
}

func QueryMailStateByIdOrUserId(mailStateQuery model.MailState) (tx *gorm.DB, data model.MailState) {
	utils.PrintObj(mailStateQuery, "QueryMailStateByIdOrUserId")

	datas := []model.MailState{}
	result := utils.DatabaseManagerSlave.Debug().Raw(`
		select * 
		from mail_state
		where ( user_id = ? or id = ? )
			and delete_at is null

	`, mailStateQuery.UserId, mailStateQuery.Id).Scan(&datas)

	if len(datas) == 0 {
		result.Error = gorm.ErrRecordNotFound
	} else {
		data = datas[0]
	}

	return result, data
}

func QueryMailStateIncludeDel(mailStateQuery model.MailState) (tx *gorm.DB, data model.MailState) {
	result := utils.DatabaseManagerSlave.Unscoped().Where(mailStateQuery).First(&data)

	return result, data
}

func UpdateMailState(mailStateWhere, mailState model.MailState) (tx *gorm.DB, data model.MailState) {
	utils.PrintObj(mailState, "UpdateMailState")

	result := utils.DatabaseManager.Debug().Where(mailStateWhere).Updates(&mailState).Find(&mailState)

	return result, mailState
}

func GetMailStateList(mailStateWhere model.MailState) (*gorm.DB, []model.MailState) {
	models := []model.MailState{}

	result := utils.DatabaseManagerSlave.Where(mailStateWhere).Order("create_at desc").Find(&models)

	utils.PrintObj(result.RowsAffected, "SelectMailState RowsAffected")

	return result, models
}

func HardDeleteMailState(mailState model.MailState) (tx *gorm.DB, data model.MailState) {
	utils.PrintObj("delete mailState", "HardDeleteMailState")

	result := utils.DatabaseManager.Debug().Unscoped().Where(&mailState).Delete(&mailState)
	return result, mailState
}
