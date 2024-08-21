package controller

import (
	"context"
	"fmt"
	"service/protos/Echo"
	"service/protos/WebServices"
	"time"
)

type EchoController struct {
	WebServices.EchoServicesServer
}

func (controller *EchoController) Echo(ctx context.Context, req *Echo.EchoRequest) (*Echo.EchoReply, error) {

	reply := &Echo.EchoReply{Result: true,
		Msg: fmt.Sprintf("[%s]:from server: %s",
			time.Now().Format("2006-01-02 15-04-05-Z07"),
			req.Msg)}

	if req.Msg == "" {
		panic("echo is empty")
	} else {
		return reply, nil
	}

}

func (controller *EchoController) EchoToken(ctx context.Context, req *Echo.EchoRequest) (*Echo.EchoReply, error) {
	return &Echo.EchoReply{Result: true,
		Msg: fmt.Sprintf("[%s]:from server: %s",
			time.Now().Format("2006-01-02 15-04-05-Z07"),
			req.Msg)}, nil
}
