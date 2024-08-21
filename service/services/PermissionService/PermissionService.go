package PermissionService

import (
	// "service/model"
	// "service/protos/Common"
	// "service/repositories"
	"context"
	"service/model"
	"service/protos/Common"
	"service/protos/Permission"
	"service/repositories"

	"service/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetPermissionGroupDefaultList(ctx context.Context, req *Permission.GetDefaultPermissionGroupListRequest) (*Permission.GetDefaultPermissionGroupListReply, error) {
	reply := &Permission.GetDefaultPermissionGroupListReply{}

	// get default
	res, pGroups, _ := repositories.GetPermissionGroupList(model.PermissionGroup{UserId: utils.ParseUUID(DefaultPermissionGroupUserId)})

	if res.Error != nil {
		return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
			Code:        Common.ErrorCodes_UNKNOWN_ERROR,
			InternalMsg: res.Error.Error(),
		}))
	}

	// convert to proto
	models := []*Permission.PermissionGroupModel{}

	for _, groupItem := range pGroups {

		permissionSeletedIds, err := utils.ParseJsonWithType[[]string](groupItem.SeletedIdArrayJson)
		if err != nil {
			return reply, utils.ReturnUnKnownError(res.Error)
		}

		model := &Permission.PermissionGroupModel{
			PermissionGroupId:    groupItem.Id.String(),
			PermissionGroupName:  groupItem.Name,
			PermissionSeletedIds: permissionSeletedIds,
		}
		models = append(models, model)
	}

	reply.Models = models

	return reply, status.Errorf(codes.OK, "")
}

/// 棄用 (暫時不讓使用者客製權限 改為先準備預設權限Group)
// var (
// 	permissionGroupLimit = 5
// )

// func GetAllPermission(ctx context.Context, req *Permission.GetAllPermissionRequest) (*Permission.GetAllPermissionReply, error) {
// 	reply := &Permission.GetAllPermissionReply{}
// 	reply.Models = permissionUtils.PermissionModels

// 	return reply, status.Errorf(codes.OK, "")
// }

// func AddPermissionGroup(ctx context.Context, req *Permission.AddPermissionGroupRequest) (*Common.CommonReply, error) {
// 	reply := &Common.CommonReply{}

// 	// check param
// 	if !permissionUtils.ValidPermissionId(req.Model.PermissionSeletedIds) ||
// 		!utils.ValidString(req.Model.PermissionGroupName, 1, 100) {
// 		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
// 	}

// 	// get basic data
// 	basicInfo := utils.GetBasicInfo(ctx)
// 	if !basicInfo.Success {
// 				return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
// 	Code:        Common.ErrorCodes_UNKNOWN_ERROR,
// 	InternalMsg: "basicInfo error",
// }))
// 	}

// 	// check permission group count
// 	checkRes, _, count := repositories.GetPermissionGroupList(model.PermissionGroup{UserId: basicInfo.AdminId})
// 	if checkRes.Error != nil {
// 					return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
// 	Code:        Common.ErrorCodes_UNKNOWN_ERROR,
// 	InternalMsg: res.Error.Error(),
// }))
// 	}

// 	utils.PrintObj([]int{int(count), int(permissionGroupLimit)}, "count vs permissionGroupLimit")

// 	if int(count) >= permissionGroupLimit {
// 		return reply, status.Errorf(codes.OutOfRange, "out of limit")
// 	}

// 	// create
// 	pGroup := model.PermissionGroup{
// 		UserId:             basicInfo.AdminId,
// 		Name:               req.Model.PermissionGroupName,
// 		SeletedIdArrayJson: utils.ToJson(req.Model.PermissionSeletedIds),
// 	}

// 	res, pGroup := repositories.CreatePermissionGroup(pGroup)
// 	if res.Error != nil {
// 					return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
// 	Code:        Common.ErrorCodes_UNKNOWN_ERROR,
// 	InternalMsg: res.Error.Error(),
// }))
// 	}

// 	return reply, status.Errorf(codes.OK, "")
// }

// func UpdatePermissionGroup(ctx context.Context, req *Permission.UpdatePermissionGroupRequest) (*Common.CommonReply, error) {
// 	reply := &Common.CommonReply{}

// 	// check param
// 	if !utils.ValidId(req.Model.PermissionGroupId) ||
// 		!permissionUtils.ValidPermissionId(req.Model.PermissionSeletedIds) ||
// 		!utils.ValidString(req.Model.PermissionGroupName, 1, 100) {

// 		return reply, status.Errorf(codes.InvalidArgument, utils.GetError(utils.ErrorType{Code: Common.ErrorCodes_INVAILD_PARAM}))
// 	}

// 	// get basic data
// 	basicInfo := utils.GetBasicInfo(ctx)
// 	if !basicInfo.Success {
// 				return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
// 	Code:        Common.ErrorCodes_UNKNOWN_ERROR,
// 	InternalMsg: "basicInfo error",
// }))
// 	}

// 	// update
// 	pGroupWhere := model.PermissionGroup{
// 		Id: utils.ParseUUID(req.Model.PermissionGroupId),
// 	}

// 	pGroup := model.PermissionGroup{
// 		UserId:             basicInfo.AdminId,
// 		Name:               req.Model.PermissionGroupName,
// 		SeletedIdArrayJson: utils.ToJson(req.Model.PermissionSeletedIds),
// 	}

// 	res, pGroup := repositories.UpdatePermissionGroup(pGroupWhere, pGroup)
// 	if res.Error != nil {
// 					return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
// 	Code:        Common.ErrorCodes_UNKNOWN_ERROR,
// 	InternalMsg: res.Error.Error(),
// }))
// 	}

// 	return reply, status.Errorf(codes.OK, "")
// }

// func GetPermissionGroupList(ctx context.Context, req *Permission.GetPermissionGroupListRequest) (*Permission.GetPermissionGroupListReply, error) {
// 	reply := &Permission.GetPermissionGroupListReply{}

// 	// get basic data
// 	basicInfo := utils.GetBasicInfo(ctx)
// 	if !basicInfo.Success {
// 				return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
// 	Code:        Common.ErrorCodes_UNKNOWN_ERROR,
// 	InternalMsg: "basicInfo error",
// }))
// 	}

// 	// get list by admin id
// 	pGroupWhere := model.PermissionGroup{
// 		UserId: basicInfo.AdminId,
// 	}

// 	res, pGroups, _ := repositories.GetPermissionGroupList(pGroupWhere)
// 	if res.Error != nil {
// 					return reply, status.Errorf(codes.Internal, utils.GetError(utils.ErrorType{
// 	Code:        Common.ErrorCodes_UNKNOWN_ERROR,
// 	InternalMsg: res.Error.Error(),
// }))
// 	}

// 	models := []*Permission.PermissionGroupModel{}

// 	// convert to proto
// 	for _, groupItem := range pGroups {

// 		permissionSeletedIds, err := utils.ParseJsonWithType[[]string](groupItem.SeletedIdArrayJson)
// 		if err != nil {
// 			return reply, err
// 		}

// 		model := &Permission.PermissionGroupModel{
// 			PermissionGroupId:    groupItem.Id.String(),
// 			PermissionGroupName:  groupItem.Name,
// 			PermissionSeletedIds: permissionSeletedIds,
// 		}
// 		models = append(models, model)
// 	}

// 	reply.Models = models
// 	return reply, status.Errorf(codes.OK, "")
// }
