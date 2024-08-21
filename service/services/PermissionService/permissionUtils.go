package PermissionService

import (
	"service/model"
	"service/protos/Permission"
	"service/repositories"
	"service/utils"
)

var (
	PermissionModels             = GetAllPermissionModels()
	DefaultPermissionGroupUserId = "66acc133-ae77-4793-bd73-ce7b10740468"
)

func InitDefailtPermissionGroup() error {
	utils.PrintObj("", "InitDefailtPermissionGroup")

	// declare model
	defaultPermissionGroup := []model.PermissionGroup{
		{
			Id:                 utils.ParseUUID("abea1cbe-3f10-449c-be3d-227c120ca56e"),
			UserId:             utils.ParseUUID(DefaultPermissionGroupUserId),
			Name:               "聯絡人管理",
			SeletedIdArrayJson: utils.ToJson([]string{"201", "202", "203", "204"}),
		},
	}

	for _, pGroup := range defaultPermissionGroup {
		res, _ := repositories.QueryPermissionGroup(model.PermissionGroup{Id: pGroup.Id})

		if res.Error != nil {
			if utils.IsErrorNotFound(res.Error) {
				utils.PrintObj("cant found default pGroup. wait for create")

				res, _ := repositories.CreatePermissionGroup(pGroup)

				if res.Error != nil {
					utils.PrintObj(res.Error.Error(), "InitDefailtPermissionGrocp err")
					return res.Error
				}

				utils.PrintObj("create success, pGroup name: " + pGroup.Name)
				return nil
			} else {

				utils.PrintObj(res.Error.Error(), "InitDefailtPermissionGrocp err")
				return res.Error
			}
		}

		utils.PrintObj("found default pGroup.pGroup name: " + pGroup.Name)
	}

	return nil
}

func ValidPermissionId(arr []string) bool {
	utils.PrintObj(arr, "ValidPermissionId")

	pMap := GetPermissionIdsMap()
	for _, item := range arr {
		_, ok := pMap[item]
		if !ok {
			utils.PrintObj("found id not exist")
			return false
		}
	}

	checkDup := CheckDuplicateInArray(arr)
	if !checkDup {
		utils.PrintObj("found Duplicate id")
		return false
	}

	return true
}

func CheckDuplicateInArray[T any](arr []T) bool {
	visited := make(map[any]bool, 0)
	for i := 0; i < len(arr); i++ {
		if visited[arr[i]] {
			return false
		} else {
			visited[arr[i]] = true
		}
	}
	return true
}

func CheckUserPermission(user *model.User, funcName string) bool {
	utils.PrintObj(funcName, "CheckUserPermission")

	if !utils.IsUserChild(user) {
		utils.PrintObj("admin wont need permission check", "")
		return true
	}

	// get permission group
	res, data := repositories.QueryPermissionGroup(model.PermissionGroup{
		Id: user.PermissionGroupId.UUID,
	})

	if res.Error != nil {
		utils.PrintObj("cant found permission group", "")
		return false
	}

	// check permission
	// get func name id
	funcNameMap := GetPermissionRpcNameMap()
	val, ok := funcNameMap[funcName]
	if !ok {
		utils.PrintObj("cant found permission by fun name in map", "")
		return false
	}

	utils.PrintObj(val, "")

	// try to find id in user permission group
	selectedIdArray, err := utils.ParseJsonWithType[[]string](data.SeletedIdArrayJson)
	if err != nil {
		utils.PrintObj(err.Error())
		return false
	}

	for _, id := range selectedIdArray {
		if id == val {
			utils.PrintObj("found permission id in group", "")
			return true
		}
	}

	utils.PrintObj("cant found permission id in group", "")
	return false
}

func GetPermissionIdsMap() map[string]string {
	permissionMap := map[string]string{}
	packages := PermissionModels
	for _, packageModel := range packages {
		models := packageModel.PermissionModels
		if len(models) == 0 {
			utils.PrintObj("packageModel permissionModels has no models")
		}

		for _, modelItem := range models {
			permissionMap[modelItem.Id] = modelItem.RpcName
		}
	}

	return permissionMap
}

func GetPermissionRpcNameMap() map[string]string {
	permissionMap := map[string]string{}
	packages := PermissionModels
	for _, packageModel := range packages {
		models := packageModel.PermissionModels
		if len(models) == 0 {
			utils.PrintObj("packageModel permissionModels has no models")
		}

		for _, modelItem := range models {
			permissionMap[modelItem.RpcName] = modelItem.Id
		}
	}

	return permissionMap
}

func GetPermissionIdsArray() []string {
	arr := []string{}
	packages := PermissionModels
	for _, packageModel := range packages {
		models := packageModel.PermissionModels
		if len(models) == 0 {
			utils.PrintObj("packageModel permissionModels has no models")
		}

		for _, modelItem := range models {
			arr = append(arr, modelItem.Id)
		}
	}

	return arr
}

func AddPermission(modelPackage *Permission.PermissionPackage,
	models []*Permission.PermissionModel,
	modelPackages []*Permission.PermissionPackage) (modelPackagesReturn []*Permission.PermissionPackage) {

	modelPackage.PermissionModels = models
	modelPackages = append(modelPackages, modelPackage)

	return modelPackages
}
