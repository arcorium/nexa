// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: authorization/v1/role.proto

package authorizationv1

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
	RoleService_Create_FullMethodName            = "/nexa.proto.generated.authorization.v1.RoleService/Create"
	RoleService_Update_FullMethodName            = "/nexa.proto.generated.authorization.v1.RoleService/Update"
	RoleService_Delete_FullMethodName            = "/nexa.proto.generated.authorization.v1.RoleService/Delete"
	RoleService_GetUser_FullMethodName           = "/nexa.proto.generated.authorization.v1.RoleService/GetUser"
	RoleService_Find_FullMethodName              = "/nexa.proto.generated.authorization.v1.RoleService/Find"
	RoleService_FindAll_FullMethodName           = "/nexa.proto.generated.authorization.v1.RoleService/FindAll"
	RoleService_AddUser_FullMethodName           = "/nexa.proto.generated.authorization.v1.RoleService/AddUser"
	RoleService_RemoveUser_FullMethodName        = "/nexa.proto.generated.authorization.v1.RoleService/RemoveUser"
	RoleService_AppendPermissions_FullMethodName = "/nexa.proto.generated.authorization.v1.RoleService/AppendPermissions"
	RoleService_RemovePermissions_FullMethodName = "/nexa.proto.generated.authorization.v1.RoleService/RemovePermissions"
)

// RoleServiceClient is the client API for RoleService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RoleServiceClient interface {
	Create(ctx context.Context, in *RoleCreateRequest, opts ...grpc.CallOption) (*RoleCreateResponse, error)
	Update(ctx context.Context, in *RoleUpdateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Delete(ctx context.Context, in *RoleDeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetUser(ctx context.Context, in *GetUserRolesRequest, opts ...grpc.CallOption) (*GetUserRolesResponse, error)
	Find(ctx context.Context, in *RoleFindRequest, opts ...grpc.CallOption) (*RoleFindResponse, error)
	FindAll(ctx context.Context, in *common.PagedElementInput, opts ...grpc.CallOption) (*RoleFindAllResponse, error)
	AddUser(ctx context.Context, in *AddUserRolesRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	RemoveUser(ctx context.Context, in *RemoveUserRolesRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	AppendPermissions(ctx context.Context, in *RoleAppendPermissionsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	RemovePermissions(ctx context.Context, in *RoleRemovePermissionsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type roleServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRoleServiceClient(cc grpc.ClientConnInterface) RoleServiceClient {
	return &roleServiceClient{cc}
}

func (c *roleServiceClient) Create(ctx context.Context, in *RoleCreateRequest, opts ...grpc.CallOption) (*RoleCreateResponse, error) {
	out := new(RoleCreateResponse)
	err := c.cc.Invoke(ctx, RoleService_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleServiceClient) Update(ctx context.Context, in *RoleUpdateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RoleService_Update_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleServiceClient) Delete(ctx context.Context, in *RoleDeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RoleService_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleServiceClient) GetUser(ctx context.Context, in *GetUserRolesRequest, opts ...grpc.CallOption) (*GetUserRolesResponse, error) {
	out := new(GetUserRolesResponse)
	err := c.cc.Invoke(ctx, RoleService_GetUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleServiceClient) Find(ctx context.Context, in *RoleFindRequest, opts ...grpc.CallOption) (*RoleFindResponse, error) {
	out := new(RoleFindResponse)
	err := c.cc.Invoke(ctx, RoleService_Find_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleServiceClient) FindAll(ctx context.Context, in *common.PagedElementInput, opts ...grpc.CallOption) (*RoleFindAllResponse, error) {
	out := new(RoleFindAllResponse)
	err := c.cc.Invoke(ctx, RoleService_FindAll_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleServiceClient) AddUser(ctx context.Context, in *AddUserRolesRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RoleService_AddUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleServiceClient) RemoveUser(ctx context.Context, in *RemoveUserRolesRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RoleService_RemoveUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleServiceClient) AppendPermissions(ctx context.Context, in *RoleAppendPermissionsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RoleService_AppendPermissions_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleServiceClient) RemovePermissions(ctx context.Context, in *RoleRemovePermissionsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, RoleService_RemovePermissions_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoleServiceServer is the server API for RoleService service.
// All implementations must embed UnimplementedRoleServiceServer
// for forward compatibility
type RoleServiceServer interface {
	Create(context.Context, *RoleCreateRequest) (*RoleCreateResponse, error)
	Update(context.Context, *RoleUpdateRequest) (*emptypb.Empty, error)
	Delete(context.Context, *RoleDeleteRequest) (*emptypb.Empty, error)
	GetUser(context.Context, *GetUserRolesRequest) (*GetUserRolesResponse, error)
	Find(context.Context, *RoleFindRequest) (*RoleFindResponse, error)
	FindAll(context.Context, *common.PagedElementInput) (*RoleFindAllResponse, error)
	AddUser(context.Context, *AddUserRolesRequest) (*emptypb.Empty, error)
	RemoveUser(context.Context, *RemoveUserRolesRequest) (*emptypb.Empty, error)
	AppendPermissions(context.Context, *RoleAppendPermissionsRequest) (*emptypb.Empty, error)
	RemovePermissions(context.Context, *RoleRemovePermissionsRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedRoleServiceServer()
}

// UnimplementedRoleServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRoleServiceServer struct {
}

func (UnimplementedRoleServiceServer) Create(context.Context, *RoleCreateRequest) (*RoleCreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedRoleServiceServer) Update(context.Context, *RoleUpdateRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedRoleServiceServer) Delete(context.Context, *RoleDeleteRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedRoleServiceServer) GetUser(context.Context, *GetUserRolesRequest) (*GetUserRolesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedRoleServiceServer) Find(context.Context, *RoleFindRequest) (*RoleFindResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Find not implemented")
}
func (UnimplementedRoleServiceServer) FindAll(context.Context, *common.PagedElementInput) (*RoleFindAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindAll not implemented")
}
func (UnimplementedRoleServiceServer) AddUser(context.Context, *AddUserRolesRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUser not implemented")
}
func (UnimplementedRoleServiceServer) RemoveUser(context.Context, *RemoveUserRolesRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveUser not implemented")
}
func (UnimplementedRoleServiceServer) AppendPermissions(context.Context, *RoleAppendPermissionsRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendPermissions not implemented")
}
func (UnimplementedRoleServiceServer) RemovePermissions(context.Context, *RoleRemovePermissionsRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemovePermissions not implemented")
}
func (UnimplementedRoleServiceServer) mustEmbedUnimplementedRoleServiceServer() {}

// UnsafeRoleServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RoleServiceServer will
// result in compilation errors.
type UnsafeRoleServiceServer interface {
	mustEmbedUnimplementedRoleServiceServer()
}

func RegisterRoleServiceServer(s grpc.ServiceRegistrar, srv RoleServiceServer) {
	s.RegisterService(&RoleService_ServiceDesc, srv)
}

func _RoleService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoleCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleServiceServer).Create(ctx, req.(*RoleCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoleUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleService_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleServiceServer).Update(ctx, req.(*RoleUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoleDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleServiceServer).Delete(ctx, req.(*RoleDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleService_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRolesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleServiceServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleService_GetUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleServiceServer).GetUser(ctx, req.(*GetUserRolesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleService_Find_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoleFindRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleServiceServer).Find(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleService_Find_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleServiceServer).Find(ctx, req.(*RoleFindRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleService_FindAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.PagedElementInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleServiceServer).FindAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleService_FindAll_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleServiceServer).FindAll(ctx, req.(*common.PagedElementInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleService_AddUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddUserRolesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleServiceServer).AddUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleService_AddUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleServiceServer).AddUser(ctx, req.(*AddUserRolesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleService_RemoveUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveUserRolesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleServiceServer).RemoveUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleService_RemoveUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleServiceServer).RemoveUser(ctx, req.(*RemoveUserRolesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleService_AppendPermissions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoleAppendPermissionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleServiceServer).AppendPermissions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleService_AppendPermissions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleServiceServer).AppendPermissions(ctx, req.(*RoleAppendPermissionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoleService_RemovePermissions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoleRemovePermissionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoleServiceServer).RemovePermissions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoleService_RemovePermissions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoleServiceServer).RemovePermissions(ctx, req.(*RoleRemovePermissionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RoleService_ServiceDesc is the grpc.ServiceDesc for RoleService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RoleService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nexa.proto.generated.authorization.v1.RoleService",
	HandlerType: (*RoleServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _RoleService_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _RoleService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _RoleService_Delete_Handler,
		},
		{
			MethodName: "GetUser",
			Handler:    _RoleService_GetUser_Handler,
		},
		{
			MethodName: "Find",
			Handler:    _RoleService_Find_Handler,
		},
		{
			MethodName: "FindAll",
			Handler:    _RoleService_FindAll_Handler,
		},
		{
			MethodName: "AddUser",
			Handler:    _RoleService_AddUser_Handler,
		},
		{
			MethodName: "RemoveUser",
			Handler:    _RoleService_RemoveUser_Handler,
		},
		{
			MethodName: "AppendPermissions",
			Handler:    _RoleService_AppendPermissions_Handler,
		},
		{
			MethodName: "RemovePermissions",
			Handler:    _RoleService_RemovePermissions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "authorization/v1/role.proto",
}
