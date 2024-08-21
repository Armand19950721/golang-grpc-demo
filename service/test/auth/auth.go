package auth

import (
	"service/model"
	"service/repositories"
	"service/test/testUtils"
	"service/utils"
	"time"

	"service/protos/Auth"
	"service/protos/Common"
	"service/protos/Echo"
	"service/protos/WebServices"
)

var (
	c = WebServices.NewAuthServicesClient(testUtils.Conn)
	e = WebServices.NewEchoServicesClient(testUtils.Conn)
)

func loginUser(email string) {
	basicInfo := testUtils.GetBasicInfo()

	ctx := testUtils.GetCtx()
	utils.PrintTitle("Login")
	model := &Auth.LoginRequest{
		Email:    email,
		Password: "123456",
	}
	r, err := c.Login(ctx, model)

	testUtils.DisplayResult(r, err, false)

	basicInfo.Token = r.Token
	basicInfo.AdminId = r.UserId
	testUtils.SetBasicInfo(basicInfo)
}

func Logout() {
	ctx := testUtils.GetCtx()

	utils.PrintTitle("logout")

	r, err := c.Logout(ctx, &Auth.LogoutRequest{})

	testUtils.DisplayResult(r, err, false)
}

func RegDevGuest() {
	ctx := testUtils.GetCtx()
	programId := testUtils.GetBasicInfo().ProgramId
	utils.PrintTitle("register")
	model := &Auth.RegisterRequest{
		Email:       "yuei@spe3d.co",
		Password:    "1234567890",
		Name:        "yuei",
		PhoneNumber: "0912345678",
		ProgramId:   programId,
	}
	r, err := c.Register(ctx, model)

	utils.PrintObj(model)
	testUtils.DisplayResult(r, err, false)
}

func LoginDevGuest() {
	ctx := testUtils.GetCtx()
	utils.PrintTitle("Login")
	model := &Auth.LoginRequest{
		Email:    "yuei@spe3d.co",
		Password: "1234567890",
	}
	r, err := c.Login(ctx, model)

	utils.PrintObj(model)
	testUtils.DisplayResult(r, err, false)
}

func PinkoiAuthTest() {
	LoginDevPinkoi()
}

func RegDevPinkoi() {
	utils.PrintTitle("RegDevPinkoi")

	ctx := testUtils.GetCtx()
	programId := testUtils.GetBasicInfo().ProgramId
	utils.PrintTitle("register")
	model := &Auth.RegisterRequest{
		Email:       "pinkoi@spe3d.co",
		Password:    "123456",
		Name:        "pinkoi",
		PhoneNumber: "0912345678",
		ProgramId:   programId,
	}
	r, err := c.Register(ctx, model)

	utils.PrintObj(model)
	testUtils.DisplayResult(r, err, false)

	utils.PrintTitle("confirm")
	r1, err := c.ConfirmRegisterEmail(ctx, &Auth.ConfirmRegisterEmailRequest{
		UrlKey: r.UrlKey,
	})

	testUtils.DisplayResult(r1, err, false)
}

func LoginDevPinkoi() {
	utils.PrintTitle("RegDevPinkoi")

	ctx := testUtils.GetCtx()
	utils.PrintTitle("Login")
	model := &Auth.LoginRequest{
		Email:    "pinkoi@spe3d.co",
		Password: "123456",
	}
	r, err := c.Login(ctx, model)

	utils.PrintObj(model)
	testUtils.DisplayResult(r, err, false)
}

func RegAndEmailFlow() {
	ctx := testUtils.GetCtx()
	programId := testUtils.GetBasicInfo().ProgramId

	///
	utils.PrintTitle("register fail")

	acc := testUtils.RandomIntString() + "@metacommerce.test.com"
	modelItem := &Auth.RegisterRequest{
		Email:       acc,
		Password:    "12345678901",
		Name:        "armand",
		PhoneNumber: "0912" + testUtils.RandomIntString(),
		ProgramId:   "9b26a91d6bb",
	}

	r, err := c.Register(ctx, modelItem)
	testUtils.DisplayResult(r, err, true)

	///
	utils.PrintTitle("register")

	modelItem = &Auth.RegisterRequest{
		Email:       acc,
		Password:    "12345678901",
		Name:        "armand",
		PhoneNumber: "0912" + testUtils.RandomIntString(),
		ProgramId:   programId,
	}

	r, err = c.Register(ctx, modelItem)
	testUtils.DisplayResult(r, err, false)

	// set info
	info := testUtils.AdminBasicInfo{
		Account:   modelItem.Email,
		Pwd:       modelItem.Password,
		ProgramId: programId,
		AdminId:   r.UserId,
	}
	testUtils.SetBasicInfo(info)

	///

	utils.PrintTitle("CheckRegisterConfirmState should get INVITED")

	re2, err := c.CheckRegisterConfirmState(ctx, &Auth.CheckRegisterConfirmStateRequest{
		UserId: r.UserId,
	})

	testUtils.DisplayResult(re2, err, false)

	if re2.MailState != Common.MailState_MAIL_STATE_INVITED {
		panic("auth MailState err")
	}

	///
	utils.PrintTitle("ResendRegisterEmail fail in 12hr")

	r2, err := c.ResendRegisterEmail(ctx, &Auth.ResendRegisterEmailRequest{
		UserId: r.UserId,
	})

	testUtils.DisplayResult(r2, err, true)

	///

	setMailStateUpdateAt(r.UserId, -13)

	///
	utils.PrintTitle("ResendRegisterEmail")

	r2, err = c.ResendRegisterEmail(ctx, &Auth.ResendRegisterEmailRequest{
		UserId: r.UserId,
	})

	testUtils.DisplayResult(r2, err, false)
	///

	utils.PrintTitle("Login before email confirm")

	r3, err := c.Login(ctx, &Auth.LoginRequest{
		Email:    modelItem.Email,
		Password: modelItem.Password,
	})

	testUtils.DisplayResult(r3, err, false)

	if r3.Code != Common.ErrorCodes_USER_EMAIL_INVAILD {
		panic("login return fail")
	}

	///
	mailStateId := r.UrlKey
	utils.PrintTitle("redeem token before email confirm")

	r3, err = c.RedeemToken(ctx, &Auth.RedeemTokenRequest{
		MailStateId: mailStateId,
	})

	testUtils.DisplayResult(r3, err, true)

	///
	utils.PrintTitle("confirm reg")

	r1, err := c.ConfirmRegisterEmail(ctx, &Auth.ConfirmRegisterEmailRequest{
		UrlKey: r.UrlKey,
	})

	testUtils.DisplayResult(r1, err, false)

	///
	utils.PrintTitle("confirm reg again err => EMAIL_IS_VALIDATED")

	r1, err = c.ConfirmRegisterEmail(ctx, &Auth.ConfirmRegisterEmailRequest{
		UrlKey: r.UrlKey,
	})

	testUtils.DisplayResult(r1, err, true)

	errParseCode, _, _ := testUtils.ParseErrorCode(err)

	if errParseCode != Common.ErrorCodes_EMAIL_IS_VALIDATED {
		panic("resend err is wrong")
	}

	///
	utils.PrintTitle("resend when already confirm err => EMAIL_IS_VALIDATED")

	r2, err = c.ResendRegisterEmail(ctx, &Auth.ResendRegisterEmailRequest{
		UserId: r.UserId,
	})

	testUtils.DisplayResult(r2, err, true)

	errParseCode, _, _ = testUtils.ParseErrorCode(err)

	if errParseCode != Common.ErrorCodes_EMAIL_IS_VALIDATED {
		panic("resend err is wrong")
	}

	///

	utils.PrintTitle("CheckRegisterConfirmState should get SUCCESS")

	re2, err = c.CheckRegisterConfirmState(ctx, &Auth.CheckRegisterConfirmStateRequest{
		UserId: r.UserId,
	})

	testUtils.DisplayResult(re2, err, false)

	if re2.MailState != Common.MailState_MAIL_STATE_SUCCESS {
		panic("auth MailState err")
	}

	///
	utils.PrintTitle("ResendRegisterEmail state err")

	r2, err = c.ResendRegisterEmail(ctx, &Auth.ResendRegisterEmailRequest{
		UserId: r.UserId,
	})

	testUtils.DisplayResult(r2, err, true)

	///
	utils.PrintTitle("redeem token by mail state id")

	r3, err = c.RedeemToken(ctx, &Auth.RedeemTokenRequest{
		MailStateId: mailStateId,
	})

	testUtils.DisplayResult(r3, err, false)
	// ///
	// utils.PrintTitle("redeem token by user id")

	// r3, err = c.RedeemToken(ctx, &Auth.RedeemTokenRequest{
	// 	UserId: r.UserId,
	// })

	// testUtils.DisplayResult(r3, err, false)

	/// save token in local
	basicInfo := testUtils.GetBasicInfo()

	basicInfo.Token = r3.Token
	basicInfo.AdminId = r3.UserId
	testUtils.SetBasicInfo(basicInfo)

	///
	utils.PrintTitle("redeem token by user id fail, already redeemed")

	r3, err = c.RedeemToken(ctx, &Auth.RedeemTokenRequest{
		UserId: r.UserId,
	})

	testUtils.DisplayResult(r3, err, true)

	///

	utils.PrintTitle("echo token")

	rR, err := e.EchoToken(ctx, &Echo.EchoRequest{
		Msg: "hi~~",
	})
	testUtils.DisplayResult(rR, err, false)

	///

	utils.PrintTitle("Login test")

	r3, err = c.Login(ctx, &Auth.LoginRequest{
		Email:    modelItem.Email,
		Password: modelItem.Password,
	})

	testUtils.DisplayResult(r3, err, false)

}

func setMailStateUpdateAt(userId string, hr int) {
	utils.PrintObj(userId+","+utils.ToString(hr), "setMailStateUpdateAt")

	///
	utils.PrintTitle("set main state update time is " + utils.ToString(hr))

	where := model.MailState{UserId: utils.ParseUUID(userId)}

	res, data := repositories.QueryMailState(where)

	if res.Error != nil {
		panic(res.Error.Error())
	}

	data.UpdateAt = data.UpdateAt.Add(time.Duration(hr) * time.Hour)

	utils.PrintObj(data, "data")

	res, dataUpdate := repositories.UpdateMailState(where, data)

	if res.Error != nil {
		panic(res.Error.Error())
	}
	utils.PrintObj(dataUpdate, "dataUpdate")
}

func Flow() {
	ctx := testUtils.GetCtx()

	RegAndEmailFlow()

	basicInfo := testUtils.GetBasicInfo()

	utils.PrintTitle("Login after email confirm")

	r3, err := c.Login(ctx, &Auth.LoginRequest{
		Email:    basicInfo.Account,
		Password: basicInfo.Pwd,
	})

	if r3.Code != Common.ErrorCodes_SUCCESS {
		panic("login return fail")
	}

	testUtils.DisplayResult(r3, err, false)

	basicInfo.Token = r3.Token
	testUtils.SetBasicInfo(basicInfo)

	utils.PrintTitle("Login fail pwd")

	r3, err = c.Login(ctx, &Auth.LoginRequest{
		Email:    basicInfo.Account,
		Password: "12345999999999999999999",
	})

	testUtils.DisplayResult(r3, err, true)

	utils.PrintTitle("Login fail account")

	r3, err = c.Login(ctx, &Auth.LoginRequest{
		Email:    "12e3123@metacommerce.test.com",
		Password: "12345999999999999999999",
	})

	testUtils.DisplayResult(r3, err, true)

	utils.PrintTitle("Login fail param")

	r3, err = c.Login(ctx, &Auth.LoginRequest{
		Email:    "123123",
		Password: "33",
	})

	testUtils.DisplayResult(r3, err, true)
}
