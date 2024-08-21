package controller

import (
	"service/protos/Auth"
	"service/protos/Common"
	"service/protos/WebServices"
	"service/services/AuthService"

	"context"
)

type AuthController struct {
	WebServices.UnimplementedAuthServicesServer
}

func (auth *AuthController) Login(ctx context.Context, req *Auth.LoginRequest) (*Auth.AuthReply, error) {
	return AuthService.LoginService(ctx, req)
}

func (auth *AuthController) Logout(ctx context.Context, req *Auth.LogoutRequest) (*Common.CommonReply, error) {
	return AuthService.LogoutService(ctx, req)
}

func (auth *AuthController) Register(ctx context.Context, req *Auth.RegisterRequest) (*Auth.RegisterReply, error) {
	return AuthService.RegisterService(ctx, req)
}

func (auth *AuthController) ConfirmRegisterEmail(ctx context.Context, req *Auth.ConfirmRegisterEmailRequest) (*Common.CommonReply, error) {
	return AuthService.ConfirmRegisterEmail(ctx, req)
}

func (auth *AuthController) ResendRegisterEmail(ctx context.Context, req *Auth.ResendRegisterEmailRequest) (*Auth.ResendRegisterEmailReply, error) {
	return AuthService.ResendRegisterEmail(ctx, req)
}

func (auth *AuthController) CheckRegisterConfirmState(ctx context.Context, req *Auth.CheckRegisterConfirmStateRequest) (*Auth.CheckRegisterConfirmStateReply, error) {
	return AuthService.CheckRegisterConfirmState(ctx, req)
}

func (auth *AuthController) RedeemToken(ctx context.Context, req *Auth.RedeemTokenRequest) (*Auth.AuthReply, error) {
	return AuthService.RedeemToken(ctx, req)
}
