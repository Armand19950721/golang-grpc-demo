package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"service/controller"
	"service/protos/Common"
	"service/protos/WebServices"
	"service/utils"
	"service/utils/interceptor"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// check info
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.FailedPrecondition, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_UNKNOWN_ERROR, ReturnMsg: "missing metadata"}))
	}

	interceptor.DisplayAndSaveInfo(req, info, md, ctx)

	// check and get user data
	checkState, newCtx := interceptor.Valid(md, info, ctx)
	if checkState != Common.ErrorCodes_SUCCESS {
		return nil, status.Errorf(codes.PermissionDenied, utils.GetError(utils.ErrorType{Code: checkState}))
	}

	// handle panic
	defer func() {
		err := recover()
		if err != nil {
			utils.PrintObj(err, "recover panic err")
			utils.SaveLog("panic", err)
		}
	}()

	// time display
	start := time.Now() // before execute
	message, err := handler(newCtx, req)
	elapsed := time.Since(start) // after execute

	log.Printf("spending time %s", elapsed)

	// unknow err handle
	if err != nil {
		return nil, err
	}
	if message != nil {
		utils.PrintObj(message, "message")
	}

	return message, err
}

func InitGrpcServer() {
	lis, err := net.Listen("tcp", utils.GetEnv("GRPC_SERVER_IP"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))

	WebServices.RegisterEchoServicesServer(s, &controller.EchoController{})
	WebServices.RegisterAuthServicesServer(s, &controller.AuthController{})
	WebServices.RegisterUserServicesServer(s, &controller.UserController{})
	WebServices.RegisterAccountSettingsServicesServer(s, &controller.AccountSettingsController{})

	fmt.Printf("grpc server start at %s \n", lis.Addr().String())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}

}
