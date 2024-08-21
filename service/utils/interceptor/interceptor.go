package interceptor

import (
	"context"
	"net"
	"strings"

	"service/protos/Common"
	permissionUtils "service/services/PermissionService"
	"service/utils"
	programUtils "service/utils/program"
	tokenUtil "service/utils/token"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

var (
	_ = godotenv.Load()
)

func getRouteInfo(info *grpc.UnaryServerInfo) (funcName, endPointType, serviceName string) {
	methodParseArr := strings.Split(info.FullMethod, "/")
	methodParseArrTwo := strings.Split(methodParseArr[1], ".")

	funcName = methodParseArr[2]
	endPointType = methodParseArrTwo[0]
	serviceName = methodParseArrTwo[1]

	utils.PrintObj([]string{funcName, endPointType, serviceName},
		" >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> getRouteInfo <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	return
}

func checkIsNeedValid(serviceName, funcName string) bool {
	// golang does not have set. so use map instead of
	tokenServiceWhiteList := map[string]string{
		"ProgramServices": "",
	}

	tokenFuncWhiteList := map[string]string{
		"Echo":                      "",
		"Login":                     "",
		"Register":                  "",
		"ForgotPassword":            "",
		"ConfirmUserChild":          "",
		"ConfirmRegisterEmail":      "",
		"ResendRegisterEmail":       "",
		"CheckRegisterConfirmState": "",
		"RedeemToken":               "",
		"GetViewerData":             "",
		"GetConfirmUserChildInfo":   "",
		"AddCount":                  "",
	}

	_, isServiceExist := tokenServiceWhiteList[serviceName]
	_, isFuncExist := tokenFuncWhiteList[funcName]

	// utils.PrintObj([]bool{isServiceExist, isFuncExist}, "in white list")

	// either rpc fun or service is true
	// will be a "pass". means no need valid
	if isServiceExist || isFuncExist {
		return false
	} else {
		return true
	}
}

func DisplayAndSaveInfo(req interface{}, info *grpc.UnaryServerInfo, md metadata.MD, ctx context.Context) {
	utils.PrintObj(req, "req")
	// utils.PrintObj(info, "info")
	// utils.PrintObj(md, "metadata")
	// utils.PrintObj(getPeerAddr(ctx), "peer addr")
	utils.PrintObj(getRealAddr(ctx), "real addr")
}

func getRealAddr(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	utils.PrintObj(ok, "ok")
	if !ok {
		return ""
	}

	rips := md.Get("x-real-ip")

	if len(rips) == 0 {
		return ""
	}

	return rips[0]
}

func getPeerAddr(ctx context.Context) string {
	var ip string
	if pr, ok := peer.FromContext(ctx); ok {
		if tcpAddr, ok := pr.Addr.(*net.TCPAddr); ok {
			ip = tcpAddr.IP.String()
		} else {
			ip = pr.Addr.String()
		}
	}
	return ip
}

func Valid(md metadata.MD, info *grpc.UnaryServerInfo, ctx context.Context) (Common.ErrorCodes, context.Context) {
	funcName, _, serviceName := getRouteInfo(info)

	if checkIsNeedValid(serviceName, funcName) {
		// get token
		bearerToken := utils.GetMetaDataField(md, "authorization") // Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6Ik...
		utils.PrintObj(bearerToken, "bearerToken")
		token := utils.ExtractToken(bearerToken)

		// check token
		status, user := tokenUtil.ValidToken(token)
		if status != Common.ErrorCodes_SUCCESS {
			return status, ctx
		}

		md.Append("user_data_json", utils.ToJson(user))

		// check email valid
		if !*user.EmailValid {
			return Common.ErrorCodes_USER_EMAIL_INVAILD, ctx
		}

		// check permission
		if !permissionUtils.CheckUserPermission(user, funcName) {
			return Common.ErrorCodes_PERMISSION_DENIED, ctx
		}

		// get program data
		programJson, status := programUtils.GetProgramDataJson(user)
		if status != Common.ErrorCodes_SUCCESS {
			return status, ctx
		}

		md.Append("program_data_json", programJson)
	}

	newCtx := metadata.NewIncomingContext(ctx, md)

	return Common.ErrorCodes_SUCCESS, newCtx
}
