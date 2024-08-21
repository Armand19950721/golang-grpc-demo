package controller

import (
	"service/protos/Common"
	"service/protos/User"
	"service/protos/WebServices"
	"service/services/UserService"

	"context"
)

type UserController struct {
	WebServices.UnimplementedUserServicesServer
}

func (user *UserController) AddUserChild(ctx context.Context, req *User.AddUserChildRequest) (*Common.CommonReply, error) {
	return UserService.AddUserChild(ctx, req)
}

func (user *UserController) UpdateUserChild(ctx context.Context, req *User.UpdateUserChildRequest) (*Common.CommonReply, error) {
	return UserService.UpdateUserChild(ctx, req)
}

func (user *UserController) DeleteUserChild(ctx context.Context, req *User.DeleteUserChildRequest) (*Common.CommonReply, error) {
	return UserService.DeleteUserChild(ctx, req)
}

func (user *UserController) GetUserChildList(ctx context.Context, req *User.GetUserChildListRequest) (*User.GetUserChildListReply, error) {
	return UserService.GetUserChildList(ctx, req)
}

func (user *UserController) ChangePassword(ctx context.Context, req *User.ChangePasswordRequest) (*Common.CommonReply, error) {
	return UserService.ChangePassword(ctx, req)
}

func (user *UserController) ForgotPassword(ctx context.Context, req *User.ForgotPasswordRequest) (*Common.CommonReply, error) {
	return UserService.ForgotPassword(ctx, req)
}

func (user *UserController) UpdateUserProgram(ctx context.Context, req *User.UpdateUserProgramRequest) (*Common.CommonReply, error) {
	return UserService.UpdateUserProgram(ctx, req)
}

func (user *UserController) GetConfirmUserChildInfo(ctx context.Context, req *User.GetConfirmUserChildInfoRequest) (*User.GetConfirmUserChildInfoReply, error) {
	return UserService.GetConfirmUserChildInfo(ctx, req)
}

func (user *UserController) ConfirmUserChild(ctx context.Context, req *User.ConfirmUserChildRequest) (*Common.CommonReply, error) {
	return UserService.ConfirmUserChild(ctx, req)
}

func (user *UserController) ResendChildInvite(ctx context.Context, req *User.ResendChildInviteRequest) (*User.ResendChildInviteReply, error) {
	return UserService.ResendChildInvite(ctx, req)
}

func (user *UserController) CancelChildInvite(ctx context.Context, req *User.CancelChildInviteRequest) (*Common.CommonReply, error) {
	return UserService.CancelChildInvite(ctx, req)
}
