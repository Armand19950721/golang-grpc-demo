package accountSettings

import (
	"service/protos/AccountSettings"
	"service/protos/WebServices"
	"service/test/testUtils"
	"service/utils"
)

var (
	c = WebServices.NewAccountSettingsServicesClient(testUtils.Conn)
)

func Flow() {
	ctx := testUtils.GetCtx()
	utils.PrintTitle("UpdateBasicInfo")

	r, err := c.UpdateBasicInfo(ctx, &AccountSettings.UpdateBasicInfoRequest{
		Model: &AccountSettings.BasicInfoModel{
			FirstName:   "FirstName_" + testUtils.RandomIntString(),
			LastName:    "LastName_" + testUtils.RandomIntString(),
			CompanyName: "CompanyName_" + testUtils.RandomIntString(),
			GuiNumber:   "GuiNumber_" + testUtils.RandomIntString(),
		},
	})

	testUtils.DisplayResult(r, err, false)

	///
	utils.PrintTitle("GetBasicInfo")

	r1, err := c.GetBasicInfo(ctx, &AccountSettings.GetBasicInfoRequest{})

	testUtils.DisplayResult(r1, err, false)

	///
	utils.PrintTitle("GetLoginInfo one")

	r2, err := c.GetLoginInfo(ctx, &AccountSettings.GetLoginInfoRequest{})

	testUtils.DisplayResult(r2, err, false)

	///
	utils.PrintTitle("UpdateLoginInfo")

	r3, err := c.UpdateLoginInfo(ctx, &AccountSettings.UpdateLoginInfoRequest{
		Model: &AccountSettings.LoginInfoModel{
			AccountName: "AccountName_" + testUtils.RandomIntString(),
		},
	})

	testUtils.DisplayResult(r3, err, false)

	///
	utils.PrintTitle("GetLoginInfo two")

	r2, err = c.GetLoginInfo(ctx, &AccountSettings.GetLoginInfoRequest{})

	testUtils.DisplayResult(r2, err, false)
}
