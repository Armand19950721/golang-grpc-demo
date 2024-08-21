package ArContentService

import (
	"service/model"
	ArContentPb "service/protos/ArContent"
	"service/protos/ArContentTemplate"
	ArContentTemplatePb "service/protos/ArContentTemplate"
	"service/protos/Common"
	"service/utils"
	"service/utils/s3Utils"

	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var (
	ArContentCollet    = getArContentCollet()
	vectorDefault      = Common.Vector{X: "0", Y: "0", Z: "0"}
	vectorDefaultScale = Common.Vector{X: "1", Y: "1", Z: "1"}
	colorWhite         = "#ffffff"
	colorBlack         = "#000000"
)

func getArContentCollet() map[ArContentPb.ArContentCategoryEnum]map[ArContentPb.ArContentTypeEnum][]ArContentPb.ArContentTemplateEnum {
	// Implement
	// map[category] []type
	// map[type] []template

	headTypes := map[ArContentPb.ArContentTypeEnum][]ArContentPb.ArContentTemplateEnum{
		ArContentPb.ArContentTypeEnum_TYPE_EAR: {
			ArContentPb.ArContentTemplateEnum_TEMPLATE_EARRING,
		},
		ArContentPb.ArContentTypeEnum_TYPE_HAIR: {
			ArContentPb.ArContentTemplateEnum_TEMPLATE_HAIR,
		},
	}

	faceTypes := map[ArContentPb.ArContentTypeEnum][]ArContentPb.ArContentTemplateEnum{
		ArContentPb.ArContentTypeEnum_TYPE_GLASSES: {
			ArContentPb.ArContentTemplateEnum_TEMPLATE_GLASSES,
		},
		ArContentPb.ArContentTypeEnum_TYPE_LENS: {
			ArContentPb.ArContentTemplateEnum_TEMPLATE_CONTECT_LENSES,
		},
	}

	collet := map[ArContentPb.ArContentCategoryEnum]map[ArContentPb.ArContentTypeEnum][]ArContentPb.ArContentTemplateEnum{
		ArContentPb.ArContentCategoryEnum_CATRGORY_HEAD: headTypes,
		ArContentPb.ArContentCategoryEnum_CATRGORY_FACE: faceTypes,
	}

	return collet
}

func convertArContentTablesToProtos(rows []model.ArContent) ([]*ArContentPb.ArContentInfo, error) {
	models := []*ArContentPb.ArContentInfo{}
	for _, row := range rows {
		res, err := convertArContentTableToProto(row)
		if err != nil {
			return models, err
		}

		models = append(models, res)
	}

	return models, nil
}

func convertArContentTableToProto(row model.ArContent) (*ArContentPb.ArContentInfo, error) {
	res := &ArContentPb.ArContentInfo{}

	// convert template
	templateSettingByte, err := convertTemplateJsonToProtoByte(row.Template, row.TemplateSetting)
	if err != nil {
		return res, err
	}

	// convert viewer
	// utils.PrintObj(row.ViewerSetting, "ParseJsonWithType")
	viewerSetting, err := utils.ParseJsonWithType[*ArContentPb.ArViewerSetting](row.ViewerSetting)
	if err != nil {
		return res, err
	}

	if viewerSetting.LogoImageName != "" {
		viewerSetting.LogoImagePath = utils.GetImagePath(viewerSetting.LogoImageName, ArContentPb.ArContentImageType_VIEWER_IMAGE)
	}

	res = &ArContentPb.ArContentInfo{
		ArContentId:     row.Id.String(),
		Name:            row.Name,
		Tag:             row.Tag,
		IsOn:            row.IsOn,
		CategoryEnum:    ArContentPb.ArContentCategoryEnum(row.Category),
		TypeEnum:        ArContentPb.ArContentTypeEnum(row.Type),
		TemplateEnum:    ArContentPb.ArContentTemplateEnum(row.Template),
		CreateTime:      utils.ParseDateToString(row.CreateAt),
		UpdateTime:      utils.ParseDateToString(row.UpdateAt),
		TemplateSetting: templateSettingByte,
		ViewerSetting:   viewerSetting,
		ThumbnailPath:   utils.GetImagePath(row.ThumbnailName, ArContentPb.ArContentImageType_THUMBNAIL),
		ThumbnailName:   row.ThumbnailName,
	}

	return res, err
}

func convertArContentTableToProtoForArLink(row model.ArContent) (*ArContentPb.ArContentInfo, error) {
	res := &ArContentPb.ArContentInfo{}

	templateSettingByte, err := convertTemplateJsonToProtoByte(row.Template, row.TemplateSetting)
	if err != nil {
		return res, err
	}

	viewerSetting, err := utils.ParseJsonWithType[*ArContentPb.ArViewerSetting](row.ViewerSetting)
	if err != nil {
		return res, err
	}

	if viewerSetting.LogoImageName != "" {
		viewerSetting.LogoImagePath = utils.GetImagePath(viewerSetting.LogoImageName, ArContentPb.ArContentImageType_VIEWER_IMAGE)
	}

	res = &ArContentPb.ArContentInfo{
		ArContentId:     row.Id.String(),
		IsOn:            row.IsOn,
		CategoryEnum:    ArContentPb.ArContentCategoryEnum(row.Category),
		TypeEnum:        ArContentPb.ArContentTypeEnum(row.Type),
		TemplateEnum:    ArContentPb.ArContentTemplateEnum(row.Template),
		TemplateSetting: templateSettingByte,
		ViewerSetting:   viewerSetting,
	}

	return res, err
}

func checkTemplatePermission(basicInfo utils.BasicInfo, arContentData model.ArContent) bool {
	//if program is includeing ENUM (all). will pass the check
	pass, err := checkEnumPass(basicInfo, ArContentPb.ArContentCollet_TEMPLATE)
	if err != nil {
		utils.PrintObj(err.Error(), "checkTemplatePermission err")
	}

	utils.PrintObj(pass, "checkTemplatePermission pass")

	if pass {
		return true
	}

	// check program can update this template
	hasTemplatePermission := false
	for _, item := range basicInfo.Program.Templates {
		if item == arContentData.Template {
			hasTemplatePermission = true
		}
	}

	utils.PrintObj(hasTemplatePermission, "checkTemplatePermission")

	if !hasTemplatePermission || !pass {
		return false
	}

	return true
}

func GetObjectDefaultSetting() *ArContentTemplatePb.ObjectSetting {
	return &ArContentTemplatePb.ObjectSetting{
		Position: &vectorDefault,
		Rotation: &vectorDefault,
		Scale:    &vectorDefaultScale,
	}
}

func convertTemplateJsonToProtoByte(template ArContentPb.ArContentTemplateEnum, templateSettingJson string) ([]byte, error) {
	switch template {
	case ArContentPb.ArContentTemplateEnum_TEMPLATE_GLASSES:
		templateProto, err := utils.ParseJsonWithType[ArContentTemplate.TemplateGlasses](templateSettingJson)
		if err != nil {
			return []byte(""), status.Errorf(codes.Internal, err.Error())
		}

		for i := 0; i < len(templateProto.Settings); i++ {
			if templateProto.Settings[i].Setting.ImageName != "" {
				templateProto.Settings[i].Setting.ImagePath = utils.GetImagePath(templateProto.Settings[i].Setting.ImageName, ArContentPb.ArContentImageType_TEMPLATE_IMAGE)
			}
		}

		return proto.Marshal(&templateProto)
	case ArContentPb.ArContentTemplateEnum_TEMPLATE_EARRING:
		templateProto, err := utils.ParseJsonWithType[ArContentTemplate.TemplateEarring](templateSettingJson)
		if err != nil {
			return []byte(""), status.Errorf(codes.Internal, err.Error())
		}

		for i := 0; i < len(templateProto.Left.Settings); i++ {
			if templateProto.Left.Settings[i].Setting.ImageName != "" {
				templateProto.Left.Settings[i].Setting.ImagePath = utils.GetImagePath(templateProto.Left.Settings[i].Setting.ImageName, ArContentPb.ArContentImageType_TEMPLATE_IMAGE)
			}
		}

		for i := 0; i < len(templateProto.Right.Settings); i++ {
			if templateProto.Right.Settings[i].Setting.ImageName != "" {
				templateProto.Right.Settings[i].Setting.ImagePath = utils.GetImagePath(templateProto.Right.Settings[i].Setting.ImageName, ArContentPb.ArContentImageType_TEMPLATE_IMAGE)
			}
		}

		return proto.Marshal(&templateProto)
	case ArContentPb.ArContentTemplateEnum_TEMPLATE_CONTECT_LENSES:
		templateProto, err := utils.ParseJsonWithType[ArContentTemplate.TemplateContactLenses](templateSettingJson)
		if err != nil {
			return []byte(""), status.Errorf(codes.Internal, err.Error())
		}

		if templateProto.ImageName != "" {
			templateProto.ImagePath = utils.GetImagePath(templateProto.ImageName, ArContentPb.ArContentImageType_TEMPLATE_IMAGE)
		}

		return proto.Marshal(&templateProto)
	case ArContentPb.ArContentTemplateEnum_TEMPLATE_HAIR:
		templateProto, err := utils.ParseJsonWithType[*ArContentTemplatePb.TemplateHair](templateSettingJson)
		if err != nil {
			return []byte(""), status.Errorf(codes.Internal, err.Error())
		}

		return proto.Marshal(templateProto)
	}

	return []byte(""), status.Errorf(codes.FailedPrecondition, "there has no match template for convert template setting json")
}

func copyTemplateImageProcess(oldFilename string) (err error, newFileName string) {

	if oldFilename != "" {
		newFileName = uuid.New().String() + utils.GetSuffixByFileName(oldFilename)

		err := s3Utils.CopyFileWithBucket(oldFilename, ArContentPb.ArContentImageType_TEMPLATE_IMAGE, newFileName, ArContentPb.ArContentImageType_TEMPLATE_IMAGE)
		if err != nil {
			return err, newFileName
		}
	}

	return nil, newFileName
}

func copyAndUpdateTemplateImage(templateSettingJson string, templateEnum ArContentPb.ArContentTemplateEnum) (err error, newTemplateJson string) {

	switch templateEnum {
	case ArContentPb.ArContentTemplateEnum_TEMPLATE_EARRING:
		templateProto, err := utils.ParseJsonWithType[ArContentTemplate.TemplateEarring](templateSettingJson)
		if err != nil {
			return err, ""
		}

		for _, item := range templateProto.Left.Settings {

			filename := item.Setting.ImageName
			err, newFileName := copyTemplateImageProcess(filename)
			if err != nil {
				return err, ""
			}

			item.Setting.ImageName = newFileName
			newTemplateJson = utils.ToJson(&templateProto)
		}

		for _, item := range templateProto.Right.Settings {

			filename := item.Setting.ImageName
			err, newFileName := copyTemplateImageProcess(filename)
			if err != nil {
				return err, ""
			}

			item.Setting.ImageName = newFileName
			newTemplateJson = utils.ToJson(&templateProto)
		}

	case ArContentPb.ArContentTemplateEnum_TEMPLATE_GLASSES:
		templateProto, err := utils.ParseJsonWithType[ArContentTemplate.TemplateGlasses](templateSettingJson)
		if err != nil {
			return err, ""
		}

		for _, item := range templateProto.Settings {

			filename := item.Setting.ImageName
			err, newFileName := copyTemplateImageProcess(filename)
			if err != nil {
				return err, ""
			}

			item.Setting.ImageName = newFileName
			newTemplateJson = utils.ToJson(&templateProto)
		}

	case ArContentPb.ArContentTemplateEnum_TEMPLATE_CONTECT_LENSES:
		templateProto, err := utils.ParseJsonWithType[ArContentTemplate.TemplateContactLenses](templateSettingJson)
		if err != nil {
			return err, ""
		}

		filename := templateProto.ImageName
		err, newFileName := copyTemplateImageProcess(filename)
		if err != nil {
			return err, ""
		}

		templateProto.ImageName = newFileName
		newTemplateJson = utils.ToJson(&templateProto)
	default:
		utils.PrintObj(templateEnum, "templtae dont need to process")
		newTemplateJson = templateSettingJson
		return nil, newTemplateJson
	}

	return nil, newTemplateJson
}

func deleteTemplateImageProcess(filename string) error {
	if filename != "" {
		// delete
		err := s3Utils.DeleteFile(filename, ArContentPb.ArContentImageType_TEMPLATE_IMAGE)
		if err != nil {
			utils.PrintObj(err.Error(), "deleteTemplateImage")
		}
	}
	return nil
}

func deleteTemplateImage(templateSettingJson string, templateEnum ArContentPb.ArContentTemplateEnum) error {

	switch templateEnum {
	case ArContentPb.ArContentTemplateEnum_TEMPLATE_EARRING:
		templateProto, err := utils.ParseJsonWithType[ArContentTemplate.TemplateEarring](templateSettingJson)
		if err != nil {
			return err
		}

		for _, item := range templateProto.Left.Settings {
			filename := item.Setting.ImageName

			err := deleteTemplateImageProcess(filename)
			if err != nil {
				utils.PrintObj(err.Error(), "deleteTemplateImage")
			}
		}

		for _, item := range templateProto.Right.Settings {
			filename := item.Setting.ImageName

			err := deleteTemplateImageProcess(filename)
			if err != nil {
				utils.PrintObj(err.Error(), "deleteTemplateImage")
			}
		}

	case ArContentPb.ArContentTemplateEnum_TEMPLATE_GLASSES:
		templateProto, err := utils.ParseJsonWithType[ArContentTemplate.TemplateGlasses](templateSettingJson)
		if err != nil {
			return err
		}

		for _, item := range templateProto.Settings {
			filename := item.Setting.ImageName

			err := deleteTemplateImageProcess(filename)
			if err != nil {
				utils.PrintObj(err.Error(), "deleteTemplateImage")
			}
		}

	case ArContentPb.ArContentTemplateEnum_TEMPLATE_CONTECT_LENSES:
		templateProto, err := utils.ParseJsonWithType[ArContentTemplate.TemplateContactLenses](templateSettingJson)
		if err != nil {
			return err
		}

		filename := templateProto.ImageName

		err = deleteTemplateImageProcess(filename)
		if err != nil {
			utils.PrintObj(err.Error(), "deleteTemplateImage")
		}

	default:
		utils.PrintObj(templateEnum, "templtae dont need to process")
		return nil
	}

	return nil
}

func updateEarringSide(side *ArContentTemplatePb.EarringSettingSide) (*ArContentTemplatePb.EarringSettingSide, error) {
	for _, item := range side.Settings {

		utils.PrintObj([]string{
			item.Setting.ImageUploadName,
			item.Setting.ImageName,
			item.Setting.ImagePath,
		}, "ImageUploadName,ImageName,ImagePath")

		// check exist
		imageUploadName := item.Setting.ImageUploadName

		if imageUploadName != "" {
			// upload
			err := s3Utils.MoveFileWithBucket(imageUploadName, ArContentPb.ArContentImageType_TEMP, ArContentPb.ArContentImageType_TEMPLATE_IMAGE)
			if err != nil {
				utils.PrintObj(err.Error(), "MoveFileWithBucket")
				// return side, err
			}

			// delete old img
			err = deleteTemplateImageProcess(item.Setting.ImageName)
			if err != nil {
				utils.PrintObj(err.Error(), "deleteTemplateImageProcess")
				// return side, err
			}

			// set field
			item.Setting.ImageName = imageUploadName
			item.Setting.ImageUploadName = ""
		}
	}

	return side, nil
}

func updateByTemplateByte(templateByte []byte, templateEnum ArContentPb.ArContentTemplateEnum) (string, error) {
	result := ""

	switch templateEnum {
	case ArContentPb.ArContentTemplateEnum_TEMPLATE_EARRING:
		// upload new image from local temp
		templateProtoNew := &ArContentTemplate.TemplateEarring{}
		err := proto.Unmarshal(templateByte, templateProtoNew)
		if err != nil {
			return "", err
		}

		//update side
		templateProtoNew.Left, err = updateEarringSide(templateProtoNew.Left)
		if err != nil {
			return "", err
		}

		utils.PrintObj(templateProtoNew.Left, "templateProtoNew.Left")

		templateProtoNew.Right, err = updateEarringSide(templateProtoNew.Right)
		if err != nil {
			return "", err
		}

		utils.PrintObj(templateProtoNew.Right, "templateProtoNew.Right")

		result = utils.ToJson(templateProtoNew)

	case ArContentPb.ArContentTemplateEnum_TEMPLATE_GLASSES:
		// upload new image from local temp
		templateProtoNew := &ArContentTemplate.TemplateGlasses{}
		err := proto.Unmarshal(templateByte, templateProtoNew)
		if err != nil {
			return "", err
		}

		for _, item := range templateProtoNew.Settings {
			utils.PrintObj([]string{
				item.Setting.ImageUploadName,
				item.Setting.ImageName,
				item.Setting.ImagePath,
			}, "ImageUploadName,ImageName,ImagePath")

			// check exist
			imageUploadName := item.Setting.ImageUploadName

			if imageUploadName != "" {
				// upload
				err := s3Utils.MoveFileWithBucket(imageUploadName, ArContentPb.ArContentImageType_TEMP, ArContentPb.ArContentImageType_TEMPLATE_IMAGE)
				if err != nil {
					utils.PrintObj(err.Error(), "updateTemplateByTemplateByte")
					return "", err
				}

				// delete old img
				err = deleteTemplateImageProcess(item.Setting.ImageName)
				if err != nil {
					utils.PrintObj(err.Error(), "deleteTemplateImageProcess")
					return "", err
				}

				// set field
				item.Setting.ImageName = item.Setting.ImageUploadName
				item.Setting.ImageUploadName = ""
			}
		}

		result = utils.ToJson(templateProtoNew)

	case ArContentPb.ArContentTemplateEnum_TEMPLATE_HAIR:
		templateProto := &ArContentTemplate.TemplateHair{}
		err := proto.Unmarshal(templateByte, templateProto)
		if err != nil {
			return "", err
		}

		result = utils.ToJson(templateProto)

	case ArContentPb.ArContentTemplateEnum_TEMPLATE_CONTECT_LENSES:
		templateProtoNew := &ArContentTemplate.TemplateContactLenses{}
		err := proto.Unmarshal(templateByte, templateProtoNew)
		if err != nil {
			return "", err
		}

		imageUploadName := *templateProtoNew.ImageUploadName

		if imageUploadName != "" {
			// upload
			err := s3Utils.MoveFileWithBucket(imageUploadName, ArContentPb.ArContentImageType_TEMP, ArContentPb.ArContentImageType_TEMPLATE_IMAGE)
			if err != nil {
				utils.PrintObj(err.Error(), "updateTemplateByTemplateByte")
				return "", err
			}
			// set field
			templateProtoNew.ImageName = *templateProtoNew.ImageUploadName
			*templateProtoNew.ImageUploadName = ""
		}

		result = utils.ToJson(templateProtoNew)

	default:
		utils.PrintObj(templateEnum, "setDefaultContentData err")
		return result, status.Errorf(codes.NotFound, "cant found match template")
	}

	return result, status.Errorf(codes.OK, "")
}

func getEarringSide() *ArContentTemplatePb.EarringSettingSide {
	return &ArContentTemplatePb.EarringSettingSide{
		OverAllPosition: &vectorDefault,
		OverAllRotation: &vectorDefault,
		OverAllScale:    &vectorDefaultScale,
		Settings: []*ArContentTemplatePb.EarringSetting{
			{Location: ArContentTemplatePb.EarringObjectLocationEnum_EARRING_FIRST, Setting: GetObjectDefaultSetting()},
			{Location: ArContentTemplatePb.EarringObjectLocationEnum_EARRING_SECOND, Setting: GetObjectDefaultSetting()},
			{Location: ArContentTemplatePb.EarringObjectLocationEnum_EARRING_THIRD, Setting: GetObjectDefaultSetting()},
		},
	}
}

func getDefaultTransform() *Common.Transform {
	return &Common.Transform{
		Position: &vectorDefault,
		Rotation: &vectorDefault,
		Scale:    &vectorDefaultScale,
	}
}

func getDefaultTemplate(templateEnum ArContentPb.ArContentTemplateEnum) (string, error) {
	result := ""

	switch templateEnum {
	case ArContentPb.ArContentTemplateEnum_TEMPLATE_EARRING:
		result = utils.ToJson(ArContentTemplatePb.TemplateEarring{
			Left:  getEarringSide(),
			Right: getEarringSide(),
		})
	case ArContentPb.ArContentTemplateEnum_TEMPLATE_GLASSES:
		result = utils.ToJson(ArContentTemplatePb.TemplateGlasses{
			OverAllPosition: &vectorDefault,
			OverAllRotation: &vectorDefault,
			OverAllScale:    &vectorDefaultScale,
			Settings: []*ArContentTemplatePb.GlassesSetting{
				{Location: ArContentTemplatePb.GlassesObjectLocationEnum_GLASSES_LEFT, Setting: GetObjectDefaultSetting()},
				{Location: ArContentTemplatePb.GlassesObjectLocationEnum_GLASSES_FRONT, Setting: GetObjectDefaultSetting()},
				{Location: ArContentTemplatePb.GlassesObjectLocationEnum_GLASSES_RIGHT, Setting: GetObjectDefaultSetting()},
			},
		})
	case ArContentPb.ArContentTemplateEnum_TEMPLATE_HAIR:
		result = utils.ToJson(ArContentTemplatePb.TemplateHair{
			ColorCode:     "#000000",
			AlphaFeather:  "0",
			AlphaSoftness: "0",
		})
	case ArContentPb.ArContentTemplateEnum_TEMPLATE_CONTECT_LENSES:
		result = utils.ToJson(ArContentTemplatePb.TemplateContactLenses{
			Left:      getDefaultTransform(),
			Right:     getDefaultTransform(),
			ColorCode: colorWhite,
			ImageName: "",
			Radius:    "1",
		})
	default:
		utils.PrintObj(templateEnum, "setDefaultContentData err")
		return result, status.Errorf(codes.Internal, "cant found match template")
	}

	return result, status.Errorf(codes.OK, "")
}

func getDefaultViewerSetting() *ArContentPb.ArViewerSetting {
	defaultBool := true

	viewerSetting := &ArContentPb.ArViewerSetting{
		LogoImagePath:       "",
		AllowCapture:        &defaultBool,
		BorderColor:         "#a0defd",
		CameraButtonColor:   "#ef0606",
		UploadLogoImageName: "",
		ViewerRightButton:   getDefaultLinkButton("right"),
		ViewerLeftButton:    getDefaultLinkButton("left"),
	}

	return viewerSetting
}

func getDefaultLinkButton(cmd string) *ArContentPb.LinkButton {
	defaultBool := true

	model := &ArContentPb.LinkButton{
		Enable:    &defaultBool,
		BtnColor:  colorWhite,
		TextColor: colorBlack,
	}

	if cmd == "left" {
		model.BtnText = "Left button"
	}

	if cmd == "right" {
		model.BtnText = "Right button"
	}

	return model
}

func checkArContentCollet(
	basicInfo utils.BasicInfo,
	categoryEnum ArContentPb.ArContentCategoryEnum,
	typeEnum ArContentPb.ArContentTypeEnum,
	templateEnum ArContentPb.ArContentTemplateEnum) bool {

	categoryPass := false
	typePass := false
	templatePass := false

	categoryPass, err := checkEnumPass(basicInfo, ArContentPb.ArContentCollet_CATEGORY)
	if err != nil {
		return false
	}

	if !categoryPass {
		for _, item := range basicInfo.Program.Categories {
			if item == categoryEnum {
				categoryPass = true
			}
		}
	}

	typePass, err = checkEnumPass(basicInfo, ArContentPb.ArContentCollet_TYPE)
	if err != nil {
		return false
	}

	if !typePass {
		for _, item := range basicInfo.Program.Types {
			if item == typeEnum {
				typePass = true
			}
		}
	}

	templatePass, err = checkEnumPass(basicInfo, ArContentPb.ArContentCollet_TEMPLATE)
	if err != nil {
		return false
	}

	if !templatePass {
		for _, item := range basicInfo.Program.Templates {
			if item == templateEnum {
				templatePass = true
			}
		}
	}

	utils.PrintObj([]bool{categoryPass, typePass, templatePass}, "checkArContentCollet")

	if !categoryPass || !typePass || !templatePass {
		return false
	}

	return true
}

func checkEnumPass(basicInfo utils.BasicInfo, collet ArContentPb.ArContentCollet) (pass bool, err error) {
	// check program is include the enum
	switch collet {
	case ArContentPb.ArContentCollet_CATEGORY:
		for _, categoryProgram := range basicInfo.Program.Categories {
			if categoryProgram == ArContentPb.ArContentCategoryEnum_CATRGORY_ALL {
				pass = true
				break
			}
		}
	case ArContentPb.ArContentCollet_TYPE:
		for _, typeProgram := range basicInfo.Program.Types {
			if typeProgram == ArContentPb.ArContentTypeEnum_TYPE_ALL {
				pass = true
				break
			}
		}
	case ArContentPb.ArContentCollet_TEMPLATE:
		for _, templateProgram := range basicInfo.Program.Templates {
			if templateProgram == ArContentPb.ArContentTemplateEnum_TEMPLATE_ALL {
				pass = true
				break
			}
		}
	default:
		return false, errors.New("collet not support")
	}

	return pass, nil
}

func getStringZeroByNumber(num int) string {
	str := ""

	for i := 0; i < num; i++ {
		str += "0"
	}

	return str
}
