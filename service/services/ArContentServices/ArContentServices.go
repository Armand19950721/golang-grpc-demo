package ArContentService

import (
	"context"
	"service/model"
	ArContentPb "service/protos/ArContent"
	"service/protos/Common"
	"service/repositories"
	"service/utils"
	"service/utils/s3Utils"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetCategoryList(ctx context.Context, req *ArContentPb.GetCategoryListRequest) (*ArContentPb.GetCategoryListReply, error) {
	reply := &ArContentPb.GetCategoryListReply{}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	//if program collet is includeing ENUM (all). will pass the program check
	pass, err := checkEnumPass(basicInfo, ArContentPb.ArContentCollet_CATEGORY)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	utils.PrintObj(pass, "pass")

	// get all category
	categoriesEnum := []ArContentPb.ArContentCategoryEnum{}

	for categoryEnum := range ArContentCollet {
		// filter by program
		for _, categoryProgramEnum := range basicInfo.Program.Categories {
			// if match then append
			if categoryEnum == categoryProgramEnum || pass {
				categoriesEnum = append(categoriesEnum, categoryEnum)
			}
		}
	}

	reply.CategoryEnum = categoriesEnum

	return reply, status.Errorf(codes.OK, "")
}

func GetTypeList(ctx context.Context, req *ArContentPb.GetTypeListRequest) (*ArContentPb.GetTypeListReply, error) {
	reply := &ArContentPb.GetTypeListReply{}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	//if program collet is includeing ENUM (all). will pass the program check
	pass, err := checkEnumPass(basicInfo, ArContentPb.ArContentCollet_TYPE)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	// get match type
	typesEnum := []ArContentPb.ArContentTypeEnum{}

	for type_ := range ArContentCollet[req.CategoryEnum] {
		// filter by program
		for _, typeProgramEnum := range basicInfo.Program.Types {
			// if match then append
			if type_ == typeProgramEnum || pass {
				typesEnum = append(typesEnum, type_)
			}
		}
	}

	reply.TypeEnum = typesEnum

	return reply, status.Errorf(codes.OK, "")
}

func GetTemplateList(ctx context.Context, req *ArContentPb.GetTemplateListRequest) (*ArContentPb.GetTemplateListReply, error) {
	reply := &ArContentPb.GetTemplateListReply{}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	//if program collet is includeing ENUM (all). will pass the program check
	pass, err := checkEnumPass(basicInfo, ArContentPb.ArContentCollet_TEMPLATE)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	// get match template
	templatesTemp := []ArContentPb.ArContentTemplateEnum{}

	for typeEnum := range ArContentCollet[req.CategoryEnum] { // find []type
		if typeEnum == req.TypeEnum {
			foundTemplateArray, exist := ArContentCollet[req.CategoryEnum][typeEnum] // find template[]

			if exist {
				templatesTemp = foundTemplateArray
			} else {
				return reply, status.Errorf(codes.Internal, "fail to get template")
			}
		}
	}

	// filter
	templateEnums := []ArContentPb.ArContentTemplateEnum{}

	if pass {
		// no need program filter if pass
		templateEnums = templatesTemp
	} else {
		// filter by program
		for _, templateEnum := range templatesTemp {
			for _, templateProgramEnum := range basicInfo.Program.Templates {
				// if match then append
				if templateEnum == templateProgramEnum {
					templateEnums = append(templateEnums, ArContentPb.ArContentTemplateEnum(templateEnum))
				}
			}
		}
	}

	reply.TemplateEnum = templateEnums

	return reply, status.Errorf(codes.OK, "")
}

func CreateArContent(ctx context.Context, req *ArContentPb.CreateArContentRequest) (*ArContentPb.CreateArContentReply, error) {
	reply := &ArContentPb.CreateArContentReply{}

	// check param
	if !utils.ValidString(req.Name, 1, 100) ||
		!utils.ValidString(req.Tag, 1, 100, "nullable") ||
		!utils.ValidString(req.UploadThumbnailName, 1, -1) {

		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
	}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	// // get ar content for check repeat // 20221006 ray說先不擋重複名稱
	// res, count := repositories.GetArContentCount(model.ArContent{
	// 	Name:   req.Name,
	// 	UserId: basicInfo.AdminId,
	// })

	// if res.Error != nil && !utils.IsErrorNotFound(res.Error) {
	// 	return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
	// 		Code:        Common.ErrorCodes_UNKNOWN_ERROR,
	// 		InternalMsg: res.Error.Error(),
	// 	}))
	// }

	// if count != 0 {
	// 	return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
	// 		Code: Common.ErrorCodes_REPEATED_ERROR,
	// 	}))
	// }

	// check collet permission
	pass := checkArContentCollet(basicInfo, req.CategoryEnum, req.TypeEnum, req.TemplateEnum)
	if !pass {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code: Common.ErrorCodes_PROGRAM_NOT_SUPPORT,
		}))
	}

	// move file from temp to thumbnail folder in s3
	err := s3Utils.MoveFileWithBucket(req.UploadThumbnailName, ArContentPb.ArContentImageType_TEMP, ArContentPb.ArContentImageType_THUMBNAIL)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	// get default for template setting
	templateSetting, err := getDefaultTemplate(req.TemplateEnum)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	// create ar content
	isOnDefault := true
	createArContent := model.ArContent{
		Name:            req.Name,
		Tag:             req.Tag,
		ThumbnailName:   req.UploadThumbnailName,
		Category:        req.CategoryEnum,
		Type:            req.TypeEnum,
		Template:        req.TemplateEnum,
		TemplateSetting: templateSetting,
		ViewerSetting:   utils.ToJson(getDefaultViewerSetting()),
		UserId:          basicInfo.AdminId,
		ViewerUrlId:     utils.GetNullableString(""),
		IsOn:            &isOnDefault,
	}
	res, row := repositories.CreateArContent(createArContent)

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// convert
	convertRes, err := convertArContentTableToProto(row)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	reply.Data = convertRes

	return reply, status.Errorf(codes.OK, "")
}

func UpdateArContentTemplate(ctx context.Context, req *ArContentPb.UpdateArContentTemplateRequest) (*ArContentPb.UpdateArContentTemplateReply, error) {
	reply := &ArContentPb.UpdateArContentTemplateReply{}

	// check param
	if !utils.ValidId(req.ArContentId) {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
	}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	// get ar content
	arContentQuery := model.ArContent{
		Id:     utils.ParseUUID(req.ArContentId),
		UserId: basicInfo.AdminId,
	}

	result, arContentData := repositories.GetArContent(arContentQuery)
	if result.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: result.Error.Error(),
		}))
	}

	// check program can update this template
	if !checkTemplatePermission(basicInfo, arContentData) {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code: Common.ErrorCodes_PROGRAM_NOT_SUPPORT,
		}))
	}

	// convert incoming template data
	newTemplateJson, err := updateByTemplateByte(req.TemplateData, arContentData.Template)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	// update ar content
	arContentData.TemplateSetting = newTemplateJson

	result, updateData := repositories.UpdateArContentTemplate(model.ArContent{Id: arContentData.Id}, arContentData)

	if result.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: result.Error.Error(),
		}))
	}

	// convert
	arContentProto, err := convertArContentTableToProto(updateData)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	reply.Data = arContentProto
	return reply, status.Errorf(codes.OK, "")
}

func UpdateArContentViewer(ctx context.Context, req *ArContentPb.UpdateArContentViewerRequest) (*ArContentPb.UpdateArContentViewerReplay, error) {
	reply := &ArContentPb.UpdateArContentViewerReplay{}

	// check param
	if !utils.ValidId(req.ArContentId) {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
	}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	// get ar content
	arContentQuery := model.ArContent{
		Id:     utils.ParseUUID(req.ArContentId),
		UserId: basicInfo.AdminId,
	}

	result, arContentData := repositories.GetArContent(arContentQuery)
	if result.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: result.Error.Error(),
		}))
	}

	// check program can update this template
	if !checkTemplatePermission(basicInfo, arContentData) {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code: Common.ErrorCodes_PROGRAM_NOT_SUPPORT,
		}))
	}

	// move image
	if req.ViewerSetting.UploadLogoImageName != "" {
		err := s3Utils.MoveFileWithBucket(req.ViewerSetting.UploadLogoImageName, ArContentPb.ArContentImageType_TEMP, ArContentPb.ArContentImageType_VIEWER_IMAGE)
		if err != nil {
			return reply, utils.ReturnUnKnownError(err)
		}

		req.ViewerSetting.LogoImageName = req.ViewerSetting.UploadLogoImageName
		req.ViewerSetting.UploadLogoImageName = ""
	}

	// update ar content
	arContentData.ViewerSetting = utils.ToJson(req.ViewerSetting)

	result, updateData := repositories.UpdateArContentTemplate(model.ArContent{Id: arContentData.Id}, arContentData)
	if result.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: result.Error.Error(),
		}))
	}

	// convert
	arContentProto, err := convertArContentTableToProto(updateData)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	reply.Data = arContentProto

	return reply, status.Errorf(codes.OK, "")
}

func GetArContent(ctx context.Context, req *ArContentPb.GetArContentRequest) (*ArContentPb.GetArContentReply, error) {
	reply := &ArContentPb.GetArContentReply{}

	// check param
	if !utils.ValidId(req.ArContentId) {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
	}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	// get ar content
	res, data := repositories.GetArContent(model.ArContent{
		Id:     utils.ParseUUID(req.ArContentId),
		UserId: basicInfo.AdminId,
	})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// convert
	model, err := convertArContentTableToProto(data)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	reply.Data = model

	return reply, status.Errorf(codes.OK, "")
}

func UpdateArContent(ctx context.Context, req *ArContentPb.UpdateArContentRequest) (*ArContentPb.UpdateArContentReply, error) {
	reply := &ArContentPb.UpdateArContentReply{}

	// check param
	if !utils.ValidId(req.ArContentId) ||
		!utils.ValidString(req.Name, 1, 100) ||
		!utils.ValidString(req.Tag, 1, 100, "nullable") ||
		!utils.ValidString(*req.UploadThumbnailName, 1, -1, "nullable") {

		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))

	}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	// get ar content
	res, data := repositories.GetArContent(model.ArContent{
		Id:     utils.ParseUUID(req.ArContentId),
		UserId: basicInfo.AdminId,
	})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// update thumbnail
	isUpdateThumbnail := utils.ValidString(*req.UploadThumbnailName, 1, -1)
	oldThumbnail := data.ThumbnailName

	if isUpdateThumbnail {
		data.ThumbnailName = *req.UploadThumbnailName

		// move file from temp to thumbnail folder in s3
		err := s3Utils.MoveFileWithBucket(*req.UploadThumbnailName, ArContentPb.ArContentImageType_TEMP, ArContentPb.ArContentImageType_THUMBNAIL)
		if err != nil {
			return reply, utils.ReturnUnKnownError(res.Error)
		}
	}

	// update
	data.Name = req.Name
	data.Tag = req.Tag

	res = repositories.UpdateArContent(model.ArContent{Id: utils.ParseUUID(req.ArContentId), UserId: basicInfo.AdminId}, data)
	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// delete old thumbnail
	if isUpdateThumbnail {
		err := s3Utils.DeleteFile(oldThumbnail, ArContentPb.ArContentImageType_THUMBNAIL)
		if err != nil {
			return reply, utils.ReturnUnKnownError(res.Error)
		}
	}

	// convert
	model, err := convertArContentTableToProto(data)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	reply.Data = model

	return reply, status.Errorf(codes.OK, "")
}

func GetArContentList(ctx context.Context, req *ArContentPb.GetArContentListRequest) (*ArContentPb.GetArContentListReply, error) {
	reply := &ArContentPb.GetArContentListReply{}

	// check param
	if !utils.ValidString(req.Keyword, 1, 100, "nullable") ||
		req.PageInfo == nil {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))

	}

	req.PageInfo = utils.ValidPageInfo(req.PageInfo)

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	// query
	queryResult, arContents, count := repositories.GetArContentList(basicInfo, req)
	if queryResult.Error != nil {
		return reply, status.Errorf(codes.Internal, queryResult.Error.Error())
	}

	// convert
	convertRes, err := convertArContentTablesToProtos(arContents)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	reply.Data = convertRes
	reply.PageInfo = &Common.PageInfoReply{TotalCount: count}

	return reply, status.Errorf(codes.OK, "")
}

func UpdateArContentIsOn(ctx context.Context, req *ArContentPb.UpdateArContentIsOnRequest) (*Common.CommonReply, error) {
	reply := &Common.CommonReply{}

	// check param
	if !utils.ValidId(req.ArContentId) {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))

	}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	// get ar content
	res, data := repositories.GetArContent(model.ArContent{
		Id:     utils.ParseUUID(req.ArContentId),
		UserId: basicInfo.AdminId,
	})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// update
	data.IsOn = req.IsOn

	res = repositories.UpdateArContent(model.ArContent{
		Id:     utils.ParseUUID(req.ArContentId),
		UserId: basicInfo.AdminId,
	}, data)

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	return reply, status.Errorf(codes.OK, "")
}

func DeleteArContent(ctx context.Context, req *ArContentPb.DeleteArContentRequest) (*Common.CommonReply, error) {
	reply := &Common.CommonReply{}

	// check param
	if !utils.ValidId(req.ArContentId) {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))

	}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	// get ar content
	res, data := repositories.GetArContent(model.ArContent{
		Id:     utils.ParseUUID(req.ArContentId),
		UserId: basicInfo.AdminId,
	})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// delete ar content
	res, _ = repositories.DeleteArContent(model.ArContent{
		Id:     utils.ParseUUID(req.ArContentId),
		UserId: basicInfo.AdminId,
	}, data)

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// delete thumbnail image
	err := s3Utils.DeleteFile(data.ThumbnailName, ArContentPb.ArContentImageType_THUMBNAIL)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	// delete template image
	err = deleteTemplateImage(data.TemplateSetting, ArContentPb.ArContentTemplateEnum(data.Template))
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	// delete viewer image
	viewerSetting, err := utils.ParseJsonWithType[*ArContentPb.ArViewerSetting](data.ViewerSetting)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	err = s3Utils.DeleteFile(viewerSetting.LogoImageName, ArContentPb.ArContentImageType_VIEWER_IMAGE)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	return reply, status.Errorf(codes.OK, "")
}

func DuplicateArContent(ctx context.Context, req *ArContentPb.DuplicateArContentRequest) (*ArContentPb.DuplicateArContentReply, error) {
	reply := &ArContentPb.DuplicateArContentReply{}

	// check param
	if !utils.ValidId(req.ArContentId) {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))

	}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	// get ar content
	res, data := repositories.GetArContent(model.ArContent{
		Id:     utils.ParseUUID(req.ArContentId),
		UserId: basicInfo.AdminId,
	})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// copy thumbnail image
	if data.ThumbnailName != "" {
		newFileName := uuid.New().String() + utils.GetSuffixByFileName(data.ThumbnailName)

		err := s3Utils.CopyFileWithBucket(data.ThumbnailName, ArContentPb.ArContentImageType_THUMBNAIL, newFileName, ArContentPb.ArContentImageType_THUMBNAIL)
		if err != nil {
			return reply, status.Errorf(codes.Internal, err.Error()+"(copy thumbnail image)")
		}

		data.ThumbnailName = newFileName
	}

	// copy template image
	err, newTemplateSetting := copyAndUpdateTemplateImage(data.TemplateSetting, ArContentPb.ArContentTemplateEnum(data.Template))
	if err != nil {
		return reply, status.Errorf(codes.Internal, err.Error()+"(copy template image)")
	}
	utils.PrintObj(newTemplateSetting, "newTemplateSetting")

	data.TemplateSetting = newTemplateSetting

	// copy viewer image
	viewerSetting, err := utils.ParseJsonWithType[*ArContentPb.ArViewerSetting](data.ViewerSetting)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	if viewerSetting.LogoImageName != "" {
		newFileName := uuid.New().String() + utils.GetSuffixByFileName(viewerSetting.LogoImageName)

		err := s3Utils.CopyFileWithBucket(viewerSetting.LogoImageName, ArContentPb.ArContentImageType_VIEWER_IMAGE, newFileName, ArContentPb.ArContentImageType_VIEWER_IMAGE)
		if err != nil {
			return reply, status.Errorf(codes.Internal, err.Error()+"(copy viewer image)"+viewerSetting.LogoImageName)
		}
		viewerSetting.LogoImageName = newFileName
		data.ViewerSetting = utils.ToJson(viewerSetting)
	}

	// give new info
	data.Id = uuid.New()
	data.Name = utils.AutoRename(data.Name)
	data.ViewerUrlId = utils.GetNullableString("")
	data.CreateAt = time.Now().UTC()
	data.UpdateAt = time.Now().UTC()

	// create
	res, newArContent := repositories.CreateArContent(data)
	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// convert
	model, err := convertArContentTableToProto(newArContent)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	reply.Data = model

	return reply, status.Errorf(codes.OK, "")
}

func GetArLink(ctx context.Context, req *ArContentPb.GetArLinkRequest) (*ArContentPb.GetArLinkReply, error) {
	reply := &ArContentPb.GetArLinkReply{}

	// check param
	if !utils.ValidId(req.ArContentId) {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))

	}

	// get basic data
	basicInfo := utils.GetBasicInfo(ctx)
	if !basicInfo.Success {
		return reply, utils.ReturnBasicInfoError()
	}

	// find ar content
	res, data := repositories.GetArContent(model.ArContent{
		UserId: basicInfo.AdminId,
		Id:     utils.ParseUUID(req.ArContentId),
	})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// upadat if "ViewerUrlId" is null
	if !data.ViewerUrlId.Valid {

		data.ViewerUrlId = utils.GetNullableString(utils.ToString(utils.GetRandomInt(100000000000000, 999999999999999)))

		res = repositories.UpdateArContent(model.ArContent{
			Id:     utils.ParseUUID(req.ArContentId),
			UserId: basicInfo.AdminId,
		}, data)

		if res.Error != nil {
			return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
				Code:        Common.ErrorCodes_UNKNOWN_ERROR,
				InternalMsg: res.Error.Error(),
			}))
		}
	}

	reply.ArContentViewerPath = utils.GetDomain() + utils.GetEnv("AR_CONTENT_VIEWER_LINK_PATH") + data.ViewerUrlId.String

	return reply, status.Errorf(codes.OK, "")
}

func GetViewerData(ctx context.Context, req *ArContentPb.GetViewerDataRequest) (*ArContentPb.GetViewerDataReply, error) {
	reply := &ArContentPb.GetViewerDataReply{}

	// get ar content
	result, data := repositories.GetArContent(model.ArContent{
		ViewerUrlId: utils.GetNullableString(req.ArContentViewerUrlId),
	})

	if result.Error != nil {
		if utils.IsErrorNotFound(result.Error) {
			return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{
				Code: Common.ErrorCodes_DATA_NOT_FOUND,
			}))
		} else {
			return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
				Code:        Common.ErrorCodes_UNKNOWN_ERROR,
				InternalMsg: result.Error.Error(),
			}))
		}
	}

	// convert
	convertRes, err := convertArContentTableToProtoForArLink(data)
	if err != nil {
		return reply, utils.ReturnUnKnownError(err)
	}

	reply.AdminId = data.UserId.String()
	reply.Data = convertRes

	return reply, status.Errorf(codes.OK, "")
}

func UpdateArContentThreeDModel(ctx context.Context, req *ArContentPb.UpdateArContentThreeDModelRequest) (*Common.CommonReply, error) {
	reply := &Common.CommonReply{}

	// valid
	if !utils.ValidId(req.ArContentId) ||
		req.UploadedThreeDModelFilename == "" {
		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
	}

	utils.PrintObj([]string{req.ArContentId, req.UploadedThreeDModelFilename}, "UpdateArContentThreeDModel")

	// check ar content
	result, _ := repositories.GetArContent(model.ArContent{
		Id: utils.ParseUUID(req.ArContentId),
	})

	if result.Error != nil {
		if utils.IsErrorNotFound(result.Error) {
			return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{
				Code: Common.ErrorCodes_DATA_NOT_FOUND,
			}))
		} else {
			return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
				Code:        Common.ErrorCodes_UNKNOWN_ERROR,
				InternalMsg: result.Error.Error(),
			}))
		}
	}

	return reply, status.Errorf(codes.OK, "")
}
