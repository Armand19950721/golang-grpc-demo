package server

import (
	"log"
	"net/http"
	"service/controller"
	"service/protos/Common"
	"service/utils"
	tokenUtil "service/utils/token"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitGinServer() {
	r := setRouter()
	err := r.Run(utils.GetEnv("GIN_SERVER_IP"))
	if err != nil {
		log.Fatalf("upload server run fail")
		return
	}
}

func setRouter() *gin.Engine {
	r := gin.Default()

	corsConf := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"},
		AllowHeaders: []string{"authorization", "Authorization", "Content-Type", "Upgrade", "Origin",
			"Connection", "Accept-Encoding", "Accept-Language", "Host", "Access-Control-Request-Method",
			"Access-Control-Request-Headers"},
	}

	r.Use(cors.New(corsConf))
	r.Use(middleware)
	r.POST("/api/upload_file", controller.UploadImage)
	r.GET("/api/static"+utils.GetEnv("AWS_S3_PATH_TEMP")+":imageName", controller.GetTempImage)
	r.GET("/api/static"+utils.GetEnv("AWS_S3_PATH_THUMBNAIL")+":imageName", controller.GetThumbnailImage)
	r.GET("/api/static"+utils.GetEnv("AWS_S3_PATH_TEMPLATE_IMAGE")+":imageName", controller.GetTemplateImage)
	r.GET("/api/static"+utils.GetEnv("AWS_S3_PATH_VIEWER_IMAGE")+":imageName", controller.GetViewerImage)
	r.GET("/api/static"+utils.GetEnv("AWS_S3_PATH_STATIC_IMAGE")+":imageName", controller.GetStaticImage)

	return r
}

func noNeedToken(uri string) bool {
	utils.PrintObj(uri, "uri")

	arr := strings.Split(uri, "/")
	utils.PrintObj(arr, "arr")

	if len(arr) >= 2 {
		secondRoute := arr[1]

		if secondRoute == "api" {
			thirdRoute := arr[2]

			if strings.Contains(thirdRoute, "static") {
				return true
			}

			if strings.Contains(thirdRoute, "DeleteJames") {
				return true
			}

			if strings.Contains(thirdRoute, "echo") {
				return true
			}
		}
	}

	return false
}

func middleware(ctx *gin.Context) {
	// check route need token
	if noNeedToken(ctx.Request.RequestURI) {
		utils.PrintObj(true, "noNeedToken")

		ctx.Next()
		return
	}
	utils.PrintObj(false, "noNeedToken")

	// get token
	bearerToken := ctx.Request.Header.Get("authorization")

	utils.PrintObj(bearerToken, "bearerToken")

	token := utils.ExtractToken(bearerToken)

	// check token
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, utils.GetErrorGin(utils.ErrorType{
			Code: Common.ErrorCodes_INVAILD_TOKEN,
		}))
		ctx.Abort()
		return
	}

	if state, _ := tokenUtil.ValidToken(token); state != Common.ErrorCodes_SUCCESS {
		ctx.JSON(http.StatusUnauthorized, utils.GetErrorGin(utils.ErrorType{
			Code: Common.ErrorCodes_INVAILD_TOKEN,
		}))
		ctx.Abort()
		return
	}

	ctx.Next()
}
