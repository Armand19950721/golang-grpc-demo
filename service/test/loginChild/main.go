package loginChild

import (
	"service/protos/Auth"
	"service/protos/WebServices"
	"service/test/testUtils"
	"service/utils"
)

var (
	c = WebServices.NewAuthServicesClient(testUtils.Conn)
)

func Flow() { //repeat at user test child flow
	ctx := testUtils.GetCtx()
	utils.PrintTitle("Login child")
	basicInfo := testUtils.GetBasicInfo()

	r, err := c.Login(ctx, &Auth.LoginRequest{
		Email:    basicInfo.Child.Account,
		Password: basicInfo.Child.Pwd,
	})

	testUtils.DisplayResult(r, err, false)
	basicInfo.Child.Token = r.Token

	testUtils.SetBasicInfo(basicInfo)
}
