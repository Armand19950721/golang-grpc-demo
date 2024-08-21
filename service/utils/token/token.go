package token

import (
	"service/model"
	"service/protos/Common"
	"service/repositories"
	"service/utils"

	"github.com/google/uuid"
)

func ValidToken(token string) (Common.ErrorCodes, *model.User) {
	utils.PrintObj(token, "validToken")
	userModelEmpty := &model.User{}

	// valid param
	if token == "" || // find token in metadata
		!utils.ValidString(token, 1, -1) { // check str length

		utils.PrintObj("token format err")
		return Common.ErrorCodes_INVAILD_TOKEN, userModelEmpty
	}

	// check redis and get user id
	userId := utils.GetRedis("token:" + token)

	if userId == "" {
		utils.PrintObj("token not found")
		return Common.ErrorCodes_INVAILD_TOKEN, userModelEmpty
	}

	// query user
	resultQuery, userData := repositories.QueryUser(model.User{Id: uuid.MustParse(userId)})

	if resultQuery.Error != nil {
		utils.PrintObj(resultQuery.Error.Error())
		return Common.ErrorCodes_INVAILD_TOKEN, userModelEmpty
	}

	utils.PrintObj("token pass")
	return Common.ErrorCodes_SUCCESS, &userData
}
