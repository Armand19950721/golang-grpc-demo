package repositories

import (
	"service/model"
	"service/protos/ArContent"
	"service/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateArContent(arContentData model.ArContent) (*gorm.DB, model.ArContent) {
	result := utils.DatabaseManager.Create(&arContentData)

	return result, arContentData
}

func GetArContentList(basicInfo utils.BasicInfo, req *ArContent.GetArContentListRequest) (tx *gorm.DB, models []model.ArContent, totalCount int64) {
	result := utils.DatabaseManagerSlave.Order("create_at desc").Where("user_id = ?", basicInfo.AdminId)

	if req.Keyword != "" {
		result.Where("name ILIKE ?", utils.GetSqlLikeString(req.Keyword)).Or("tag ILIKE ?", utils.GetSqlLikeString(req.Keyword))
	}

	result.Debug().Find(&models).Limit(int(req.PageInfo.PageItemCount)).Offset(int(req.PageInfo.PageItemCount) * (int(req.PageInfo.CurrentPageNum) - 1)).Count(&totalCount)

	utils.PrintObj(totalCount, "totalCount")

	return result, models, totalCount
}

func GetArContent(arContentWhere model.ArContent) (tx *gorm.DB, data model.ArContent) {
	result := utils.DatabaseManagerSlave.Where(arContentWhere).First(&data)

	return result, data
}

func GetArContentCount(arContentWhere model.ArContent) (tx *gorm.DB, count int64) {
	data := model.ArContent{}

	result := utils.DatabaseManagerSlave.Debug().Where(arContentWhere).Find(&data).Count(&count)
	utils.PrintObj(count, "GetArContentCount")

	return result, count
}

func UpdateArContent(ArContentWhere, arContent model.ArContent) (tx *gorm.DB) {
	result := utils.DatabaseManager.Where(ArContentWhere).Updates(&arContent)

	return result
}

func DeleteArContent(arContentWhere, arContent model.ArContent) (*gorm.DB, model.ArContent) {
	result := utils.DatabaseManager.Where(arContentWhere).Delete(&arContent)

	return result, arContent
}

func UpdateArContentTemplate(arContentWhere, arContent model.ArContent) (*gorm.DB, model.ArContent) {
	result := utils.DatabaseManager.Where(arContentWhere).Updates(&arContent)

	return result, arContent
}

func UpdateArContentViewer(id uuid.UUID, userId uuid.UUID, viewerSetting string) (*gorm.DB, model.ArContent) {
	qResult, data := GetArContent(model.ArContent{Id: id, UserId: userId})

	if qResult != nil {
		return qResult, data
	}
	data.ViewerSetting = viewerSetting

	result := utils.DatabaseManager.Save(&data)

	return result, data
}
