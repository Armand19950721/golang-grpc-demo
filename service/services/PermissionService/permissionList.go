package PermissionService

import (
	"service/protos/Permission"
)

func GetAllPermissionModels() []*Permission.PermissionPackage {
	modelPackages := []*Permission.PermissionPackage{}

	modelPackages = AddPermission(
		&Permission.PermissionPackage{
			Name:         "使用者管理",
			ServicesName: "UserServices",
		},
		[]*Permission.PermissionModel{
			{Id: "101", RpcName: "UpdateCompany", Name: "更新公司資料"},
			{Id: "102", RpcName: "GetCompany", Name: "取得公司資料"},
		},
		modelPackages,
	)

	modelPackages = AddPermission(
		&Permission.PermissionPackage{
			Name:         "聯絡人管理",
			ServicesName: "ContactServices",
		},
		[]*Permission.PermissionModel{
			{Id: "201", RpcName: "AddContact", Name: "新增聯絡人"},
			{Id: "202", RpcName: "GetListContact", Name: "取得所有聯絡人"},
			{Id: "203", RpcName: "UpdateContact", Name: "更新聯絡人"},
			{Id: "204", RpcName: "DeleteContact", Name: "刪除聯絡人"},
		},
		modelPackages,
	)

	modelPackages = AddPermission(
		&Permission.PermissionPackage{
			Name:         "權限管理",
			ServicesName: "PermissionServices",
		},
		[]*Permission.PermissionModel{
			// {Id: "301", RpcName: "GetAllPermission", Name: "取得所有權限條目"},
			{Id: "302", RpcName: "AddPermissionGroup", Name: "加入權限群組"},
			{Id: "303", RpcName: "UpdatePermissionGroup", Name: "更新權限群組"},
			{Id: "304", RpcName: "GetPermissionGroupList", Name: "取得所有權限群組"},
		},
		modelPackages,
	)
	// utils.PrintObj(modelPackages)

	return modelPackages
}
