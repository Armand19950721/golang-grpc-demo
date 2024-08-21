package controller

import (
	"errors"
	"net/http"
	"service/protos/ArContent"
	"service/protos/Common"
	"service/utils"
	"service/utils/s3Utils"

	"github.com/gin-gonic/gin"
)

func GetTempImage(ctx *gin.Context) {
	ctx = downloadFile(ctx, ArContent.ArContentImageType_TEMP)
}

func GetThumbnailImage(ctx *gin.Context) {
	ctx = downloadFile(ctx, ArContent.ArContentImageType_THUMBNAIL)
}

func GetTemplateImage(ctx *gin.Context) {
	ctx = downloadFile(ctx, ArContent.ArContentImageType_TEMPLATE_IMAGE)
}

func GetViewerImage(ctx *gin.Context) {
	ctx = downloadFile(ctx, ArContent.ArContentImageType_VIEWER_IMAGE)
}

func GetStaticImage(ctx *gin.Context) {
	ctx = downloadFile(ctx, ArContent.ArContentImageType_STATIC_IMAGE)
}

func UploadImage(ctx *gin.Context) {
	uploadRes, uploadedFileName := uploadFile(ctx, 1, Common.UploadFileType_IMAGE)

	if uploadRes.Code != Common.ErrorCodes_SUCCESS {
		ctx.JSON(http.StatusOK, uploadRes)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"FileName": uploadedFileName,
		"FilePath": utils.GetDomainAPI() + utils.GetEnv("GIN_IMAGE_ROUTE") + utils.GetEnv("AWS_S3_PATH_TEMP") + uploadedFileName,
		"Code":     Common.ErrorCodes_SUCCESS,
		"Message":  "ok",
	})
}

func Upload3DModel(ctx *gin.Context) {
	uploadRes, uploadedFileName := uploadFile(ctx, 1, Common.UploadFileType_THREE_D_FILE)

	if uploadRes.Code != Common.ErrorCodes_SUCCESS {
		ctx.JSON(http.StatusOK, uploadRes)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"FileName": uploadedFileName,
		"FilePath": utils.GetDomainAPI() + utils.GetEnv("GIN_IMAGE_ROUTE") + utils.GetEnv("AWS_S3_PATH_TEMP") + uploadedFileName,
		"Code":     Common.ErrorCodes_SUCCESS,
		"Message":  "ok",
	})
}

func downloadFile(ctx *gin.Context, type_ ArContent.ArContentImageType) *gin.Context {
	fileName := ctx.Param("imageName")
	folderPath, err := utils.GetFolderPath(type_)

	if err != nil {
		ctx.JSON(http.StatusOK, utils.GetErrorGin(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: err.Error(),
		}))
		return ctx
	}

	byteArr, err := s3Utils.DownloadFile(folderPath + fileName)

	if err != nil {
		ctx.JSON(http.StatusOK, utils.GetErrorGin(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: err.Error(),
		}))
		return ctx
	}

	ctx.Data(0, "file", byteArr)
	return ctx
}

func uploadFile(ctx *gin.Context, sizeMB int64, _type Common.UploadFileType) (result utils.ErrorType, uploadedFileName string) {
	// check size
	var limitSizeMB int64 = 1
	var maxBytes int64 = 1024 * 1024 * limitSizeMB
	var writer http.ResponseWriter = ctx.Writer

	ctx.Request.Body = http.MaxBytesReader(writer, ctx.Request.Body, maxBytes)

	if err := ctx.Request.ParseMultipartForm(maxBytes); err != nil {
		utils.PrintObj(err.Error(), "check size")
		result = utils.ErrorType{
			Code:      Common.ErrorCodes_UPLOAD_FILE_SIZE_INVALID,
			ReturnMsg: "max " + utils.ToString(int(limitSizeMB)) + " MB",
		}
		return
	}

	// read file
	file, err := ctx.FormFile("file") //get file
	if err != nil {
		result = utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: err.Error(),
		}
		return
	}

	// get file type
	fileSuffix := utils.GetSuffix(file)

	// check file type
	suffixWhiteList := map[string]string{}

	switch _type {
	case Common.UploadFileType_IMAGE:
		suffixWhiteList[".jpg"] = ""
		suffixWhiteList[".jpeg"] = ""
		suffixWhiteList[".png"] = ""
	case Common.UploadFileType_THREE_D_FILE:
		suffixWhiteList[".glb"] = ""
		suffixWhiteList[".gltf"] = ""
	}

	utils.PrintObj(fileSuffix, "fileSuffix")

	if _, isExist := suffixWhiteList[fileSuffix]; !isExist {
		result = utils.ErrorType{
			Code: Common.ErrorCodes_UPLOAD_FILE_TYPE_NOT_SUPPORT,
		}
		return
	}

	// get reader
	reader, err := file.Open()
	if err != nil {
		result = utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: err.Error(),
		}
		return
	}

	// get folder path
	path := ""
	getPathErr := errors.New("")

	switch _type {
	case Common.UploadFileType_IMAGE:
		path, getPathErr = utils.GetFolderPath(ArContent.ArContentImageType_TEMP)
	case Common.UploadFileType_THREE_D_FILE:
		path, getPathErr = utils.GetFolderPath(ArContent.ArContentImageType_TEMP) // also upload to s3 temp folder
	}

	if getPathErr != nil {
		result = utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: getPathErr.Error(),
		}
		return
	}

	// s3 upload
	uploadFileName := utils.GetNewFileName(file)
	err = s3Utils.UploadFile(reader, uploadFileName, path)
	if err != nil {
		result = utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: err.Error(),
		}
		return
	}

	result = utils.ErrorType{
		Code: Common.ErrorCodes_SUCCESS,
	}
	uploadedFileName = uploadFileName

	return
}
