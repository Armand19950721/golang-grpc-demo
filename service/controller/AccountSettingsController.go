package controller

import (
	"context"
	"service/protos/AccountSettings"
	"service/protos/WebServices"
	"service/services/AccountSettingsServices"
)

type AccountSettingsController struct {
	WebServices.UnimplementedAccountSettingsServicesServer
}

func (company *AccountSettingsController) GetBasicInfo(ctx context.Context, req *AccountSettings.GetBasicInfoRequest) (*AccountSettings.GetBasicInfoReply, error) {
	return AccountSettingsServices.GetBasicInfo(ctx, req)
}

func (company *AccountSettingsController) UpdateBasicInfo(ctx context.Context, req *AccountSettings.UpdateBasicInfoRequest) (*AccountSettings.UpdateBasicInfoReply, error) {
	return AccountSettingsServices.UpdateBasicInfo(ctx, req)
}

func (company *AccountSettingsController) GetLoginInfo(ctx context.Context, req *AccountSettings.GetLoginInfoRequest) (*AccountSettings.GetLoginInfoReply, error) {
	return AccountSettingsServices.GetLoginInfo(ctx, req)
}

func (company *AccountSettingsController) UpdateLoginInfo(ctx context.Context, req *AccountSettings.UpdateLoginInfoRequest) (*AccountSettings.UpdateLoginInfoReply, error) {
	return AccountSettingsServices.UpdateLoginInfo(ctx, req)
}
