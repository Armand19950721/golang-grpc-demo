package thirtyPartyInterceptor

import (
	"context"
	"strings"

	"service/protos/Common"
	"service/protos/ThirdPartyCommon"
	"service/utils"
	programUtils "service/utils/program"
	tokenUtil "service/utils/token"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

func Valid(md metadata.MD, info *grpc.UnaryServerInfo, ctx context.Context) (ThirdPartyCommon.StatusCode, context.Context) {
	funcName, _, _ := getRouteInfo(info)

	if funcName == "Verification" {
		return ThirdPartyCommon.StatusCode_SUCCESS, ctx
	}

	// get token
	bearerToken := utils.GetMetaDataField(md, "authorization") // Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6Ik...

	utils.PrintObj(bearerToken, "bearerToken")

	token := utils.ExtractToken(bearerToken)

	// check token and get user
	status, user := tokenUtil.ValidToken(token)
	if status != Common.ErrorCodes_SUCCESS {
		return ThirdPartyCommon.StatusCode_INVAILD_TOKEN, ctx
	}

	md.Append("user_data_json", utils.ToJson(user))

	// get program data
	programJson, status := programUtils.GetProgramDataJson(user)
	if status != Common.ErrorCodes_SUCCESS {
		return ThirdPartyCommon.StatusCode_INTERNAL_ERROR, ctx
	}

	md.Append("program_data_json", programJson)

	// renew ctx
	newCtx := metadata.NewIncomingContext(ctx, md)

	return ThirdPartyCommon.StatusCode_SUCCESS, newCtx
}
