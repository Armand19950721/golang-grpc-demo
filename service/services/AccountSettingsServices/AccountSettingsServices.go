package AccountSettingsServices

import (
	"service/model"
	AccountSettings "service/protos/AccountSettings"
	"service/protos/Common"
	"service/repositories"
	"service/utils"

	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetBasicInfo(ctx context.Context, req *AccountSettings.GetBasicInfoRequest) (*AccountSettings.GetBasicInfoReply, error) {
	reply := &AccountSettings.GetBasicInfoReply{}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	// query
	res, data := repositories.QueryAccountSettings(model.AccountSettings{
		UserId: basicInfo.AdminId,
	})

	notFound := false

	if res.Error != nil {
		if utils.IsErrorNotFound(res.Error) {
			notFound = true
		} else {
			return reply, utils.ReturnUnKnownError(res.Error)
		}
	}

	// convert
	if notFound {
		// return empty model
		reply.Model = &AccountSettings.BasicInfoModel{}
	} else {
		// convert json data
		convertData, err := utils.ParseJsonWithType[AccountSettings.BasicInfoModel](data.SettingsJson)
		if err != nil {
			return reply, utils.ReturnError(Common.ErrorCodes_UNKNOWN_ERROR, err.Error(), "")
		}

		reply.Model = &convertData
	}

	return reply, status.Errorf(codes.OK, "")
}
func UpdateBasicInfo(ctx context.Context, req *AccountSettings.UpdateBasicInfoRequest) (*AccountSettings.UpdateBasicInfoReply, error) {
	reply := &AccountSettings.UpdateBasicInfoReply{}

	// check param
	if !utils.ValidString(req.Model.CompanyName, 1, 100, "nullable") ||
		!utils.ValidString(req.Model.FirstName, 1, 100, "nullable") ||
		!utils.ValidString(req.Model.LastName, 1, 100, "nullable") ||
		!utils.ValidString(req.Model.GuiNumber, 1, 100, "nullable") {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
	}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	// query
	res, _ := repositories.QueryAccountSettings(model.AccountSettings{
		UserId: basicInfo.AdminId,
	})

	notFound := false

	if res.Error != nil {
		if utils.IsErrorNotFound(res.Error) {
			notFound = true
		} else {
			return reply, utils.ReturnUnKnownError(res.Error)
		}
	}

	// convert json
	jsonData := utils.ToJson(req.Model)
	if jsonData == "" {
		return reply, utils.ReturnUnKnownError(res.Error)
	}

	// prepare model
	data := model.AccountSettings{
		UserId:       basicInfo.AdminId,
		SettingsJson: jsonData,
	}

	// insert or update
	if notFound {
		// insert
		res, _ = repositories.CreateAccountSettings(data)

		if res.Error != nil {
			return reply, utils.ReturnUnKnownError(res.Error)
		}
	} else {
		// update
		res, _ = repositories.UpdateAccountSettings(model.AccountSettings{
			UserId: basicInfo.AdminId,
		}, data)

		if res.Error != nil {
			return reply, utils.ReturnUnKnownError(res.Error)
		}
	}

	return reply, status.Errorf(codes.OK, "")
}
func GetLoginInfo(ctx context.Context, req *AccountSettings.GetLoginInfoRequest) (*AccountSettings.GetLoginInfoReply, error) {
	reply := &AccountSettings.GetLoginInfoReply{}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	// convert
	replyModel := AccountSettings.LoginInfoModel{
		AccountName: basicInfo.User.Name,
		EmailAddres: basicInfo.User.Email,
	}

	reply.Model = &replyModel

	return reply, status.Errorf(codes.OK, "")
}
func UpdateLoginInfo(ctx context.Context, req *AccountSettings.UpdateLoginInfoRequest) (*AccountSettings.UpdateLoginInfoReply, error) {
	reply := &AccountSettings.UpdateLoginInfoReply{}

	// check param
	if !utils.ValidString(req.Model.AccountName, 1, 100) {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
	}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	// update
	user := basicInfo.User
	user.Name = req.Model.AccountName

	res, _ := repositories.UpdateUser(model.User{
		Id: basicInfo.AdminId,
	}, user)

	if res.Error != nil {
		return reply, utils.ReturnUnKnownError(res.Error)
	}

	return reply, status.Errorf(codes.OK, "")
}
