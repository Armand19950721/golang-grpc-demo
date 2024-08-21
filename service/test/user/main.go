package user

import (
	"errors"
	"service/model"
	"service/protos/Auth"
	"service/protos/Common"
	"service/protos/Permission"
	"service/protos/User"
	"service/protos/WebServices"
	"service/repositories"
	"service/test/testUtils"
	"service/utils"
	"time"
)

var (
	u = WebServices.NewUserServicesClient(testUtils.Conn)
	p = WebServices.NewPermissionServicesClient(testUtils.Conn)
	a = WebServices.NewAuthServicesClient(testUtils.Conn)
	r = &Common.CommonReply{}
)

func UpdateUserProgram() {
	ctx := testUtils.GetCtx()
	utils.PrintTitle("UpdateUserProgram")
	basicInfo := testUtils.GetBasicInfo()

	r, err := u.UpdateUserProgram(ctx, &User.UpdateUserProgramRequest{
		ProgramId: basicInfo.ProgramId,
	})

	testUtils.DisplayResult(r, err, false)
}

func UserChildFlow() {
	ExpireUnitTest()
	UserFlowOther()
	UserChildFlowRegular()
	UserChildFlowDeleteInvite()
}

func ExpireUnitTest() {
	utils.PrintTitle("ExpireUnitTest")
	m13Hr := time.Now().Add(time.Duration(-13) * time.Hour)
	m10Hr := time.Now().Add(time.Duration(-10) * time.Hour)
	expire, endDate := utils.CheckExpire(12*time.Hour, m13Hr)

	utils.PrintObj(expire)
	utils.PrintObj(endDate)
	if !expire {
		panic("m13Hr test fail")
	}

	expire, endDate = utils.CheckExpire(12*time.Hour, m10Hr)

	utils.PrintObj(expire)
	utils.PrintObj(endDate)
	if expire {
		panic("m10Hr test fail")
	}
}

func UserFlowOther() {
	ctx := testUtils.GetCtx()
	permissionGroupRes := GetPermissionGroupId()
	basicInfo := testUtils.GetBasicInfo()

	basicInfo.Child.Account = testUtils.RandomIntString() + "@metacommerce.test.com"

	///
	utils.PrintTitle("add user child")
	modelUser := &User.UserChildModel{
		Email:             basicInfo.Child.Account,
		Name:              "name_" + testUtils.RandomIntString(),
		PermissionGroupId: permissionGroupRes.Models[0].PermissionGroupId,
	}

	r, err := u.AddUserChild(ctx, &User.AddUserChildRequest{
		Model: modelUser,
	})

	testUtils.DisplayResult(r, err, false)

	///
	utils.PrintTitle("add user child repeat")
	r, err = u.AddUserChild(ctx, &User.AddUserChildRequest{
		Model: modelUser,
	})

	testUtils.DisplayResult(r, err, true)

	///
	utils.PrintTitle("get user child list")
	pageInfo := Common.PageInfoRequest{
		CurrentPageNum: 1,
		PageItemCount:  10,
	}
	rList, err := u.GetUserChildList(ctx, &User.GetUserChildListRequest{
		PageInfo: &pageInfo,
	})

	testUtils.DisplayResult(rList, err, false)
	userId := rList.Models[0].UserId

	///
	utils.PrintTitle("Login child fail. email vaild fail")

	r7, err := a.Login(ctx, &Auth.LoginRequest{
		Email:    basicInfo.Child.Account,
		Password: "1234455",
	})

	testUtils.DisplayResult(r7, err, true)

	///
	utils.PrintTitle("set main state update time is -13 hr")
	where := model.MailState{UserId: utils.ParseUUID(userId)}

	res, data := repositories.QueryMailState(where)

	if res.Error != nil {
		panic(res.Error.Error())
	}

	data.UpdateAt = data.UpdateAt.Add(time.Duration(-13) * time.Hour)

	utils.PrintObj(data, "data")

	res, dataUpdate := repositories.UpdateMailState(where, data)

	if res.Error != nil {
		panic(res.Error.Error())
	}
	utils.PrintObj(dataUpdate, "dataUpdate")

	///
	utils.PrintTitle("ResendChildInvite")

	r2, err := u.ResendChildInvite(ctx, &User.ResendChildInviteRequest{
		UserId: userId,
	})

	testUtils.DisplayResult(r2, err, false)

	///
	utils.PrintTitle("ConfirmUserChild")

	basicInfo.Child.Pwd = "123456"

	r3, err := u.ConfirmUserChild(ctx, &User.ConfirmUserChildRequest{
		UrlKey:   r2.UrlKey,
		Password: basicInfo.Child.Pwd,
	})

	testUtils.DisplayResult(r3, err, false)

	///
	utils.PrintTitle("ResendChildInvite state err")

	r2, err = u.ResendChildInvite(ctx, &User.ResendChildInviteRequest{
		UserId: userId,
	})

	testUtils.DisplayResult(r2, err, true)

	///
	utils.PrintTitle("Login child")

	r6, err := a.Login(ctx, &Auth.LoginRequest{
		Email:    basicInfo.Child.Account,
		Password: basicInfo.Child.Pwd,
	})

	testUtils.DisplayResult(r6, err, false)
	basicInfo.Child.Token = r6.Token

	testUtils.SetBasicInfo(basicInfo)

	///
	utils.PrintTitle("update user child")
	modelTwo := &User.UserChildModel{
		UserId:            userId,
		Name:              " name chenge  name chenge  name chenge ",
		PermissionGroupId: permissionGroupRes.Models[0].PermissionGroupId,
	}
	r, err = u.UpdateUserChild(ctx, &User.UpdateUserChildRequest{
		Model: modelTwo,
	})
	testUtils.DisplayResult(r, err, false)

	///
	utils.PrintTitle("delete user child")
	utils.PrintObj(modelTwo.UserId, "delete id")
	r, err = u.DeleteUserChild(ctx, &User.DeleteUserChildRequest{
		Model: modelTwo,
	})

	testUtils.DisplayResult(r, err, false)

	///
	utils.PrintTitle("check cancel")
	res, _ = repositories.QueryUserIncludeDel(model.User{Id: utils.ParseUUID(userId)})
	if !utils.IsErrorNotFound(res.Error) {
		panic("user child is not delete success")
	}

	testUtils.DisplayResult(r, err, false)
}

func UserChildFlowDeleteInvite() {
	ctx := testUtils.GetCtx()
	permissionGroupRes := GetPermissionGroupId()

	///
	utils.PrintTitle("add user child")
	modelUser := &User.UserChildModel{
		Email:             testUtils.RandomIntString() + "@metacommerce.test.com",
		Name:              "name_" + testUtils.RandomIntString(),
		PermissionGroupId: permissionGroupRes.Models[0].PermissionGroupId,
	}

	r, err := u.AddUserChild(ctx, &User.AddUserChildRequest{
		Model: modelUser,
	})

	testUtils.DisplayResult(r, err, false)

	///
	utils.PrintTitle("get user child list")
	pageInfo := Common.PageInfoRequest{
		CurrentPageNum: 1,
		PageItemCount:  10,
	}
	rList, err := u.GetUserChildList(ctx, &User.GetUserChildListRequest{
		PageInfo: &pageInfo,
	})

	testUtils.DisplayResult(rList, err, false)

	///
	utils.PrintTitle("CancelChildInvite")

	userId := rList.Models[0].UserId
	r1, err := u.CancelChildInvite(ctx, &User.CancelChildInviteRequest{
		UserId: userId,
	})

	testUtils.DisplayResult(r1, err, false)

	///
	utils.PrintTitle("check cancel")
	res, _ := repositories.QueryUserIncludeDel(model.User{Id: utils.ParseUUID(userId)})
	if !utils.IsErrorNotFound(res.Error) {
		panic("user child is not delete success")
	}
	res, _ = repositories.QueryMailStateIncludeDel(model.MailState{Id: utils.ParseUUID(userId)})
	if !utils.IsErrorNotFound(res.Error) {
		panic("MailState is not delete success")
	}
}

func UserChildFlowRegular() {
	ctx := testUtils.GetCtx()
	permissionGroupRes := GetPermissionGroupId()

	///
	utils.PrintTitle("add user child")
	modelUser := &User.UserChildModel{
		Email:             testUtils.RandomIntString() + "@metacommerce.test.com",
		Name:              "name_" + testUtils.RandomIntString(),
		PermissionGroupId: permissionGroupRes.Models[0].PermissionGroupId,
	}

	r, err := u.AddUserChild(ctx, &User.AddUserChildRequest{
		Model: modelUser,
	})

	testUtils.DisplayResult(r, err, false)

	///
	utils.PrintTitle("get user child list")
	pageInfo := Common.PageInfoRequest{
		CurrentPageNum: 1,
		PageItemCount:  10,
	}
	rList, err := u.GetUserChildList(ctx, &User.GetUserChildListRequest{
		PageInfo: &pageInfo,
	})

	testUtils.DisplayResult(rList, err, false)

	if rList.Models[0].MailState != Common.MailState_MAIL_STATE_INVITED {
		testUtils.DisplayResult(rList.Models[0].MailState, errors.New("check MailState fail"), false)
	}

	if rList.Models[0].MailExpireDate == "" {
		testUtils.DisplayResult(rList.Models[0].MailExpireDate, errors.New("check MailExpireDate fail"), false)
	}

	userId := rList.Models[0].UserId

	///
	utils.PrintTitle("ResendChildInvite fail 12hr")

	r2, err := u.ResendChildInvite(ctx, &User.ResendChildInviteRequest{
		UserId: userId,
	})

	testUtils.DisplayResult(r2, err, true)

	///
	utils.PrintTitle("set main state update time is -13 hr")
	where := model.MailState{UserId: utils.ParseUUID(userId)}

	res, data := repositories.QueryMailState(where)

	if res.Error != nil {
		panic(res.Error.Error())
	}

	data.UpdateAt = data.UpdateAt.Add(time.Duration(-13) * time.Hour)

	utils.PrintObj(data, "data")

	res, dataUpdate := repositories.UpdateMailState(where, data)

	if res.Error != nil {
		panic(res.Error.Error())
	}
	utils.PrintObj(dataUpdate, "dataUpdate")

	///
	utils.PrintTitle("ResendChildInvite ok")

	r2, err = u.ResendChildInvite(ctx, &User.ResendChildInviteRequest{
		UserId: userId,
	})

	testUtils.DisplayResult(r2, err, false)

	///
	utils.PrintTitle("check mail updatetime")

	res, mailState := repositories.QueryMailState(model.MailState{
		Id: utils.ParseUUID(r2.UrlKey),
	})

	if res.Error != nil {
		testUtils.DisplayResult("", res.Error, false)
	}

	if mailState.UpdateAt == mailState.CreateAt {
		testUtils.DisplayResult([]time.Time{mailState.UpdateAt, mailState.UpdateAt}, errors.New("check UpdateAt fail"), false)
	}

	utils.PrintObj("check ok")

	///
	utils.PrintTitle("GetConfirmUserChildInfo")

	r1, err := u.GetConfirmUserChildInfo(ctx, &User.GetConfirmUserChildInfoRequest{
		UrlKey: r2.UrlKey,
	})

	testUtils.DisplayResult(r1, err, false)

	///
	utils.PrintTitle("ConfirmUserChild")

	r3, err := u.ConfirmUserChild(ctx, &User.ConfirmUserChildRequest{
		UrlKey:   r2.UrlKey,
		Password: "123456",
	})

	testUtils.DisplayResult(r3, err, false)
}

func ForgotPwd() {
	basicInfo := testUtils.GetBasicInfo()

	ctx := testUtils.GetCtx()
	utils.PrintTitle("forgotPwd")

	r, err := u.ForgotPassword(ctx, &User.ForgotPasswordRequest{
		Email: basicInfo.Account,
	})

	testUtils.DisplayResult(r, err, false)

	utils.PrintTitle("forgotPwd login")
	basicInfo = testUtils.GetBasicInfo() // get new info
	basicInfo.Pwd = basicInfo.Remark     // new pwd will save in remak field when call forgot pwd in dev env

	rR, err := a.Login(ctx, &Auth.LoginRequest{
		Email:    basicInfo.Account,
		Password: basicInfo.Pwd,
	})
	testUtils.DisplayResult(rR, err, false)
	// save token
	basicInfo.Token = rR.Token
	testUtils.SetBasicInfo(basicInfo)
	ctx = testUtils.GetCtx()

	utils.PrintTitle("forgotPwd child")

	r, err = u.ForgotPassword(ctx, &User.ForgotPasswordRequest{
		Email: basicInfo.Child.Account,
	})

	testUtils.DisplayResult(r, err, false)

	utils.PrintTitle("forgotPwd login child")
	basicInfo = testUtils.GetBasicInfo()   // get new info
	basicInfo.Child.Pwd = basicInfo.Remark // new pwd will save in remak field when call forgot pwd in dev env

	rR, err = a.Login(ctx, &Auth.LoginRequest{
		Email:    basicInfo.Child.Account,
		Password: basicInfo.Child.Pwd,
	})

	testUtils.DisplayResult(rR, err, false)
	// save token
	basicInfo.Child.Token = rR.Token
	testUtils.SetBasicInfo(basicInfo)
}

func TestSeats() {
	ctx := testUtils.GetCtx()
	utils.PrintTitle("testSeats")
	permissionGroupRes := GetPermissionGroupId()

	utils.PrintTitle("add user child")
	for i := 0; i < 4; i++ {
		model := &User.UserChildModel{
			Email:             testUtils.RandomIntString() + "@metacommerce.test.com",
			Password:          "1234569222222222",
			Name:              "name_" + testUtils.RandomIntString(),
			PermissionGroupId: permissionGroupRes.Models[0].PermissionGroupId,
		}

		r, err := u.AddUserChild(ctx, &User.AddUserChildRequest{
			Model: model,
		})

		testUtils.DisplayResult(r, err, false)

		// if err != nil {
		// 	break
		// }
	}
}

func ChangePasswordFlow() {
	ctx := testUtils.GetCtx()
	info := testUtils.GetBasicInfo()
	//
	old := info.Pwd
	new := testUtils.RandomIntString() + "123456"
	utils.PrintTitle("change pwd")
	r, err := u.ChangePassword(ctx, &User.ChangePasswordRequest{
		OldPassword: old,
		NewPassword: new,
	})

	testUtils.DisplayResult(r, err, false)

	utils.PrintTitle("change pwd and login")

	r2, err := a.Login(ctx, &Auth.LoginRequest{
		Email:    info.Account,
		Password: new,
	})

	testUtils.DisplayResult(r2, err, false)
	info.Token = r2.Token
	testUtils.SetBasicInfo(info)

	utils.PrintTitle("change pwd and login fail old")

	r3, err := a.Login(ctx, &Auth.LoginRequest{
		Email:    info.Account,
		Password: old,
	})

	testUtils.DisplayResult(r3, err, true)

	utils.PrintTitle("change pwd and login fail")

	r4, err := a.Login(ctx, &Auth.LoginRequest{
		Email:    info.Account,
		Password: "testtesttesttesttest",
	})

	testUtils.DisplayResult(r4, err, true)
}

func GetPermissionGroupId() *Permission.GetDefaultPermissionGroupListReply {
	ctx := testUtils.GetCtx()
	utils.PrintTitle("GetPermissionGroupDefaultList")

	permissionGroupRes, err := p.GetPermissionGroupDefaultList(ctx, &Permission.GetDefaultPermissionGroupListRequest{})
	testUtils.DisplayResult(r, err, false)

	return permissionGroupRes
}
