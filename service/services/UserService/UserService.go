package UserService

import (
	"service/model"
	"service/protos/Common"
	"service/protos/User"
	"service/repositories"
	// "service/test/testUtils"
	"context"
	"service/utils"
	EmailUtils "service/utils/EmailUtils"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AddUserChild(ctx context.Context, req *User.AddUserChildRequest) (*Common.CommonReply, error) {
	reply := &Common.CommonReply{}

	// check param
	if !utils.ValidEmail(req.Model.Email) ||
		// !utils.ValidString(req.Model.Password, 1, 100) ||
		!utils.ValidString(req.Model.Name, 1, 100) ||
		!utils.ValidId(req.Model.PermissionGroupId) {

		return reply, status.Errorf(codes.InvalidArgument, "param invalid")
	}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	// check user email limit
	validEmail, err := utils.CheckUserEmailLimit(basicInfo.User.Id)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	if !validEmail {
		return reply, utils.ReturnError(Common.ErrorCodes_REACH_DAILY_MAIL_LIMIT, "", "daily limit reached")
	}

	// check count
	userWhere := model.User{
		ParentId: utils.SetNullableUUID(basicInfo.AdminId.String()),
	}

	res, count := repositories.GetUserCount(userWhere)
	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	utils.PrintObj([]int{int(count), int(basicInfo.Program.Seats)}, "count vs program seats")

	if int(count) >= int(basicInfo.Program.Seats) {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:      Common.ErrorCodes_SEATS_LIMIT_REACHED,
			ReturnMsg: "out of program seats limit",
		}))
	}

	// check account(email) repeat
	res, _ = repositories.QueryUser(model.User{Email: req.Model.Email})

	if res.Error == nil {
		return reply, status.Error(codes.AlreadyExists, "")
	}

	// check group id
	res, _ = repositories.QueryPermissionGroup(model.PermissionGroup{
		Id: utils.ParseUUID(req.Model.PermissionGroupId),
		// UserId: basicInfo.AdminId,
	})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
			ReturnMsg:   "permission group not found",
		}))
	}

	// create user
	valid := false

	res, user := repositories.CreateUser(model.User{
		Email:             req.Model.Email,
		Name:              req.Model.Name,
		PermissionGroupId: utils.SetNullableUUID(req.Model.PermissionGroupId),
		ParentId:          utils.SetNullableUUID(utils.GetParentId(basicInfo).String()),
		EmailValid:        &valid,
	})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// create mail state
	res, mailState := repositories.CreateMailState(model.MailState{
		AdminId:   basicInfo.AdminId,
		UserId:    user.Id,
		MailType:  Common.MailType_MAIL_TYPE_CHILD,
		MailState: Common.MailState_MAIL_STATE_INVITED,
		UpdateAt:  time.Now(),
	})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// send email for confirm
	err = sendChildInvite(user.Email, mailState.Id)

	if err != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	return reply, status.Errorf(codes.OK, "")
}

func UpdateUserChild(ctx context.Context, req *User.UpdateUserChildRequest) (*Common.CommonReply, error) {
	reply := &Common.CommonReply{}

	// check param
	if !utils.ValidId(req.Model.UserId) ||
		!utils.ValidString(req.Model.Name, 1, 100) ||
		!utils.ValidId(req.Model.PermissionGroupId) {

		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
	}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	// check group id
	res, _ := repositories.QueryPermissionGroup(model.PermissionGroup{
		Id: utils.ParseUUID(req.Model.PermissionGroupId),
		// UserId: basicInfo.AdminId,
	})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
			ReturnMsg:   "permission group not found",
		}))
	}

	// update
	res, _ = repositories.UpdateUser(
		model.User{Id: utils.ParseUUID(req.Model.UserId)},
		model.User{
			Name:              req.Model.Name,
			PermissionGroupId: utils.SetNullableUUID(req.Model.PermissionGroupId),
		},
	)

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	return reply, status.Errorf(codes.OK, "")
}

func GetUserChildList(ctx context.Context, req *User.GetUserChildListRequest) (*User.GetUserChildListReply, error) {
	reply := &User.GetUserChildListReply{}
	// check param
	req.PageInfo = utils.ValidPageInfo(req.PageInfo)

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	// get child list
	res, users, totalCount := repositories.GetUserList(
		model.User{ParentId: utils.SetNullableUUID(utils.GetParentId(basicInfo).String())},
		req.PageInfo,
	)

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// get mail state
	res, mailStates := repositories.GetMailStateList(model.MailState{AdminId: basicInfo.AdminId})
	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// convert
	models := []*User.UserChildModel{}

	for _, user := range users {
		state := Common.MailState_MAIL_STATE_NONE
		mailExpireDate := ""

		// get mail state
		for _, mailStateItem := range mailStates {
			if user.Id == mailStateItem.UserId {
				state, mailExpireDate = getMailStateByModel(mailStateItem)
			}
		}

		models = append(models, &User.UserChildModel{
			UserId:            user.Id.String(),
			Email:             user.Email,
			Name:              user.Name,
			PermissionGroupId: user.PermissionGroupId.UUID.String(),
			MailState:         state,
			MailExpireDate:    mailExpireDate,
		})
	}

	reply.Models = models
	reply.PageInfo = &Common.PageInfoReply{TotalCount: totalCount}

	return reply, status.Errorf(codes.OK, "")
}

func DeleteUserChild(ctx context.Context, req *User.DeleteUserChildRequest) (*Common.CommonReply, error) {
	reply := &Common.CommonReply{}

	// check param
	if !utils.ValidId(req.Model.UserId) {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
	}

	// delete
	res, _ := repositories.HardDeleteUser(model.User{Id: utils.ParseUUID(req.Model.UserId)})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	return reply, status.Errorf(codes.OK, "")
}

func ForgotPassword(ctx context.Context, req *User.ForgotPasswordRequest) (*Common.CommonReply, error) {
	reply := &Common.CommonReply{}

	// check param
	if !utils.ValidEmail(req.Email) {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
	}

	// vaild user
	res, user := repositories.QueryUser(model.User{Email: req.Email})

	if res.Error != nil {
		if utils.IsErrorNotFound(res.Error) {
			return reply, utils.ReturnError(Common.ErrorCodes_USER_EMAIL_INVAILD, "", "account not exist")
		} else {
			return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
				Code:        Common.ErrorCodes_UNKNOWN_ERROR,
				InternalMsg: res.Error.Error(),
			}))
		}
	}

	// check user email limit
	valid, err := utils.CheckUserEmailLimit(user.Id)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	if !valid {
		return reply, utils.ReturnError(Common.ErrorCodes_REACH_DAILY_MAIL_LIMIT, "", "daily limit reached")
	}

	// create new pwd
	newPwd := utils.GetRandomShortString()

	// for dev test start (ignore)
	// testInfo := testUtils.GetBasicInfo("foo")
	// testInfo.Remark = newPwd
	// utils.PrintObj(testInfo, "test forgot pwd. save remark")
	// testUtils.SetBasicInfo(testInfo, "foo")
	// end

	// pwd encode
	newPwdEncode, err := utils.GenHashPassword(newPwd)
	if err != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:      Common.ErrorCodes_UNKNOWN_ERROR,
			ReturnMsg: "encode pwd err",
		}))
	}

	// update pwd
	res, _ = repositories.UpdateUser(
		model.User{Id: user.Id},
		model.User{Password: newPwdEncode},
	)

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// send new password email
	EmailUtils.SendEmail(req.Email, EmailUtils.GetTemplateForgotPasswod(newPwd))

	return reply, status.Errorf(codes.OK, "")
}

func ChangePassword(ctx context.Context, req *User.ChangePasswordRequest) (*Common.CommonReply, error) {
	reply := &Common.CommonReply{}

	// check param
	if !utils.ValidPassword(req.OldPassword) ||
		!utils.ValidPassword(req.NewPassword) {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
	}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}
	userId := utils.ParseUUID(basicInfo.User.Id.String())

	// get user
	res, user := repositories.QueryUser(model.User{Id: userId})
	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// check pwd
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:      Common.ErrorCodes_UNKNOWN_ERROR,
			ReturnMsg: "old pwd err",
		}))
	}

	// pwd encode
	newPwd, err := utils.GenHashPassword(req.NewPassword)
	if err != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:      Common.ErrorCodes_UNKNOWN_ERROR,
			ReturnMsg: "encode pwd err",
		}))
	}

	// update pwd
	res, _ = repositories.UpdateUser(
		model.User{Id: userId},
		model.User{Password: newPwd},
	)

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	utils.PrintObj(userId, "update pwd")
	utils.PrintObj(newPwd)

	return reply, status.Errorf(codes.OK, "")
}

func UpdateUserProgram(ctx context.Context, req *User.UpdateUserProgramRequest) (*Common.CommonReply, error) {
	reply := &Common.CommonReply{}

	// check param
	if !utils.ValidId(req.ProgramId) {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
	}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	// update user program
	res, _ := repositories.UpdateUser(
		model.User{Id: basicInfo.User.Id},
		model.User{
			ProgramId: utils.SetNullableUUID(req.ProgramId),
		},
	)
	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	return reply, status.Errorf(codes.OK, "")
}

func GetConfirmUserChildInfo(ctx context.Context, req *User.GetConfirmUserChildInfoRequest) (*User.GetConfirmUserChildInfoReply, error) {
	reply := &User.GetConfirmUserChildInfoReply{}

	// check param
	if !utils.ValidId(req.UrlKey) {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
	}

	res, mailState := repositories.QueryMailState(model.MailState{
		Id: utils.ParseUUID(req.UrlKey),
	})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// check mail state is invited
	whereParam := model.MailState{
		Id: utils.ParseUUID(req.UrlKey),
	}

	res, data := repositories.QueryMailState(whereParam)
	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	if data.MailState != Common.MailState_MAIL_STATE_INVITED {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:      Common.ErrorCodes_UNKNOWN_ERROR,
			ReturnMsg: "mail state err",
		}))
	}

	// check expire
	expire, err := isMailStateExpireByUrlKey(req.UrlKey)

	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	if expire {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:      Common.ErrorCodes_MAIL_STATE_EXPIRE_ERROR,
			ReturnMsg: "mail state expire",
		}))
	}

	// get user email
	res, user := repositories.QueryUser(model.User{
		Id: mailState.UserId,
	})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	reply.Account = user.Email

	return reply, status.Errorf(codes.OK, "")
}

func ConfirmUserChild(ctx context.Context, req *User.ConfirmUserChildRequest) (*Common.CommonReply, error) {
	reply := &Common.CommonReply{}

	// check param
	if !utils.ValidId(req.UrlKey) ||
		!utils.ValidPassword(req.Password) {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
	}

	// check mail state is invited
	whereParam := model.MailState{
		Id: utils.ParseUUID(req.UrlKey),
	}

	res, data := repositories.QueryMailState(whereParam)
	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	if data.MailState != Common.MailState_MAIL_STATE_INVITED {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:      Common.ErrorCodes_UNKNOWN_ERROR,
			ReturnMsg: "mail state err",
		}))
	}

	// check expire
	expire, err := isMailStateExpireByUrlKey(req.UrlKey)

	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	if expire {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:      Common.ErrorCodes_MAIL_STATE_EXPIRE_ERROR,
			ReturnMsg: "mail state expire",
		}))
	}

	// update mail state

	updateParam := model.MailState{
		MailState: Common.MailState_MAIL_STATE_SUCCESS,
		UpdateAt:  time.Now(),
	}

	res, data = repositories.UpdateMailState(whereParam, updateParam)

	utils.PrintObj(data, "data")

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// user update valid
	valid := true
	// encode pwd
	encodePwd, err := utils.GenHashPassword(req.Password)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	res, user := repositories.UpdateUser(model.User{
		Id: data.UserId,
	}, model.User{
		Password:   encodePwd,
		EmailValid: &valid,
	})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	utils.PrintObj(user, "user")

	return reply, status.Errorf(codes.OK, "")
}

func ResendChildInvite(ctx context.Context, req *User.ResendChildInviteRequest) (*User.ResendChildInviteReply, error) {
	reply := &User.ResendChildInviteReply{}

	// check param
	if !utils.ValidId(req.UserId) {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
	}

	// query user
	res, user := repositories.QueryUser(model.User{
		Id: utils.ParseUUID(req.UserId),
	})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// check user email limit
	valid, err := utils.CheckUserEmailLimit(user.Id)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	if !valid {
		return reply, utils.ReturnError(Common.ErrorCodes_REACH_DAILY_MAIL_LIMIT, "", "daily limit reached")
	}

	// query mail state
	res, mailStateQuery := repositories.QueryMailState(model.MailState{
		UserId: utils.ParseUUID(req.UserId),
	})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	if mailStateQuery.MailState != Common.MailState_MAIL_STATE_INVITED {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:      Common.ErrorCodes_UNKNOWN_ERROR,
			ReturnMsg: "mail state err",
		}))
	}

	// check 12hr resend
	expire, _ := utils.CheckExpire(12*time.Hour, mailStateQuery.UpdateAt) // if expire means can resend
	if !expire {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:      Common.ErrorCodes_MAIL_HEAVY_FREQUENTLY,
			ReturnMsg: "still in 12 hours range",
		}))
	}

	// update mail state
	whereParam := model.MailState{
		Id: mailStateQuery.Id,
	}

	updateParam := model.MailState{
		MailState: Common.MailState_MAIL_STATE_INVITED,
		UpdateAt:  time.Now(),
	}

	res, _ = repositories.UpdateMailState(whereParam, updateParam)

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// send email for confirm
	err = sendChildInvite(user.Email, mailStateQuery.Id)

	if err != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	reply.UrlKey = mailStateQuery.Id.String()

	return reply, status.Errorf(codes.OK, "")
}

func CancelChildInvite(ctx context.Context, req *User.CancelChildInviteRequest) (*Common.CommonReply, error) {
	reply := &Common.CommonReply{}

	// check param
	if !utils.ValidId(req.UserId) {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
	}

	// hard delete child
	res, _ := repositories.HardDeleteUser(model.User{Id: utils.ParseUUID(req.UserId)})
	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// hard delete mailState
	res, _ = repositories.HardDeleteMailState(model.MailState{UserId: utils.ParseUUID(req.UserId)})
	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	return reply, status.Errorf(codes.OK, "")
}
