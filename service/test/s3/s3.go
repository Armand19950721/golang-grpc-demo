package s3

import (
	"os"
	"service/protos/ArContent"
	"service/test/testUtils"
	"service/utils"
	"service/utils/s3Utils"
)

func Flow() {
	filename := "thumbnail-1.jpg"
	imageTypes := []ArContent.ArContentImageType{ArContent.ArContentImageType_TEMP, ArContent.ArContentImageType_THUMBNAIL, ArContent.ArContentImageType_TEMPLATE_IMAGE, ArContent.ArContentImageType_VIEWER_IMAGE}

	// get file
	utils.PrintTitle("open file")
	file, err := os.Open("./static/" + filename)
	testUtils.DisplayResult("", err, false)

	for _, imageType := range imageTypes {
		utils.PrintTitle("UploadFile")
		err = s3Utils.UploadImageFile(file, filename, imageType)
		testUtils.DisplayResult("", err, false)

		utils.PrintTitle("FileExists")
		exist, err := s3Utils.FileExists(filename, imageType)
		testUtils.DisplayResult(exist, err, false)

		utils.PrintTitle("DeleteFile")
		err = s3Utils.DeleteFile(filename, imageType)
		testUtils.DisplayResult("", err, false)
	}

	// copy
	utils.PrintTitle("UploadFile test copy")
	err = s3Utils.UploadImageFile(file, filename, ArContent.ArContentImageType_TEMP)
	testUtils.DisplayResult("", err, false)

	utils.PrintTitle("FileExists test copy")
	exist, err := s3Utils.FileExists(filename, ArContent.ArContentImageType_TEMP)
	testUtils.DisplayResult(exist, err, false)

	newFileName := filename + "_new"

	utils.PrintTitle("CopyFileWithBucket test copy")
	err = s3Utils.CopyFileWithBucket(filename, ArContent.ArContentImageType_TEMP, newFileName, ArContent.ArContentImageType_THUMBNAIL)
	testUtils.DisplayResult(exist, err, false)

	utils.PrintTitle("FileExists test copy")
	exist, err = s3Utils.FileExists(newFileName, ArContent.ArContentImageType_THUMBNAIL)
	testUtils.DisplayResult(exist, err, false)

	utils.PrintTitle("DeleteFile test copy")
	err = s3Utils.DeleteFile(newFileName, ArContent.ArContentImageType_TEMP)
	utils.PrintObj(err)

	utils.PrintTitle("DeleteFile test copy")
	err = s3Utils.DeleteFile(newFileName, ArContent.ArContentImageType_THUMBNAIL)
	utils.PrintObj(err)

	// copy fail
	utils.PrintTitle("UploadFile test copy fail")
	notExistFileName := "notExistFileName"
	err = s3Utils.CopyFileWithBucket(notExistFileName, ArContent.ArContentImageType_TEMP, notExistFileName, ArContent.ArContentImageType_THUMBNAIL)
	testUtils.DisplayResult("", err, false) // 20230324 ignore the file sourse of copy function.

	// check email image
	createEmailLogoIfNotExist()
}

func createEmailLogoIfNotExist() bool {
	utils.PrintTitle("createEmailLogoIfNotExist")

	// make sure the image of email useage
	logoName := "MetARlogo.png"
	exist, err := s3Utils.FileExists(logoName, ArContent.ArContentImageType_STATIC_IMAGE)
	if err != nil {
		panic(err.Error())
	}

	if !exist {
		utils.PrintObj("cant found email logo. start create")

		file, err := os.Open("../utils/EmailUtils/static/" + logoName)
		if err != nil {
			panic(err.Error())
		}

		err = s3Utils.UploadImageFile(file, logoName, ArContent.ArContentImageType_STATIC_IMAGE)
		if err != nil {
			panic(err.Error())
		}

	} else {
		utils.PrintObj("found email logo. containue ")
	}

	return true
}
