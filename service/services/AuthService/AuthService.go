package AuthService

import (
	"service/model"
	"service/protos/Auth"
	"service/protos/Common"
	"time"

	"context"
	"service/repositories"
	"service/utils"

	EmailUtils "service/utils/EmailUtils"

	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoginService(ctx context.Context, req *Auth.LoginRequest) (*Auth.AuthReply, error) {
	reply := &Auth.AuthReply{}

	// check param
	if !utils.ValidEmail(req.Email) ||
		!utils.ValidPassword(req.Password) {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
	}

	// valid
	user := model.User{
		Email: req.Email,
	}
	result, user := repositories.QueryUser(user)

	if result.Error != nil {
		if utils.IsErrorNotFound(result.Error) {
			return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{
				Code: Common.ErrorCodes_LOGIN_CHECK_ID_FAIL,
			}))
		} else {
			return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
				Code:        Common.ErrorCodes_UNKNOWN_ERROR,
				InternalMsg: result.Error.Error(),
			}))
		}
	}

	utils.PrintObj(user, "user")

	if !utils.CheckPassword(req.Password, user.Password) {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_LOGIN_CHECK_ID_FAIL}))
	}

	// check email valid
	reply.UserId = user.Id.String()

	if !*user.EmailValid {
		reply.Code = Common.ErrorCodes_USER_EMAIL_INVAILD
		return reply, status.Errorf(codes.OK, "")
	}

	// create token
	token, expireTime, err := UpdateOrCreateTokenInRedis(user.Id.String())

	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	reply.Token = "Bearer " + token
	reply.ExpireDate = utils.ParseDateToString(expireTime)
	reply.Code = Common.ErrorCodes_SUCCESS

	return reply, status.Errorf(codes.OK, "")
}

func UpdateOrCreateTokenInRedis(userId string) (token string, expireTime time.Time, err error) {
	expireHours := 24
	expireTime = time.Now().Add(time.Duration(expireHours) * time.Hour)

	// check token
	lastToken := utils.GetRedis("userId:" + userId)
	if lastToken != "" {
		// if exsit will be delete
		utils.DeleteRedis("token:" + lastToken)
		utils.DeleteRedis("userId:" + userId)
	}

	// save new token
	token = utils.GenerateBearerToken()
	utils.SetRedis("token:"+token, userId, int16(expireHours))
	utils.SetRedis("userId:"+userId, token, int16(expireHours))

	return token, expireTime, nil
}

func LogoutService(ctx context.Context, req *Auth.LogoutRequest) (*Common.CommonReply, error) {
	reply := &Common.CommonReply{}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	// delete token
	token := utils.GetRedis("userId:" + basicInfo.User.Id.String())

	utils.DeleteRedis("token:" + token)
	utils.DeleteRedis("userId:" + basicInfo.User.Id.String())

	return reply, status.Errorf(codes.OK, "")
}

func RegisterService(ctx context.Context, req *Auth.RegisterRequest) (*Auth.RegisterReply, error) {
	reply := &Auth.RegisterReply{}

	// check param
	if !utils.ValidEmail(req.Email) ||
		!utils.ValidPassword(req.Password) ||
		!utils.ValidString(req.Name, 1, 100, "nullable") ||
		!utils.ValidString(req.PhoneNumber, 1, 15) ||
		!utils.ValidId(req.ProgramId, "nullable") {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
	}

	// check repeat
	user := model.User{
		Email: req.Email,
	}
	res, count := repositories.GetUserCount(user)

	if res.Error == nil && count > 0 {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code: Common.ErrorCodes_REPEATED_ERROR,
		}))
	}

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// check program
	res, _ = repositories.QueryProgram(model.Program{Id: utils.ParseUUID(req.ProgramId)})
	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
			ReturnMsg:   "註冊時檢查 program 失敗",
		}))
	}

	// encode pwd
	encodePwd, err := utils.GenHashPassword(req.Password)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	// create user
	valid := false

	user = model.User{
		Email:        req.Email,
		Password:     encodePwd,
		Name:         req.Name,
		BusinessType: Common.BusinessType_COMPANY,
		ProgramId:    utils.SetNullableUUID(req.ProgramId),
		EmailValid:   &valid,
	}

	result, user := repositories.CreateUser(user)
	if result.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: result.Error.Error(),
		}))
	}

	// after register then login
	if !utils.CheckPassword(req.Password, user.Password) {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: "pwd re-check err",
		}))
	}

	// create mail state
	res, mailState := repositories.CreateMailState(model.MailState{
		AdminId:   user.Id,
		UserId:    user.Id,
		MailType:  Common.MailType_MAIL_TYPE_ADMIN,
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
	err = sendAdminInvite(user.Email, mailState.Id)

	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	// create redeem token for auto login when email valid succeed
	res, _ = repositories.CreateRedeemToken(model.RedeemToken{
		UserId:      mailState.UserId.String(),
		MailStateId: mailState.Id.String(),
	})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// reply
	reply.UrlKey = mailState.Id.String()
	reply.UserId = mailState.UserId.String()

	return reply, status.Errorf(codes.OK, "")
}

func ConfirmRegisterEmail(ctx context.Context, req *Auth.ConfirmRegisterEmailRequest) (*Common.CommonReply, error) {
	reply := &Common.CommonReply{}

	// check param
	if !utils.ValidId(req.UrlKey) {
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

	if data.MailState == Common.MailState_MAIL_STATE_SUCCESS {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code: Common.ErrorCodes_EMAIL_IS_VALIDATED,
		}))
	}

	if data.MailState != Common.MailState_MAIL_STATE_INVITED {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:      Common.ErrorCodes_UNKNOWN_ERROR,
			ReturnMsg: "mail state err",
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

	res, user := repositories.UpdateUser(model.User{
		Id: data.UserId,
	}, model.User{
		EmailValid: &valid,
	})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	utils.PrintObj(user, "user")

	// check email valid
	if !*user.EmailValid {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_USER_EMAIL_INVAILD,
			InternalMsg: res.Error.Error(),
		}))
	}

	return reply, status.Errorf(codes.OK, "")
}

func ResendRegisterEmail(ctx context.Context, req *Auth.ResendRegisterEmailRequest) (*Auth.ResendRegisterEmailReply, error) {
	reply := &Auth.ResendRegisterEmailReply{}

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

	// check mail state
	if mailStateQuery.MailState == Common.MailState_MAIL_STATE_SUCCESS {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code: Common.ErrorCodes_EMAIL_IS_VALIDATED,
		}))
	}

	if mailStateQuery.MailState != Common.MailState_MAIL_STATE_INVITED {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:      Common.ErrorCodes_UNKNOWN_ERROR,
			ReturnMsg: "mail state err",
		}))
	}

	// check 12hr resend
	expire, _ := utils.CheckExpire(12*time.Minute, mailStateQuery.UpdateAt) // if expire means can resend
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
	err = sendAdminInvite(user.Email, mailStateQuery.Id)
	if err != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	reply.UrlKey = mailStateQuery.Id.String()

	return reply, status.Errorf(codes.OK, "")
}

func CheckRegisterConfirmState(ctx context.Context, req *Auth.CheckRegisterConfirmStateRequest) (*Auth.CheckRegisterConfirmStateReply, error) {
	reply := &Auth.CheckRegisterConfirmStateReply{}

	utils.PrintObj(req.UserId, "CheckRegisterConfirmState")

	// check param
	if !utils.ValidId(req.UserId) {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
	}

	// check mail state is invited
	whereParam := model.MailState{
		UserId: utils.ParseUUID(req.UserId),
	}

	res, data := repositories.QueryMailState(whereParam)

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	reply.MailState = data.MailState

	return reply, status.Errorf(codes.OK, "")
}

func RedeemToken(ctx context.Context, req *Auth.RedeemTokenRequest) (*Auth.AuthReply, error) {
	reply := &Auth.AuthReply{}

	// check param
	if !utils.ValidId(req.UserId) && !utils.ValidId(req.MailStateId) {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
	}

	// check reg mail state is success
	res, data := repositories.QueryMailStateByIdOrUserId(model.MailState{
		Id:     utils.ParseUUID(req.MailStateId),
		UserId: utils.ParseUUID(req.UserId),
	})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// check valid success first before redeem token
	if data.MailState != Common.MailState_MAIL_STATE_SUCCESS {
		utils.PrintObj("check valid success first before redeem token = fail")
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code: Common.ErrorCodes(Common.ErrorCodes_UNKNOWN_ERROR),
		}))
	}

	// check redeem info
	result, redeem := repositories.QueryRedeemToken(model.RedeemToken{
		UserId:      req.UserId,
		MailStateId: req.MailStateId,
	})

	if result.Error != nil {
		if utils.IsErrorNotFound(result.Error) {
			return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
				Code: Common.ErrorCodes_DATA_NOT_FOUND,
			}))
		} else {
			return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
				Code:        Common.ErrorCodes_UNKNOWN_ERROR,
				InternalMsg: result.Error.Error(),
			}))
		}
	}

	// delete redeem info
	result, _ = repositories.DeleteRedeemToken(model.RedeemToken{
		Id: redeem.Id,
	}, model.RedeemToken{})

	if result.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: result.Error.Error(),
		}))
	}

	utils.PrintObj(redeem, "redeem")

	// create token
	token, expireTime, err := UpdateOrCreateTokenInRedis(redeem.UserId)

	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	reply.Token = "Bearer " + token
	reply.ExpireDate = utils.ParseDateToString(expireTime)
	reply.Code = Common.ErrorCodes_SUCCESS
	reply.UserId = redeem.UserId

	return reply, status.Errorf(codes.OK, "")
}

func sendAdminInvite(email string, mailStateId uuid.UUID) error { //, redeemCode string
	domain := utils.GetDomain()
	url := domain + utils.GetEnv("CONFIRM_REGISTER_ROUTE") + mailStateId.String()
	err := EmailUtils.SendEmail(email, EmailUtils.GetTemplateRegister(url))

	if err != nil {
		return err
	}

	return nil
}
