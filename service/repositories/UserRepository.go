package repositories

import (
	"errors"
	"service/model"
	"service/protos/Common"
	"service/utils"

	"gorm.io/gorm"
)

func QueryUser(userQuery model.User) (tx *gorm.DB, data model.User) {
	result := utils.DatabaseManagerSlave.Where(userQuery).First(&data)

	return result, data
}

func QueryUserIncludeDel(userQuery model.User) (tx *gorm.DB, data model.User) {
	result := utils.DatabaseManagerSlave.Unscoped().Where(userQuery).First(&data)

	return result, data
}

func CreateUser(User model.User) (*gorm.DB, model.User) {
	result := utils.DatabaseManager.Create(&User)

	return result, User
}

func UpdateUser(userWhere, user model.User) (tx *gorm.DB, data model.User) {
	result := utils.DatabaseManager.Where(userWhere).Updates(&user).Find(&user)

	return result, user
}

func GetUserCount(userWhere model.User) (*gorm.DB, int64) {
	totalCount := int64(0)
	result := utils.DatabaseManagerSlave.Where(userWhere).Find(&[]model.User{}).Count(&totalCount)

	utils.PrintObj(totalCount, "GetUserCount")

	return result, totalCount
}

func GetUserList(userWhere model.User, pageInfo *Common.PageInfoRequest) (*gorm.DB, []model.User, int64) {
	totalCount := int64(0)
	models := []model.User{}

	result := utils.DatabaseManagerSlave.Limit(int(pageInfo.PageItemCount)).Offset(int(pageInfo.PageItemCount) * (int(pageInfo.CurrentPageNum) - 1)).Where(userWhere).Order("create_at desc").Find(&models).Count(&totalCount)

	utils.PrintObj(totalCount, "SelectUser totalCount")
	utils.PrintObj(result.RowsAffected, "SelectUser RowsAffected")

	return result, models, totalCount
}

func DeleteUser(user model.User) (tx *gorm.DB, data model.User) {
	result := utils.DatabaseManager.Delete(&user)

	return result, user
}

func HardDeleteUser(user model.User) (tx *gorm.DB, data model.User) {
	// make sure only user child can be hard delete
	res, data := QueryUser(user)
	if res.Error != nil {
		return res, data
	}

	if utils.IsUserChild(&data) {
		utils.PrintObj("delete child", "HardDeleteUser")
		result := utils.DatabaseManager.Debug().Unscoped().Delete(&user)
		return result, user
	} else {
		utils.PrintObj("cant delete admin", "HardDeleteUser")
		return &gorm.DB{Error: errors.New("only child can hard delete")}, model.User{}
	}
}
