package UserService

import (
	"service/model"
	"service/protos/Common"
	"service/repositories"
	"service/utils"
	"service/utils/EmailUtils"
	"time"

	"github.com/google/uuid"
)

func sendChildInvite(email string, mailStateId uuid.UUID) error {
	domain := utils.GetDomainAPI()
	url := domain + utils.GetEnv("DOMAIN_INVITE_ROUTE") + mailStateId.String()
	err := EmailUtils.SendEmail(email, EmailUtils.GetTemplateAddUserChild(url))

	if err != nil {
		return err
	}

	return nil
}

func getMailStateByModel(item model.MailState) (state Common.MailState, expireDateStr string) {
	expire, endDate := utils.CheckExpire(48*time.Hour, item.UpdateAt)

	if !expire {
		utils.PrintObj("INVITED", "getMailStateByModel")

		return Common.MailState_MAIL_STATE_INVITED, utils.ParseDateToString(endDate)
	}

	return Common.MailState_MAIL_STATE_EXPIRE, ""
}

func isMailStateExpireByUrlKey(urlKey string) (bool, error) {
	whereParam := model.MailState{
		Id: utils.ParseUUID(urlKey),
	}

	res, data := repositories.QueryMailState(whereParam)

	if res.Error != nil {
		return true, res.Error
	}

	stateEnum, _ := getMailStateByModel(data)
	if stateEnum == Common.MailState_MAIL_STATE_EXPIRE {
		utils.PrintObj(true, "isMailStateExpire")
		return true, nil
	}

	utils.PrintObj(false, "isMailStateExpire")
	return false, nil
}
