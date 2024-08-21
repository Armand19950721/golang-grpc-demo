package repositories

import (
	"service/model"
	"service/utils"

	"gorm.io/gorm"
)

func CreateRedeemToken(redeemToken model.RedeemToken) (*gorm.DB, model.RedeemToken) {
	result := utils.DatabaseManager.Create(&redeemToken)

	return result, redeemToken
}

func UpdateRedeemToken(redeemTokenWhere, redeemToken model.RedeemToken) (tx *gorm.DB, data model.RedeemToken) {
	result := utils.DatabaseManager.Where(redeemTokenWhere).Updates(&redeemToken)

	return result, redeemToken
}

func QueryRedeemToken(redeemTokenWhere model.RedeemToken) (tx *gorm.DB, redeemToken model.RedeemToken) {

	datas := []model.RedeemToken{}

	result := utils.DatabaseManagerSlave.Debug().Raw(`
		select * 
		from redeem_token
		where ( user_id = ? 
			or mail_state_id = ? )
			and delete_at is null

	`, redeemTokenWhere.UserId, redeemTokenWhere.MailStateId).Scan(&datas)

	if len(datas) == 0 {
		result.Error = gorm.ErrRecordNotFound
	} else {
		redeemToken = datas[0]
	}

	return result, redeemToken
}

func DeleteRedeemToken(redeemTokenWhere, redeemToken model.RedeemToken) (*gorm.DB, model.RedeemToken) {
	result := utils.DatabaseManager.Where(redeemTokenWhere).Delete(&redeemToken)

	return result, redeemToken
}
