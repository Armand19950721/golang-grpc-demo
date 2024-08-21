package childPermission

/// 棄用 (暫時不讓使用者客製權限 改為先準備預設權限Group)
// import (
// 	"service/test/testUtils"
// 	"service/utils"

// 	"service/protos/Contact"
// 	"service/protos/Echo"
// 	"service/protos/Permission"
// 	"service/protos/WebServices"
// 	"service/utils/permissionUtils"
// )

// var (
// 	p = WebServices.NewPermissionServicesClient(testUtils.Conn)
// 	c = WebServices.NewContactServicesClient(testUtils.Conn)
// 	e = WebServices.NewEchoServicesClient(testUtils.Conn)
// )

// func Flow() {
// 	ctx := testUtils.GetCtx("use Child Token")
// 	//抓取某child group id
// 	utils.PrintTitle("get list permission group")

// 	r, err := p.GetPermissionGroupList(ctx, &Permission.GetPermissionGroupListRequest{})

// 	testUtils.DisplayResult(r, err, false)

// 	// call 允許的rpc fun
// 	utils.PrintTitle("GetListContact")

// 	rl, err := c.GetListContact(ctx, &Contact.GetListContactRequest{})
// 	testUtils.DisplayResult(rl, err, false)

// 	//修改group id
// 	utils.PrintTitle("update permission group")

// 	model := r.Models[0]
// 	model.PermissionGroupName = "update_by_child"
// 	model.PermissionSeletedIds = utils.RemoveFromArrayString(
// 		permissionUtils.GetPermissionIdsArray(), "202")

// 	rld, err := p.UpdatePermissionGroup(ctx, &Permission.UpdatePermissionGroupRequest{
// 		Model: model,
// 	})

// 	testUtils.DisplayResult(rld, err, false)

// 	//測試call 被修改的rpc fun
// 	utils.PrintTitle("GetListContact err")

// 	rl, err = c.GetListContact(ctx, &Contact.GetListContactRequest{})
// 	testUtils.DisplayResult(rl, err, true)

// 	// 確認有被改到
// 	utils.PrintTitle("get list permission group")

// 	r, err = p.GetPermissionGroupList(ctx, &Permission.GetPermissionGroupListRequest{})

// 	testUtils.DisplayResult(r, err, false)

// 	// 修改group id
// 	utils.PrintTitle("update permission group to normal")

// 	model.PermissionGroupName = "update_by_child_recover"
// 	model.PermissionSeletedIds = permissionUtils.GetPermissionIdsArray()

// 	rld, err = p.UpdatePermissionGroup(ctx, &Permission.UpdatePermissionGroupRequest{
// 		Model: model,
// 	})

// 	testUtils.DisplayResult(rld, err, false)

// 	// 測試未被登錄的rpc function
// 	utils.PrintTitle("test not reg in pList")

// 	rR, err := e.EchoToken(ctx, &Echo.EchoRequest{
// 		Msg: "hi~~",
// 	})

// 	testUtils.DisplayResult(rR, err, true)
// }
