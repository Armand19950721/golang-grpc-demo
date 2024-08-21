// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.5
// source: ThirdPartyWebServices.proto

package ThirdPartyWebServices

import (
	ThirdPartyArContent "service/protos/ThirdPartyArContent"
	ThirdPartyAuth "service/protos/ThirdPartyAuth"
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ThirdPartyArContentServicesClient is the client API for ThirdPartyArContentServices service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ThirdPartyArContentServicesClient interface {
	CreateArContent(ctx context.Context, in *ThirdPartyArContent.CreateArContentRequest, opts ...grpc.CallOption) (*ThirdPartyArContent.CreateArContentReply, error)
	GetArLink(ctx context.Context, in *ThirdPartyArContent.GetArLinkRequest, opts ...grpc.CallOption) (*ThirdPartyArContent.GetArLinkReply, error)
}

type thirdPartyArContentServicesClient struct {
	cc grpc.ClientConnInterface
}

func NewThirdPartyArContentServicesClient(cc grpc.ClientConnInterface) ThirdPartyArContentServicesClient {
	return &thirdPartyArContentServicesClient{cc}
}

func (c *thirdPartyArContentServicesClient) CreateArContent(ctx context.Context, in *ThirdPartyArContent.CreateArContentRequest, opts ...grpc.CallOption) (*ThirdPartyArContent.CreateArContentReply, error) {
	out := new(ThirdPartyArContent.CreateArContentReply)
	err := c.cc.Invoke(ctx, "/ThirdPartyWebServices.ThirdPartyArContentServices/CreateArContent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *thirdPartyArContentServicesClient) GetArLink(ctx context.Context, in *ThirdPartyArContent.GetArLinkRequest, opts ...grpc.CallOption) (*ThirdPartyArContent.GetArLinkReply, error) {
	out := new(ThirdPartyArContent.GetArLinkReply)
	err := c.cc.Invoke(ctx, "/ThirdPartyWebServices.ThirdPartyArContentServices/GetArLink", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ThirdPartyArContentServicesServer is the server API for ThirdPartyArContentServices service.
// All implementations must embed UnimplementedThirdPartyArContentServicesServer
// for forward compatibility
type ThirdPartyArContentServicesServer interface {
	CreateArContent(context.Context, *ThirdPartyArContent.CreateArContentRequest) (*ThirdPartyArContent.CreateArContentReply, error)
	GetArLink(context.Context, *ThirdPartyArContent.GetArLinkRequest) (*ThirdPartyArContent.GetArLinkReply, error)
	mustEmbedUnimplementedThirdPartyArContentServicesServer()
}

// UnimplementedThirdPartyArContentServicesServer must be embedded to have forward compatible implementations.
type UnimplementedThirdPartyArContentServicesServer struct {
}

func (UnimplementedThirdPartyArContentServicesServer) CreateArContent(context.Context, *ThirdPartyArContent.CreateArContentRequest) (*ThirdPartyArContent.CreateArContentReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateArContent not implemented")
}
func (UnimplementedThirdPartyArContentServicesServer) GetArLink(context.Context, *ThirdPartyArContent.GetArLinkRequest) (*ThirdPartyArContent.GetArLinkReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArLink not implemented")
}
func (UnimplementedThirdPartyArContentServicesServer) mustEmbedUnimplementedThirdPartyArContentServicesServer() {
}

// UnsafeThirdPartyArContentServicesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ThirdPartyArContentServicesServer will
// result in compilation errors.
type UnsafeThirdPartyArContentServicesServer interface {
	mustEmbedUnimplementedThirdPartyArContentServicesServer()
}

func RegisterThirdPartyArContentServicesServer(s grpc.ServiceRegistrar, srv ThirdPartyArContentServicesServer) {
	s.RegisterService(&ThirdPartyArContentServices_ServiceDesc, srv)
}

func _ThirdPartyArContentServices_CreateArContent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ThirdPartyArContent.CreateArContentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ThirdPartyArContentServicesServer).CreateArContent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ThirdPartyWebServices.ThirdPartyArContentServices/CreateArContent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ThirdPartyArContentServicesServer).CreateArContent(ctx, req.(*ThirdPartyArContent.CreateArContentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ThirdPartyArContentServices_GetArLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ThirdPartyArContent.GetArLinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ThirdPartyArContentServicesServer).GetArLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ThirdPartyWebServices.ThirdPartyArContentServices/GetArLink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ThirdPartyArContentServicesServer).GetArLink(ctx, req.(*ThirdPartyArContent.GetArLinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ThirdPartyArContentServices_ServiceDesc is the grpc.ServiceDesc for ThirdPartyArContentServices service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ThirdPartyArContentServices_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ThirdPartyWebServices.ThirdPartyArContentServices",
	HandlerType: (*ThirdPartyArContentServicesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateArContent",
			Handler:    _ThirdPartyArContentServices_CreateArContent_Handler,
		},
		{
			MethodName: "GetArLink",
			Handler:    _ThirdPartyArContentServices_GetArLink_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ThirdPartyWebServices.proto",
}

// ThirdPartyAuthServicesClient is the client API for ThirdPartyAuthServices service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ThirdPartyAuthServicesClient interface {
	Verification(ctx context.Context, in *ThirdPartyAuth.VerificationRequest, opts ...grpc.CallOption) (*ThirdPartyAuth.VerificationReply, error)
}

type thirdPartyAuthServicesClient struct {
	cc grpc.ClientConnInterface
}

func NewThirdPartyAuthServicesClient(cc grpc.ClientConnInterface) ThirdPartyAuthServicesClient {
	return &thirdPartyAuthServicesClient{cc}
}

func (c *thirdPartyAuthServicesClient) Verification(ctx context.Context, in *ThirdPartyAuth.VerificationRequest, opts ...grpc.CallOption) (*ThirdPartyAuth.VerificationReply, error) {
	out := new(ThirdPartyAuth.VerificationReply)
	err := c.cc.Invoke(ctx, "/ThirdPartyWebServices.ThirdPartyAuthServices/Verification", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ThirdPartyAuthServicesServer is the server API for ThirdPartyAuthServices service.
// All implementations must embed UnimplementedThirdPartyAuthServicesServer
// for forward compatibility
type ThirdPartyAuthServicesServer interface {
	Verification(context.Context, *ThirdPartyAuth.VerificationRequest) (*ThirdPartyAuth.VerificationReply, error)
	mustEmbedUnimplementedThirdPartyAuthServicesServer()
}

// UnimplementedThirdPartyAuthServicesServer must be embedded to have forward compatible implementations.
type UnimplementedThirdPartyAuthServicesServer struct {
}

func (UnimplementedThirdPartyAuthServicesServer) Verification(context.Context, *ThirdPartyAuth.VerificationRequest) (*ThirdPartyAuth.VerificationReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Verification not implemented")
}
func (UnimplementedThirdPartyAuthServicesServer) mustEmbedUnimplementedThirdPartyAuthServicesServer() {
}

// UnsafeThirdPartyAuthServicesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ThirdPartyAuthServicesServer will
// result in compilation errors.
type UnsafeThirdPartyAuthServicesServer interface {
	mustEmbedUnimplementedThirdPartyAuthServicesServer()
}

func RegisterThirdPartyAuthServicesServer(s grpc.ServiceRegistrar, srv ThirdPartyAuthServicesServer) {
	s.RegisterService(&ThirdPartyAuthServices_ServiceDesc, srv)
}

func _ThirdPartyAuthServices_Verification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ThirdPartyAuth.VerificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ThirdPartyAuthServicesServer).Verification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ThirdPartyWebServices.ThirdPartyAuthServices/Verification",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ThirdPartyAuthServicesServer).Verification(ctx, req.(*ThirdPartyAuth.VerificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ThirdPartyAuthServices_ServiceDesc is the grpc.ServiceDesc for ThirdPartyAuthServices service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ThirdPartyAuthServices_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ThirdPartyWebServices.ThirdPartyAuthServices",
	HandlerType: (*ThirdPartyAuthServicesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Verification",
			Handler:    _ThirdPartyAuthServices_Verification_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ThirdPartyWebServices.proto",
}
