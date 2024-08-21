package thirdParty

import (
	"service/protos/ThirdPartyArContent"
	"service/protos/ThirdPartyAuth"
	"service/protos/ThirdPartyWebServices"
	"service/test/testUtils"
	"service/utils"
)

var (
	authConn      = ThirdPartyWebServices.NewThirdPartyAuthServicesClient(testUtils.ConnThirdParty)
	arContentConn = ThirdPartyWebServices.NewThirdPartyArContentServicesClient(testUtils.ConnThirdParty)
)

func MainProcess() {
	GetThirdPartyAuth()
	CreateArContent()
}

func GetThirdPartyAuth() {
	// get token
	ctx := testUtils.GetThirdPartyCtx()
	utils.PrintTitle("GetThirdPartyAuth")

	defaultPinkoiClientId := "0965ed36-06a0-4a69-96f8-885540151276"
	defaultPinkoiClientSecret := "sMjBEAbvaW6l8LzN07xjx2A6kfRlAQt134fgl57f8i6Uyq3jwHg="

	req := &ThirdPartyAuth.VerificationRequest{}
	req.ClientId = defaultPinkoiClientId
	req.ClientSecret = defaultPinkoiClientSecret
	req.GrantType = "bearer"

	utils.PrintTitle("Verification")

	r, err := authConn.Verification(ctx, req)
	testUtils.DisplayResultThirdParty(r, err, false)

	// save token
	basicInfo := testUtils.GetBasicInfo()
	basicInfo.Token = r.AccessToken
	basicInfo.AdminId = defaultPinkoiClientId

	testUtils.SetBasicInfo(basicInfo)

	utils.PrintTitle("Verification auth fail")

	req.ClientId = "0965ed36-06a0-4a69-96f8-885540151276"
	req.ClientSecret = "123"
	req.GrantType = "123"

	r, err = authConn.Verification(ctx, req)
	testUtils.DisplayResultThirdParty(r, err, true)

	utils.PrintTitle("Verification param fail")

	req.ClientId = "123"
	req.ClientSecret = "123"
	req.GrantType = "123"

	r, err = authConn.Verification(ctx, req)
	testUtils.DisplayResultThirdParty(r, err, true)
}

func CreateArContent() {
	ctx := testUtils.GetThirdPartyCtx()
	req := &ThirdPartyArContent.CreateArContentRequest{}

	utils.PrintTitle("CreateArContent fail")
	req.Name = ""

	r, err := arContentConn.CreateArContent(ctx, req)
	testUtils.DisplayResultThirdParty(r, err, true)

	utils.PrintTitle("CreateArContent glasses")
	req.Name = "ArContent_" + testUtils.RandomIntString()
	req.TemplateType = ThirdPartyArContent.ArContentTemplateType_TEMPLATE_GLASSES

	r, err = arContentConn.CreateArContent(ctx, req)
	testUtils.DisplayResultThirdParty(r, err, false)

	utils.PrintTitle("CreateArContent earring")
	req.Name = "ArContent_" + testUtils.RandomIntString()
	req.TemplateType = ThirdPartyArContent.ArContentTemplateType_TEMPLATE_EARRING

	r, err = arContentConn.CreateArContent(ctx, req)
	testUtils.DisplayResultThirdParty(r, err, false)

	utils.PrintTitle("GetArLink")

	req2 := &ThirdPartyArContent.GetArLinkRequest{}
	req2.ArContentId = r.ArContentId

	r2, err := arContentConn.GetArLink(ctx, req2)
	testUtils.DisplayResultThirdParty(r2, err, false)

	utils.PrintTitle("GetArLink id fail")

	req2.ArContentId = "fake"

	r2, err = arContentConn.GetArLink(ctx, req2)
	testUtils.DisplayResultThirdParty(r2, err, true)

}
