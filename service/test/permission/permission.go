package permission

import (
	"service/test/testUtils"
	"service/utils"
	"time"

	// "service/protos/Common"
	"service/protos/Permission"
	"service/protos/WebServices"
	// "service/services/permissionUtils"
)

const (
	// defaultName = "world"
	timestampFormat = time.StampNano // "Jan _2 15:04:05.000"
)

var (
	c    = WebServices.NewPermissionServicesClient(testUtils.Conn)
	name = testUtils.RandomIntString() + "group name"
)

func Flow() {
	ctx := testUtils.GetCtx()
	utils.PrintTitle("GetPermissionGroupDefaultList")

	r, err := c.GetPermissionGroupDefaultList(ctx, &Permission.GetDefaultPermissionGroupListRequest{})
	testUtils.DisplayResult(r, err, false)
}

/// 棄用 (暫時不讓使用者客製權限 改為先準備預設權限Group)
// func PGLimit() {
// 	ctx := testUtils.GetCtx()

// 	utils.PrintTitle("add permission group")
// 	for i := 0; i < 5; i++ {
// 		model := &Permission.PermissionGroupModel{
// 			PermissionGroupName:  testUtils.RandomIntString() + "group name",
// 			PermissionSeletedIds: permissionUtils.GetPermissionIdsArray(),
// 		}

// 		rR, err := c.AddPermissionGroup(ctx, &Permission.AddPermissionGroupRequest{
// 			Model: model,
// 		})
// 		testUtils.DisplayResult(rR, err, false)
// 	}
// }

// func Flow() {
// 	ctx := testUtils.GetCtx()
// 	utils.PrintTitle("permission get all")

// 	r, err := c.GetAllPermission(ctx, &Permission.GetAllPermissionRequest{})
// 	testUtils.DisplayResult(r, err, false)

// 	utils.PrintTitle("add permission group")

// 	model := &Permission.PermissionGroupModel{
// 		PermissionGroupName:  name,
// 		PermissionSeletedIds: permissionUtils.GetPermissionIdsArray(),
// 	}

// 	rR, err := c.AddPermissionGroup(ctx, &Permission.AddPermissionGroupRequest{
// 		Model: model,
// 	})

// 	testUtils.DisplayResult(rR, err, false)

// 	utils.PrintTitle("add permission group fail (wrong id)")

// 	model = &Permission.PermissionGroupModel{
// 		PermissionGroupName:  name,
// 		PermissionSeletedIds: []string{"1", "3"},
// 	}

// 	rR, err = c.AddPermissionGroup(ctx, &Permission.AddPermissionGroupRequest{
// 		Model: model,
// 	})

// 	testUtils.DisplayResult(rR, err, true)

// 	utils.PrintTitle("get list permission group")

// 	rRR, err := c.GetPermissionGroupList(ctx, &Permission.GetPermissionGroupListRequest{})

// 	testUtils.DisplayResult(rR, err, false)

// 	utils.PrintTitle("update permission group")
// 	model.PermissionGroupId = rRR.Models[0].PermissionGroupId
// 	model.PermissionGroupName += "_update"
// 	model.PermissionSeletedIds = permissionUtils.GetPermissionIdsArray()
// 	rRRR, err := c.UpdatePermissionGroup(ctx, &Permission.UpdatePermissionGroupRequest{
// 		Model: model,
// 	})

// 	testUtils.DisplayResult(rRRR, err, false)

// 	utils.PrintTitle("get list permission group two")

// 	rRRA, err := c.GetPermissionGroupList(ctx, &Permission.GetPermissionGroupListRequest{})

// 	testUtils.DisplayResult(rRRA, err, false)
// }
