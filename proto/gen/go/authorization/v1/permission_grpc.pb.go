// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: authorization/v1/permission.proto

package authorizationv1

import (
	context "context"
	common "github.com/arcorium/nexa/proto/gen/go/common"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	PermissionService_Create_FullMethodName      = "/nexa.authorization.v1.PermissionService/Create"
	PermissionService_Find_FullMethodName        = "/nexa.authorization.v1.PermissionService/Find"
	PermissionService_FindByRoles_FullMethodName = "/nexa.authorization.v1.PermissionService/FindByRoles"
	PermissionService_FindAll_FullMethodName     = "/nexa.authorization.v1.PermissionService/FindAll"
	PermissionService_Delete_FullMethodName      = "/nexa.authorization.v1.PermissionService/Delete"
)

// PermissionServiceClient is the client API for PermissionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PermissionServiceClient interface {
	Create(ctx context.Context, in *CreatePermissionRequest, opts ...grpc.CallOption) (*CreatePermissionResponse, error)
	Find(ctx context.Context, in *FindPermissionRequest, opts ...grpc.CallOption) (*FindPermissionResponse, error)
	FindByRoles(ctx context.Context, in *FindPermissionsByRoleRequest, opts ...grpc.CallOption) (*FindPermissionByRoleResponse, error)
	FindAll(ctx context.Context, in *common.PagedElementInput, opts ...grpc.CallOption) (*FindAllPermissionResponse, error)
	Delete(ctx context.Context, in *DeletePermissionRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type permissionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPermissionServiceClient(cc grpc.ClientConnInterface) PermissionServiceClient {
	return &permissionServiceClient{cc}
}

func (c *permissionServiceClient) Create(ctx context.Context, in *CreatePermissionRequest, opts ...grpc.CallOption) (*CreatePermissionResponse, error) {
	out := new(CreatePermissionResponse)
	err := c.cc.Invoke(ctx, PermissionService_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionServiceClient) Find(ctx context.Context, in *FindPermissionRequest, opts ...grpc.CallOption) (*FindPermissionResponse, error) {
	out := new(FindPermissionResponse)
	err := c.cc.Invoke(ctx, PermissionService_Find_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionServiceClient) FindByRoles(ctx context.Context, in *FindPermissionsByRoleRequest, opts ...grpc.CallOption) (*FindPermissionByRoleResponse, error) {
	out := new(FindPermissionByRoleResponse)
	err := c.cc.Invoke(ctx, PermissionService_FindByRoles_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionServiceClient) FindAll(ctx context.Context, in *common.PagedElementInput, opts ...grpc.CallOption) (*FindAllPermissionResponse, error) {
	out := new(FindAllPermissionResponse)
	err := c.cc.Invoke(ctx, PermissionService_FindAll_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionServiceClient) Delete(ctx context.Context, in *DeletePermissionRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PermissionService_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PermissionServiceServer is the server API for PermissionService service.
// All implementations must embed UnimplementedPermissionServiceServer
// for forward compatibility
type PermissionServiceServer interface {
	Create(context.Context, *CreatePermissionRequest) (*CreatePermissionResponse, error)
	Find(context.Context, *FindPermissionRequest) (*FindPermissionResponse, error)
	FindByRoles(context.Context, *FindPermissionsByRoleRequest) (*FindPermissionByRoleResponse, error)
	FindAll(context.Context, *common.PagedElementInput) (*FindAllPermissionResponse, error)
	Delete(context.Context, *DeletePermissionRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedPermissionServiceServer()
}

// UnimplementedPermissionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPermissionServiceServer struct {
}

func (UnimplementedPermissionServiceServer) Create(context.Context, *CreatePermissionRequest) (*CreatePermissionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedPermissionServiceServer) Find(context.Context, *FindPermissionRequest) (*FindPermissionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Find not implemented")
}
func (UnimplementedPermissionServiceServer) FindByRoles(context.Context, *FindPermissionsByRoleRequest) (*FindPermissionByRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByRoles not implemented")
}
func (UnimplementedPermissionServiceServer) FindAll(context.Context, *common.PagedElementInput) (*FindAllPermissionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindAll not implemented")
}
func (UnimplementedPermissionServiceServer) Delete(context.Context, *DeletePermissionRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedPermissionServiceServer) mustEmbedUnimplementedPermissionServiceServer() {}

// UnsafePermissionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PermissionServiceServer will
// result in compilation errors.
type UnsafePermissionServiceServer interface {
	mustEmbedUnimplementedPermissionServiceServer()
}

func RegisterPermissionServiceServer(s grpc.ServiceRegistrar, srv PermissionServiceServer) {
	s.RegisterService(&PermissionService_ServiceDesc, srv)
}

func _PermissionService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePermissionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PermissionService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServiceServer).Create(ctx, req.(*CreatePermissionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PermissionService_Find_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindPermissionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServiceServer).Find(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PermissionService_Find_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServiceServer).Find(ctx, req.(*FindPermissionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PermissionService_FindByRoles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindPermissionsByRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServiceServer).FindByRoles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PermissionService_FindByRoles_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServiceServer).FindByRoles(ctx, req.(*FindPermissionsByRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PermissionService_FindAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.PagedElementInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServiceServer).FindAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PermissionService_FindAll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServiceServer).FindAll(ctx, req.(*common.PagedElementInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _PermissionService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePermissionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PermissionService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServiceServer).Delete(ctx, req.(*DeletePermissionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PermissionService_ServiceDesc is the grpc.ServiceDesc for PermissionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PermissionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nexa.authorization.v1.PermissionService",
	HandlerType: (*PermissionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _PermissionService_Create_Handler,
		},
		{
			MethodName: "Find",
			Handler:    _PermissionService_Find_Handler,
		},
		{
			MethodName: "FindByRoles",
			Handler:    _PermissionService_FindByRoles_Handler,
		},
		{
			MethodName: "FindAll",
			Handler:    _PermissionService_FindAll_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _PermissionService_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "authorization/v1/permission.proto",
}