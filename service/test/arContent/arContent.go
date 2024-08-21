package arContent

import (
	"errors"
	"strings"
	"time"

	"service/test/testUtils"
	"service/utils"
	"service/utils/s3Utils"

	"service/protos/ArContent"
	"service/protos/ArContentTemplate"
	"service/protos/Common"
	"service/protos/WebServices"

	"google.golang.org/protobuf/proto"
)

var (
	c              = WebServices.NewArContentServicesClient(testUtils.Conn)
	arContentId    = ""
	arContentNames = []string{}
	arContentTags  = []string{}
)

func CreateMultiArContent() {
	ticker := time.NewTicker(1 * time.Second)
	count := 0

	for range ticker.C {

		CreateFlow(ArContent.ArContentTemplateEnum_TEMPLATE_GLASSES, false)

		if count == 10 {
			break
		}

		count++
	}

}

func TestListOrder() {
	ctx := testUtils.GetCtx()
	///
	r, err := c.GetArContentList(ctx, &ArContent.GetArContentListRequest{
		PageInfo: &Common.PageInfoRequest{
			CurrentPageNum: 1,
			PageItemCount:  10,
		},
	})

	testUtils.DisplayResult(r, err, false)
}

func getRedomVactor() *Common.Vector {
	min := 1
	max := 255

	return &Common.Vector{
		X: utils.ToString(testUtils.RandomInt(min, max)),
		Y: utils.ToString(testUtils.RandomInt(min, max)),
		Z: utils.ToString(testUtils.RandomInt(min, max)),
	}
}

type GinReply struct {
	FileName string
	FilePath string
	Code     Common.ErrorCodes
	Message  string
}

func UploadImageExpress() string {
	uploadRes := testUtils.UploadFile("thumbnail-1.jpg")
	if uploadRes == "" {
		testUtils.DisplayResult(uploadRes, errors.New("upload fail"), false)
	}
	uploadResFormatOk, err := utils.ParseJsonWithType[GinReply](uploadRes)

	utils.PrintObj(&uploadResFormatOk, "uploadResFormatOk")

	if err != nil {
		utils.PrintObj("uploadImage err", "uploadImage")
		return ""
	}

	if uploadResFormatOk.Code != Common.ErrorCodes_SUCCESS {
		testUtils.DisplayResult(&uploadResFormatOk, errors.New(uploadResFormatOk.Message), false)
	} else {
		testUtils.DisplayResult(&uploadResFormatOk, nil, false)
	}

	return uploadResFormatOk.FileName
}

func getRedomObjectSetting() *ArContentTemplate.ObjectSetting {
	utils.PrintTitle("UploadFile")
	newImageName := UploadImageExpress()

	return &ArContentTemplate.ObjectSetting{
		Position:        getRedomVactor(),
		Rotation:        getRedomVactor(),
		Scale:           getRedomVactor(),
		ImageUploadName: newImageName,
	}
}

func CreateFlow(template ArContent.ArContentTemplateEnum, expactError bool) {
	ctx := testUtils.GetCtx()

	utils.PrintTitle("GetCategoryList")

	r, err := c.GetCategoryList(ctx, &ArContent.GetCategoryListRequest{})

	testUtils.DisplayResult(r, err, false)

	utils.PrintTitle("GetTypeList")

	// set collet
	category := ArContent.ArContentCategoryEnum_CATRGORY_HEAD
	type_ := ArContent.ArContentTypeEnum_TYPE_EAR

	r1, err := c.GetTypeList(ctx, &ArContent.GetTypeListRequest{
		CategoryEnum: category,
	})

	testUtils.DisplayResult(r1, err, false)

	utils.PrintTitle("GetTemplateList")

	r2, err := c.GetTemplateList(ctx, &ArContent.GetTemplateListRequest{
		CategoryEnum: category,
		TypeEnum:     type_,
	})

	testUtils.DisplayResult(r2, err, false)

	utils.PrintTitle("UploadFile ok")

	uploadNewFileName := UploadImageExpress()

	utils.PrintObj(uploadNewFileName, "uploadNewFileName")

	utils.PrintTitle("check File upload exist ok")
	exist, err := s3Utils.FileExists(uploadNewFileName, ArContent.ArContentImageType_TEMP)
	testUtils.DisplayResult(exist, err, false)

	utils.PrintTitle("UploadFile suffix err")

	uploadRes := testUtils.UploadFile("thumbnail-1.gif")

	uploadResFormat, err := utils.ParseJsonWithType[Common.ErrorReplyGin](uploadRes)
	if err != nil {
		testUtils.DisplayResult("", err, false)
	}

	if uploadResFormat.Code != Common.ErrorCodes_SUCCESS {
		testUtils.DisplayResult(&uploadResFormat, errors.New(uploadRes), true)
	} else {
		testUtils.DisplayResult(&uploadResFormat, nil, true)
	}

	utils.PrintTitle("UploadFile size err")

	uploadRes = testUtils.UploadFile("bigfile.jpg")

	uploadResFormat, err = utils.ParseJsonWithType[Common.ErrorReplyGin](uploadRes)
	if err != nil {
		testUtils.DisplayResult("", err, false)
	}

	if uploadResFormat.Code != Common.ErrorCodes_SUCCESS {
		testUtils.DisplayResult(&uploadResFormat, errors.New(uploadRes), true)
	} else {
		testUtils.DisplayResult(&uploadResFormat, nil, true)
	}

	utils.PrintTitle("CreateArContent expactError:" + utils.ParseBoolToString(expactError))

	arName := "ar_name_" + testUtils.RandomIntString()
	arTag := "ar_tag_" + testUtils.RandomIntString() + "," + "ar_tag_" + testUtils.RandomIntString()

	r3, err := c.CreateArContent(ctx, &ArContent.CreateArContentRequest{
		Name:                arName,
		Tag:                 arTag,
		UploadThumbnailName: uploadNewFileName,
		CategoryEnum:        category,
		TypeEnum:            type_,
		TemplateEnum:        template,
	})

	testUtils.DisplayResult(r3, err, expactError)

	if expactError {
		return
	}

	basicInfo := testUtils.GetBasicInfo()
	basicInfo.ArContentId = r3.Data.ArContentId
	testUtils.SetBasicInfo(basicInfo)

	if len(r3.Data.TemplateSetting) == 0 {
		testUtils.DisplayResult(r3, errors.New("TemplateSetting fail"), false)
	}

	if utils.ToJson(r3.Data.ViewerSetting) == "{}" {
		testUtils.DisplayResult(r3, errors.New("ViewerSetting fail"), false)
	}

	if len(r3.Data.ThumbnailPath) == 0 {
		testUtils.DisplayResult(r3, errors.New("ThumbnailPath fail"), false)
	}

	arContentId = r3.Data.ArContentId
	arContentNames = append(arContentNames, r3.Data.Name)
	arContentTags = append(arContentTags, r3.Data.Tag)

	utils.PrintTitle("check File upload exist ok")
	exist, err = s3Utils.FileExists(uploadNewFileName, ArContent.ArContentImageType_THUMBNAIL)
	testUtils.DisplayResult(exist, err, false)

	// utils.PrintTitle("CreateArContent repeat")  // 20221006 ray說先不擋重複名稱

	// r3, err = c.CreateArContent(ctx, &ArContent.CreateArContentRequest{
	// 	Name:                arName,
	// 	Tag:                 arTag,
	// 	UploadThumbnailName: uploadNewFileName,
	// 	CategoryEnum:        category,
	// 	TypeEnum:            type_,
	// 	TemplateEnum:        template,
	// })

	// testUtils.DisplayResult(r3, err, true)
}

func UpdateViewSetting() {
	ctx := testUtils.GetCtx()

	utils.PrintTitle("updateViewSetting")
	newImageName := UploadImageExpress()
	boolVal := true
	viewSetting := ArContent.ArViewerSetting{
		UploadLogoImageName: newImageName,
		BorderColor:         "#FF12345",
		CameraButtonColor:   "#FF12345",
		AllowCapture:        &boolVal,
	}

	r, err := c.UpdateArContentViewer(ctx, &ArContent.UpdateArContentViewerRequest{
		ArContentId:   arContentId,
		ViewerSetting: &viewSetting,
	})

	testUtils.DisplayResult(r, err, false)

	utils.PrintTitle("updateViewSetting check")

	checkOk := true

	if *viewSetting.AllowCapture != *r.Data.ViewerSetting.AllowCapture {
		utils.PrintObj([]bool{*viewSetting.AllowCapture, *r.Data.ViewerSetting.AllowCapture})
		checkOk = markFail("AllowCapture")
	}

	if viewSetting.BorderColor != r.Data.ViewerSetting.BorderColor {
		checkOk = markFail("BorderColor")
	}

	if viewSetting.CameraButtonColor != r.Data.ViewerSetting.CameraButtonColor {
		checkOk = markFail("ButtonColor")
	}

	if utils.GetImagePath(newImageName, ArContent.ArContentImageType_VIEWER_IMAGE) != r.Data.ViewerSetting.LogoImagePath {
		checkOk = markFail("LogoImagePath")
	}

	if !checkOk {
		testUtils.DisplayResult("", errors.New("updateViewSetting fail"), false)
	} else {
		testUtils.DisplayResult("check success", nil, false)
	}
}

func TemplateGlassesUpdateFlow() {
	ctx := testUtils.GetCtx()

	utils.PrintTitle("upload template image x3")

	obj1 := getRedomObjectSetting()
	obj2 := getRedomObjectSetting()
	obj3 := getRedomObjectSetting()

	utils.PrintTitle("templateGlassesUpdateFlow")

	template := &ArContentTemplate.TemplateGlasses{
		OverAllPosition: getRedomVactor(),
		OverAllRotation: getRedomVactor(),
		OverAllScale:    getRedomVactor(),
		Settings: []*ArContentTemplate.GlassesSetting{
			{Location: ArContentTemplate.GlassesObjectLocationEnum_GLASSES_LEFT, Setting: obj1},
			{Location: ArContentTemplate.GlassesObjectLocationEnum_GLASSES_FRONT, Setting: obj2},
			{Location: ArContentTemplate.GlassesObjectLocationEnum_GLASSES_RIGHT, Setting: obj3},
		},
	}

	utils.PrintObj(template, "TemplateGlassesUpdateFlow template")

	templateBytes, _ := proto.Marshal(template)

	r5, err := c.UpdateArContentTemplate(ctx, &ArContent.UpdateArContentTemplateRequest{
		ArContentId:  arContentId,
		TemplateData: templateBytes,
	})

	testUtils.DisplayResult(r5, err, false)

	utils.PrintTitle("templateGlassesUpdateFlow check")

	templateRes := ArContentTemplate.TemplateGlasses{}
	err = proto.Unmarshal(r5.Data.TemplateSetting, &templateRes)
	if err != nil {
		testUtils.DisplayResult("", err, false)
	}
	checkOk := true
	utils.PrintObj(&templateRes, "templateRes")

	if utils.ToJson(templateRes.OverAllPosition) != utils.ToJson(template.OverAllPosition) {
		checkOk = markFail("OverAllPosition")
	}
	if utils.ToJson(templateRes.OverAllRotation) != utils.ToJson(template.OverAllRotation) {
		checkOk = markFail("OverAllRotation")
	}
	if utils.ToJson(templateRes.OverAllScale) != utils.ToJson(template.OverAllScale) {
		checkOk = markFail("OverAllScale")
	}

	for _, item := range templateRes.Settings {
		if item.Location == ArContentTemplate.GlassesObjectLocationEnum_GLASSES_LEFT {
			compareSettingGlasses(item, obj1, true)
		}
		if item.Location == ArContentTemplate.GlassesObjectLocationEnum_GLASSES_FRONT {
			compareSettingGlasses(item, obj2, true)
		}
		if item.Location == ArContentTemplate.GlassesObjectLocationEnum_GLASSES_RIGHT {
			compareSettingGlasses(item, obj3, true)
		}
	}

	if !checkOk {
		testUtils.DisplayResult("", errors.New("UpdateArContentTemplate fail"), false)
	} else {
		testUtils.DisplayResult("check success", nil, false)
	}

	utils.PrintTitle("GetArContentRequest")

	///
	r6, err := c.GetArContent(ctx, &ArContent.GetArContentRequest{
		ArContentId: arContentId,
	})

	testUtils.DisplayResult(r6, err, false)

	utils.PrintTitle("GetArContentRequest check")

	err = proto.Unmarshal(r6.Data.TemplateSetting, &templateRes)
	if err != nil {
		testUtils.DisplayResult("", err, false)
	}

	utils.PrintObj(&templateRes, "templateRes")

	///
	utils.PrintTitle("TemplateGlassesUpdateFlow template remove upload name")

	for i := 0; i < len(templateRes.Settings); i++ {
		templateRes.Settings[i].Setting.ImageUploadName = ""
	}

	utils.PrintObj(&templateRes, "templateRes")

	templateBytes, _ = proto.Marshal(&templateRes)

	r5, err = c.UpdateArContentTemplate(ctx, &ArContent.UpdateArContentTemplateRequest{
		ArContentId:  arContentId,
		TemplateData: templateBytes,
	})

	testUtils.DisplayResult(r5, err, false)

	templateBytes, _ = proto.Marshal(template)

	utils.PrintTitle("templateGlassesUpdateFlow check")

	templateResFormat := ArContentTemplate.TemplateGlasses{}
	err = proto.Unmarshal(r5.Data.TemplateSetting, &templateResFormat)
	if err != nil {
		testUtils.DisplayResult("", err, false)
	}

	utils.PrintObj(&templateResFormat, "templateResFormat")

	for _, item := range templateResFormat.Settings {
		if item.Location == ArContentTemplate.GlassesObjectLocationEnum_GLASSES_LEFT {
			compareSettingGlasses(item, obj1, false)
		}
		if item.Location == ArContentTemplate.GlassesObjectLocationEnum_GLASSES_FRONT {
			compareSettingGlasses(item, obj2, false)
		}
		if item.Location == ArContentTemplate.GlassesObjectLocationEnum_GLASSES_RIGHT {
			compareSettingGlasses(item, obj3, false)
		}
	}

	if !checkOk {
		testUtils.DisplayResult("", errors.New("UpdateArContentTemplate fail"), false)
	} else {
		testUtils.DisplayResult("check success", nil, false)
	}

	///
	r6, err = c.GetArContent(ctx, &ArContent.GetArContentRequest{
		ArContentId: arContentId,
	})

	testUtils.DisplayResult(r6, err, false)

	utils.PrintTitle("GetArContentRequest check again")

	err = proto.Unmarshal(r6.Data.TemplateSetting, &templateRes)
	if err != nil {
		testUtils.DisplayResult("", err, false)
	}

	utils.PrintObj(&templateRes, "templateRes GetArContentRequest check again")
}

func compareSettingGlasses(item *ArContentTemplate.GlassesSetting, obj *ArContentTemplate.ObjectSetting, hasUploadImage bool) bool {
	utils.PrintObj(obj, "req")
	utils.PrintObj(item.Setting, "reply")
	checkOk := true

	if item.Setting.ImageName == "" {
		checkOk = markFail("ImageName")
	}
	if item.Setting.ImagePath != utils.GetImagePath(obj.ImageUploadName, ArContent.ArContentImageType_TEMPLATE_IMAGE) {
		checkOk = markFail("ImagePath")
	}
	if utils.ToJson(item.Setting.Scale) != utils.ToJson(obj.Scale) {
		checkOk = markFail("Scale")
	}
	if utils.ToJson(item.Setting.Position) != utils.ToJson(obj.Position) {
		checkOk = markFail("Position")
	}
	if utils.ToJson(item.Setting.Rotation) != utils.ToJson(obj.Rotation) {
		checkOk = markFail("Rotation")
	}

	if !checkOk {
		testUtils.DisplayResult("", errors.New("compareSettingGlasses test fail"), false)
		return false
	}

	if !checkTemplateImageExist(item.Setting.ImageName) {
		testUtils.DisplayResult("", errors.New("cant found image fail:"+item.Setting.ImageName), false)
		return false
	}

	if hasUploadImage {
		if item.Setting.ImageUploadName == "" {
			checkOk = markFail("ImageName")
		}
	}

	return checkOk
}

func TemplateEarringUpdateFlow() {
	ctx := testUtils.GetCtx()

	utils.PrintTitle("upload template image x3")

	obj1 := getRedomObjectSetting()
	obj2 := getRedomObjectSetting()
	obj3 := getRedomObjectSetting()

	utils.PrintTitle("templateEarringUpdateFlow")

	template := &ArContentTemplate.TemplateEarring{
		Left: &ArContentTemplate.EarringSettingSide{
			OverAllPosition: getRedomVactor(),
			OverAllRotation: getRedomVactor(),
			OverAllScale:    getRedomVactor(),
			Settings: []*ArContentTemplate.EarringSetting{
				{Location: ArContentTemplate.EarringObjectLocationEnum_EARRING_FIRST, Setting: obj1},
				{Location: ArContentTemplate.EarringObjectLocationEnum_EARRING_SECOND, Setting: obj2},
				{Location: ArContentTemplate.EarringObjectLocationEnum_EARRING_THIRD, Setting: obj3},
			},
		},
		Right: &ArContentTemplate.EarringSettingSide{
			OverAllPosition: getRedomVactor(),
			OverAllRotation: getRedomVactor(),
			OverAllScale:    getRedomVactor(),
			Settings: []*ArContentTemplate.EarringSetting{
				{Location: ArContentTemplate.EarringObjectLocationEnum_EARRING_FIRST, Setting: obj1},
				{Location: ArContentTemplate.EarringObjectLocationEnum_EARRING_SECOND, Setting: obj2},
				{Location: ArContentTemplate.EarringObjectLocationEnum_EARRING_THIRD, Setting: obj3},
			},
		},
	}

	templateBytes, _ := proto.Marshal(template)

	r, err := c.UpdateArContentTemplate(ctx, &ArContent.UpdateArContentTemplateRequest{
		ArContentId:  arContentId,
		TemplateData: templateBytes,
	})

	testUtils.DisplayResult(r, err, false)

	utils.PrintTitle("templateEarringUpdateFlow check")

	templateRes := ArContentTemplate.TemplateEarring{}
	err = proto.Unmarshal(r.Data.TemplateSetting, &templateRes)
	if err != nil {
		testUtils.DisplayResult("", err, false)
	}
	checkOk := true
	utils.PrintObj(&templateRes, "templateRes")

	left := templateRes.Left
	leftOrigin := template.Left

	right := templateRes.Right
	rightOrigin := template.Right

	///
	if utils.ToJson(left.OverAllPosition) != utils.ToJson(leftOrigin.OverAllPosition) {
		checkOk = markFail("OverAllPosition")
	}
	if utils.ToJson(left.OverAllRotation) != utils.ToJson(leftOrigin.OverAllRotation) {
		checkOk = markFail("OverAllRotation")
	}
	if utils.ToJson(left.OverAllScale) != utils.ToJson(leftOrigin.OverAllScale) {
		checkOk = markFail("OverAllScale")
	}

	for _, item := range left.Settings {
		if item.Location == ArContentTemplate.EarringObjectLocationEnum_EARRING_FIRST {
			compareSettingEarring(item, obj1)
		}
		if item.Location == ArContentTemplate.EarringObjectLocationEnum_EARRING_SECOND {
			compareSettingEarring(item, obj2)
		}
		if item.Location == ArContentTemplate.EarringObjectLocationEnum_EARRING_THIRD {
			compareSettingEarring(item, obj3)
		}
	}

	if utils.ToJson(right.OverAllPosition) != utils.ToJson(rightOrigin.OverAllPosition) {
		checkOk = markFail("OverAllPosition")
	}
	if utils.ToJson(right.OverAllRotation) != utils.ToJson(rightOrigin.OverAllRotation) {
		checkOk = markFail("OverAllRotation")
	}
	if utils.ToJson(right.OverAllScale) != utils.ToJson(rightOrigin.OverAllScale) {
		checkOk = markFail("OverAllScale")
	}

	for _, item := range right.Settings {
		if item.Location == ArContentTemplate.EarringObjectLocationEnum_EARRING_FIRST {
			compareSettingEarring(item, obj1)
		}
		if item.Location == ArContentTemplate.EarringObjectLocationEnum_EARRING_SECOND {
			compareSettingEarring(item, obj2)
		}
		if item.Location == ArContentTemplate.EarringObjectLocationEnum_EARRING_THIRD {
			compareSettingEarring(item, obj3)
		}
	}

	///
	if !checkOk {
		testUtils.DisplayResult("", errors.New("UpdateArContentTemplate fail"), false)
	} else {
		testUtils.DisplayResult("check success", nil, false)
	}
}

func compareSettingEarring(item *ArContentTemplate.EarringSetting, obj *ArContentTemplate.ObjectSetting) bool {
	utils.PrintObj(obj, "req")
	utils.PrintObj(item.Setting, "reply")
	checkOk := true

	if item.Setting.ImagePath != utils.GetImagePath(obj.ImageUploadName, ArContent.ArContentImageType_TEMPLATE_IMAGE) {
		checkOk = markFail("ImagePath")
	}
	if utils.ToJson(item.Setting.Scale) != utils.ToJson(obj.Scale) {
		checkOk = markFail("Scale")
	}
	if utils.ToJson(item.Setting.Position) != utils.ToJson(obj.Position) {
		checkOk = markFail("Position")
	}
	if utils.ToJson(item.Setting.Rotation) != utils.ToJson(obj.Rotation) {
		checkOk = markFail("Rotation")
	}

	if !checkOk {
		testUtils.DisplayResult("", errors.New("compareSettingEarring test fail"), false)
		return false
	}

	if !checkTemplateImageExist(item.Setting.ImageName) {
		testUtils.DisplayResult("", errors.New("checkTemplateExist(item.Setting.ImageName) test fail"), false)
		return false
	}

	return checkOk
}

func TemplateContactLensesUpdateFlow() {
	ctx := testUtils.GetCtx()

	utils.PrintTitle("upload template image")

	obj1 := getRedomObjectSetting()

	utils.PrintTitle("templateContactLensesUpdateFlow")

	template := &ArContentTemplate.TemplateContactLenses{
		Right: &Common.Transform{
			Position: getRedomVactor(),
			Rotation: getRedomVactor(),
			Scale:    getRedomVactor(),
		},
		Left: &Common.Transform{
			Position: getRedomVactor(),
			Rotation: getRedomVactor(),
			Scale:    getRedomVactor(),
		},
		ImageUploadName: &obj1.ImageUploadName,
	}

	templateBytes, _ := proto.Marshal(template)

	r5, err := c.UpdateArContentTemplate(ctx, &ArContent.UpdateArContentTemplateRequest{
		ArContentId:  arContentId,
		TemplateData: templateBytes,
	})

	testUtils.DisplayResult(r5, err, false)

	utils.PrintTitle("templateContactLensesUpdateFlow check")

	templateRes := ArContentTemplate.TemplateContactLenses{}
	err = proto.Unmarshal(r5.Data.TemplateSetting, &templateRes)
	if err != nil {
		testUtils.DisplayResult("", err, false)
	}
	checkOk := true
	utils.PrintObj(&templateRes, "templateRes")

	///
	left := templateRes.Left
	leftOrigin := template.Left

	right := templateRes.Right
	rightOrigin := template.Right
	///
	if utils.ToJson(left.Position) != utils.ToJson(leftOrigin.Position) {
		checkOk = markFail("Position")
	}

	if utils.ToJson(left.Rotation) != utils.ToJson(leftOrigin.Rotation) {
		checkOk = markFail("Rotation")
	}

	if utils.ToJson(left.Scale) != utils.ToJson(leftOrigin.Scale) {
		checkOk = markFail("Scale")
	}
	///
	if utils.ToJson(right.Position) != utils.ToJson(rightOrigin.Position) {
		checkOk = markFail("Position")
	}

	if utils.ToJson(right.Rotation) != utils.ToJson(rightOrigin.Rotation) {
		checkOk = markFail("Rotation")
	}

	if utils.ToJson(right.Scale) != utils.ToJson(rightOrigin.Scale) {
		checkOk = markFail("Scale")
	}
	///
	if !checkOk {
		testUtils.DisplayResult("", errors.New("UpdateArContentTemplate fail"), false)
	} else {
		testUtils.DisplayResult("check success", nil, false)
	}

	if !checkTemplateImageExist(templateRes.ImageName) {
		testUtils.DisplayResult("", errors.New("checkTemplateExist(item.Setting.ImageName) test fail"), false)
	}
}

func TemplateHairUpdateFlow() {
	ctx := testUtils.GetCtx()

	utils.PrintTitle("templateHairUpdateFlow")

	template := &ArContentTemplate.TemplateHair{
		ColorCode:     "#FF12345",
		AlphaSoftness: "0.1234",
		AlphaFeather:  "0.1234",
	}

	templateBytes, _ := proto.Marshal(template)

	r5, err := c.UpdateArContentTemplate(ctx, &ArContent.UpdateArContentTemplateRequest{
		ArContentId:  arContentId,
		TemplateData: templateBytes,
	})

	testUtils.DisplayResult(r5, err, false)

	utils.PrintTitle("templateHairUpdateFlow check")

	templateRes := ArContentTemplate.TemplateHair{}
	err = proto.Unmarshal(r5.Data.TemplateSetting, &templateRes)
	if err != nil {
		testUtils.DisplayResult("", err, false)
	}
	checkOk := true
	utils.PrintObj(&templateRes, "templateRes")

	if templateRes.AlphaFeather != template.AlphaFeather {
		checkOk = markFail("AlphaFeather")
	}
	if templateRes.AlphaSoftness != template.AlphaSoftness {
		checkOk = markFail("AlphaSoftness")
	}
	if templateRes.ColorCode != template.ColorCode {
		checkOk = markFail("ColorCode")
	}

	if !checkOk {
		testUtils.DisplayResult("", errors.New("UpdateArContentTemplate fail"), false)
	} else {
		testUtils.DisplayResult("check success", nil, false)
	}
}

func TestCopyAndDeleteArContent() {
	ctx := testUtils.GetCtx()

	CreateFlow(ArContent.ArContentTemplateEnum_TEMPLATE_GLASSES, false)
	TemplateGlassesUpdateFlow()
	UpdateViewSetting()

	utils.PrintTitle("GetArContent for copy")

	r2, err := c.GetArContent(ctx, &ArContent.GetArContentRequest{
		ArContentId: arContentId,
	})

	testUtils.DisplayResult(r2, err, false)

	templateRes := ArContentTemplate.TemplateGlasses{}
	err = proto.Unmarshal(r2.Data.TemplateSetting, &templateRes)
	if err != nil {
		testUtils.DisplayResult("", err, false)
	}

	utils.PrintObj(&templateRes, "GetArContent templateRes")

	utils.PrintTitle("DuplicateArContent")

	r4, err := c.DuplicateArContent(ctx, &ArContent.DuplicateArContentRequest{
		ArContentId: arContentId,
	})

	testUtils.DisplayResult(r4, err, false)

	templateResCopy := ArContentTemplate.TemplateGlasses{}
	err = proto.Unmarshal(r4.Data.TemplateSetting, &templateResCopy)
	if err != nil {
		testUtils.DisplayResult("", err, false)
	}

	utils.PrintObj(&templateResCopy, "DuplicateArContent templateResCopy")

	utils.PrintTitle("DuplicateArContent image name check")

	for idx := range templateResCopy.Settings {
		if templateResCopy.Settings[idx].Setting.ImageName == templateRes.Settings[idx].Setting.ImageName {
			testUtils.DisplayResult("", errors.New("copy before and after the image name is the same"), false)
		} else {
			testUtils.DisplayResult("ok", nil, false)
		}
	}

	utils.PrintTitle("DuplicateArContent image exist check")

	// check new image exist
	for _, item := range templateResCopy.Settings {
		exist, err := s3Utils.FileExists(item.Setting.ImageName, ArContent.ArContentImageType_TEMPLATE_IMAGE)

		if err != nil {
			testUtils.DisplayResult(nil, err, false)
		}

		if !exist {
			testUtils.DisplayResult(exist, errors.New("fail check newFileName exsit:"+item.Setting.ImageName), false)
		} else {
			testUtils.DisplayResult(exist, nil, false)
		}
	}

	utils.PrintTitle("DeleteArContent")

	r5, err := c.DeleteArContent(ctx, &ArContent.DeleteArContentRequest{
		ArContentId: r4.Data.ArContentId,
	})

	testUtils.DisplayResult(r5, err, false)

	templateResDelete := ArContentTemplate.TemplateGlasses{}
	err = proto.Unmarshal(r4.Data.TemplateSetting, &templateResDelete)
	if err != nil {
		testUtils.DisplayResult("", err, false)
	}

	utils.PrintTitle("DeleteArContent image exist check")

	// check new image exist
	for _, item := range templateResDelete.Settings {
		exist, err := s3Utils.FileExists(item.Setting.ImageName, ArContent.ArContentImageType_TEMPLATE_IMAGE)

		if err != nil {
			testUtils.DisplayResult(nil, err, false)
		}

		if !exist {
			testUtils.DisplayResult(exist, errors.New("fail check newFileName exsit:"+item.Setting.ImageName), true)
		} else {
			testUtils.DisplayResult(exist, nil, true)
		}
	}

	utils.PrintTitle("GetArContent for delete")

	r6, err := c.GetArContent(ctx, &ArContent.GetArContentRequest{
		ArContentId: r4.Data.ArContentId,
	})

	testUtils.DisplayResult(r6, err, true)
}

func TestArContentList() {
	ctx := testUtils.GetCtx()

	CreateFlow(ArContent.ArContentTemplateEnum_TEMPLATE_GLASSES, false)
	TemplateGlassesUpdateFlow()
	UpdateViewSetting()

	CreateFlow(ArContent.ArContentTemplateEnum_TEMPLATE_CONTECT_LENSES, false)
	TemplateContactLensesUpdateFlow()
	UpdateViewSetting()

	utils.PrintObj(arContentNames)
	utils.PrintObj(arContentTags)

	utils.PrintTitle("GetArContentList test full name")

	r, err := c.GetArContentList(ctx, &ArContent.GetArContentListRequest{
		Keyword: arContentNames[0],
		PageInfo: &Common.PageInfoRequest{
			CurrentPageNum: 1,
			PageItemCount:  10,
		},
	})

	utils.PrintObj(r.PageInfo, "GetArContentList pageInfo")
	testUtils.DisplayResult(r.Data, err, false)

	utils.PrintTitle("GetArContentList test full name check")

	if r.PageInfo.TotalCount <= 0 || len(r.Data) <= 0 {
		testUtils.DisplayResult("fail", errors.New("test fail"), false)
	} else {
		testUtils.DisplayResult("ok", nil, false)
	}

	utils.PrintTitle("GetArContentList test helf name")

	r, err = c.GetArContentList(ctx, &ArContent.GetArContentListRequest{
		Keyword: strings.Split(arContentNames[0], "_")[2],
		PageInfo: &Common.PageInfoRequest{
			CurrentPageNum: 1,
			PageItemCount:  10,
		},
	})

	utils.PrintObj(r.PageInfo, "GetArContentList pageInfo")
	testUtils.DisplayResult(r.Data, err, false)

	utils.PrintTitle("GetArContentList test helf name check")

	if r.PageInfo.TotalCount <= 0 || len(r.Data) <= 0 {
		testUtils.DisplayResult("fail", errors.New("test fail"), false)
	} else {
		testUtils.DisplayResult("ok", nil, false)
	}

	utils.PrintTitle("GetArContentList test full tag")

	r, err = c.GetArContentList(ctx, &ArContent.GetArContentListRequest{
		Keyword: arContentTags[0],
		PageInfo: &Common.PageInfoRequest{
			CurrentPageNum: 1,
			PageItemCount:  10,
		},
	})

	utils.PrintObj(r.PageInfo, "GetArContentList pageInfo")
	testUtils.DisplayResult(r.Data, err, false)

	utils.PrintTitle("GetArContentList test full tag check")

	if r.PageInfo.TotalCount <= 0 || len(r.Data) <= 0 {
		testUtils.DisplayResult("fail", errors.New("test fail"), false)
	} else {
		testUtils.DisplayResult("ok", nil, false)
	}

	///
	utils.PrintObj(r.PageInfo, "GetArContentList no key word")
	r, err = c.GetArContentList(ctx, &ArContent.GetArContentListRequest{
		PageInfo: &Common.PageInfoRequest{
			CurrentPageNum: 1,
			PageItemCount:  10,
		},
	})

	testUtils.DisplayResult(r.Data, err, false)

	if r.PageInfo.TotalCount <= 0 || len(r.Data) <= 0 {
		testUtils.DisplayResult("fail", errors.New("test fail"), false)
	} else {
		testUtils.DisplayResult("ok", nil, false)
	}
}

func ArContentMainPageFlow() {
	ctx := testUtils.GetCtx()

	TestArContentList()

	utils.PrintTitle("UpdateArContentIsOn")

	boolVal := false

	utils.PrintObj(arContentId, "arContentId")
	r1, err := c.UpdateArContentIsOn(ctx, &ArContent.UpdateArContentIsOnRequest{
		ArContentId: arContentId,
		IsOn:        &boolVal,
	})

	testUtils.DisplayResult(r1, err, false)

	utils.PrintTitle("GetArContent")

	r, err := c.GetArContent(ctx, &ArContent.GetArContentRequest{
		ArContentId: arContentId,
	})

	testUtils.DisplayResult(r, err, false)

	utils.PrintTitle("UpdateArContentIsOn & get check")

	expectVal := false
	if *r.Data.IsOn != expectVal {
		testUtils.DisplayResult(r.Data.IsOn, errors.New("fail check"), false)
	} else {
		testUtils.DisplayResult(r.Data.IsOn, nil, false)
	}

	utils.PrintTitle("UpdateArContent")

	newFileName := UploadImageExpress()

	r3, err := c.UpdateArContent(ctx, &ArContent.UpdateArContentRequest{
		ArContentId:         arContentId,
		Name:                r.Data.Name + "_update",
		Tag:                 r.Data.Tag + "_update",
		UploadThumbnailName: &newFileName,
	})

	testUtils.DisplayResult(r3, err, false)

	utils.PrintTitle("UpdateArContent check")

	if !strings.Contains(r3.Data.Name, "update") {
		testUtils.DisplayResult(r3.Data.Name, errors.New("fail check Name"), false)
	}

	if !strings.Contains(r3.Data.Tag, "update") {
		testUtils.DisplayResult(r3.Data.Tag, errors.New("fail check Tag"), false)
	}

	if r3.Data.ThumbnailName != newFileName {
		testUtils.DisplayResult(r3.Data.ThumbnailName, errors.New("fail check Name:"+newFileName), false)
	}

	// check new image exist
	exist, err := s3Utils.FileExists(newFileName, ArContent.ArContentImageType_THUMBNAIL)

	if err != nil {
		testUtils.DisplayResult(nil, err, false)
	}

	if !exist {
		testUtils.DisplayResult(exist, errors.New("fail check newFileName exsit:"+newFileName), false)
	} else {
		testUtils.DisplayResult(exist, nil, false)
	}

	// check old image not exist
	oldFileName := r.Data.ThumbnailName
	exist, err = s3Utils.FileExists(oldFileName, ArContent.ArContentImageType_THUMBNAIL)

	if err != nil {
		testUtils.DisplayResult(nil, err, false)
	}

	if exist {
		testUtils.DisplayResult(exist, errors.New("fail check oldFileName exsit:"+oldFileName), false)
	} else {
		testUtils.DisplayResult(exist, nil, false)
	}

	utils.PrintTitle("GetArLink create new")

	r6, err := c.GetArLink(ctx, &ArContent.GetArLinkRequest{
		ArContentId: arContentId,
	})

	testUtils.DisplayResult(r6, err, false)

	utils.PrintTitle("GetArLink test exist")

	r7, err := c.GetArLink(ctx, &ArContent.GetArLinkRequest{
		ArContentId: arContentId,
	})

	testUtils.DisplayResult(r7, err, false)

	utils.PrintTitle("GetViewerData")

	r8, err := c.GetViewerData(ctx, &ArContent.GetViewerDataRequest{
		ArContentViewerUrlId: strings.Split(r7.ArContentViewerPath, "/")[4],
	})

	testUtils.DisplayResult(r8, err, false)

	TestCopyAndDeleteArContent()
}

func checkTemplateImageExist(imageName string) bool {
	checkOk := true
	res, err := s3Utils.FileExists(imageName, ArContent.ArContentImageType_TEMPLATE_IMAGE)
	utils.PrintObj(res, "checkTemplateImageExist")

	if err != nil {
		utils.PrintObj(err.Error(), "err")
		return false
	}

	checkOk = res

	return checkOk
}

func markFail(mark string) bool {
	utils.PrintObj(mark, "markFail")
	return false
}
