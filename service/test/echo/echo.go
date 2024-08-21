package echo

import (
	"service/test/testUtils"
	"service/utils"
	// "log"
	// "time"

	"service/protos/Echo"
	"service/protos/WebServices"
)

var (
	ctx = testUtils.GetCtx()
	c   = WebServices.NewEchoServicesClient(testUtils.Conn)
)

func Flow() {
	utils.PrintTitle("echo")

	r, err := c.Echo(ctx, &Echo.EchoRequest{
		Msg: "hi~~",
	})
	testUtils.DisplayResult(r, err, false)

	utils.PrintTitle("echo token")

	rR, err := c.EchoToken(ctx, &Echo.EchoRequest{
		Msg: "hi~~",
	})
	testUtils.DisplayResult(rR, err, false)

}

func FlowEchoEmpty() {
	utils.PrintTitle("echo")

	r, err := c.Echo(ctx, &Echo.EchoRequest{
		Msg: "",
	})
	testUtils.DisplayResult(r, err, false)
}
