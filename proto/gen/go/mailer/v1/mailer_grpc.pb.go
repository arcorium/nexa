// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: mailer/v1/mailer.proto

package mailerv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	common "nexa/proto/gen/go/common"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	MailerService_Find_FullMethodName       = "/nexa.proto.generated.mailer.v1.MailerService/Find"
	MailerService_FindByIds_FullMethodName  = "/nexa.proto.generated.mailer.v1.MailerService/FindByIds"
	MailerService_FindByTag_FullMethodName  = "/nexa.proto.generated.mailer.v1.MailerService/FindByTag"
	MailerService_Send_FullMethodName       = "/nexa.proto.generated.mailer.v1.MailerService/Send"
	MailerService_SendToUser_FullMethodName = "/nexa.proto.generated.mailer.v1.MailerService/SendToUser"
	MailerService_Update_FullMethodName     = "/nexa.proto.generated.mailer.v1.MailerService/Update"
	MailerService_Remove_FullMethodName     = "/nexa.proto.generated.mailer.v1.MailerService/Remove"
)

// MailerServiceClient is the client API for MailerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MailerServiceClient interface {
	Find(ctx context.Context, in *common.PagedElementInput, opts ...grpc.CallOption) (*FindResponse, error)
	FindByIds(ctx context.Context, in *FindMailByIdsRequest, opts ...grpc.CallOption) (*FindMailByIdsResponse, error)
	FindByTag(ctx context.Context, in *FindMailByTagRequest, opts ...grpc.CallOption) (*FindMailByTagResponse, error)
	Send(ctx context.Context, in *SendMailRequest, opts ...grpc.CallOption) (*SendMailResponse, error)
	SendToUser(ctx context.Context, in *SendMailToUserRequest, opts ...grpc.CallOption) (*SendMailToUserResponse, error)
	Update(ctx context.Context, in *UpdateMailRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Remove(ctx context.Context, in *RemoveMailRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type mailerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMailerServiceClient(cc grpc.ClientConnInterface) MailerServiceClient {
	return &mailerServiceClient{cc}
}

func (c *mailerServiceClient) Find(ctx context.Context, in *common.PagedElementInput, opts ...grpc.CallOption) (*FindResponse, error) {
	out := new(FindResponse)
	err := c.cc.Invoke(ctx, MailerService_Find_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailerServiceClient) FindByIds(ctx context.Context, in *FindMailByIdsRequest, opts ...grpc.CallOption) (*FindMailByIdsResponse, error) {
	out := new(FindMailByIdsResponse)
	err := c.cc.Invoke(ctx, MailerService_FindByIds_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailerServiceClient) FindByTag(ctx context.Context, in *FindMailByTagRequest, opts ...grpc.CallOption) (*FindMailByTagResponse, error) {
	out := new(FindMailByTagResponse)
	err := c.cc.Invoke(ctx, MailerService_FindByTag_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailerServiceClient) Send(ctx context.Context, in *SendMailRequest, opts ...grpc.CallOption) (*SendMailResponse, error) {
	out := new(SendMailResponse)
	err := c.cc.Invoke(ctx, MailerService_Send_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailerServiceClient) SendToUser(ctx context.Context, in *SendMailToUserRequest, opts ...grpc.CallOption) (*SendMailToUserResponse, error) {
	out := new(SendMailToUserResponse)
	err := c.cc.Invoke(ctx, MailerService_SendToUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailerServiceClient) Update(ctx context.Context, in *UpdateMailRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, MailerService_Update_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mailerServiceClient) Remove(ctx context.Context, in *RemoveMailRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, MailerService_Remove_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MailerServiceServer is the server API for MailerService service.
// All implementations must embed UnimplementedMailerServiceServer
// for forward compatibility
type MailerServiceServer interface {
	Find(context.Context, *common.PagedElementInput) (*FindResponse, error)
	FindByIds(context.Context, *FindMailByIdsRequest) (*FindMailByIdsResponse, error)
	FindByTag(context.Context, *FindMailByTagRequest) (*FindMailByTagResponse, error)
	Send(context.Context, *SendMailRequest) (*SendMailResponse, error)
	SendToUser(context.Context, *SendMailToUserRequest) (*SendMailToUserResponse, error)
	Update(context.Context, *UpdateMailRequest) (*emptypb.Empty, error)
	Remove(context.Context, *RemoveMailRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedMailerServiceServer()
}

// UnimplementedMailerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMailerServiceServer struct {
}

func (UnimplementedMailerServiceServer) Find(context.Context, *common.PagedElementInput) (*FindResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Find not implemented")
}
func (UnimplementedMailerServiceServer) FindByIds(context.Context, *FindMailByIdsRequest) (*FindMailByIdsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByIds not implemented")
}
func (UnimplementedMailerServiceServer) FindByTag(context.Context, *FindMailByTagRequest) (*FindMailByTagResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByTag not implemented")
}
func (UnimplementedMailerServiceServer) Send(context.Context, *SendMailRequest) (*SendMailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Send not implemented")
}
func (UnimplementedMailerServiceServer) SendToUser(context.Context, *SendMailToUserRequest) (*SendMailToUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendToUser not implemented")
}
func (UnimplementedMailerServiceServer) Update(context.Context, *UpdateMailRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedMailerServiceServer) Remove(context.Context, *RemoveMailRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}
func (UnimplementedMailerServiceServer) mustEmbedUnimplementedMailerServiceServer() {}

// UnsafeMailerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MailerServiceServer will
// result in compilation errors.
type UnsafeMailerServiceServer interface {
	mustEmbedUnimplementedMailerServiceServer()
}

func RegisterMailerServiceServer(s grpc.ServiceRegistrar, srv MailerServiceServer) {
	s.RegisterService(&MailerService_ServiceDesc, srv)
}

func _MailerService_Find_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.PagedElementInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailerServiceServer).Find(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MailerService_Find_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailerServiceServer).Find(ctx, req.(*common.PagedElementInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailerService_FindByIds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindMailByIdsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailerServiceServer).FindByIds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MailerService_FindByIds_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailerServiceServer).FindByIds(ctx, req.(*FindMailByIdsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailerService_FindByTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindMailByTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailerServiceServer).FindByTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MailerService_FindByTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailerServiceServer).FindByTag(ctx, req.(*FindMailByTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailerService_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailerServiceServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MailerService_Send_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailerServiceServer).Send(ctx, req.(*SendMailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailerService_SendToUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMailToUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailerServiceServer).SendToUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MailerService_SendToUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailerServiceServer).SendToUser(ctx, req.(*SendMailToUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailerService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailerServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MailerService_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailerServiceServer).Update(ctx, req.(*UpdateMailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MailerService_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveMailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MailerServiceServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MailerService_Remove_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MailerServiceServer).Remove(ctx, req.(*RemoveMailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MailerService_ServiceDesc is the grpc.ServiceDesc for MailerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MailerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nexa.proto.generated.mailer.v1.MailerService",
	HandlerType: (*MailerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Find",
			Handler:    _MailerService_Find_Handler,
		},
		{
			MethodName: "FindByIds",
			Handler:    _MailerService_FindByIds_Handler,
		},
		{
			MethodName: "FindByTag",
			Handler:    _MailerService_FindByTag_Handler,
		},
		{
			MethodName: "Send",
			Handler:    _MailerService_Send_Handler,
		},
		{
			MethodName: "SendToUser",
			Handler:    _MailerService_SendToUser_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _MailerService_Update_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _MailerService_Remove_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mailer/v1/mailer.proto",
}
