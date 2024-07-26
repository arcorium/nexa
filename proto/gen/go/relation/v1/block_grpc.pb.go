// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: relation/v1/block.proto

package relationv1

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
	BlockService_Block_FullMethodName         = "/nexa.relation.v1.BlockService/Block"
	BlockService_Unblock_FullMethodName       = "/nexa.relation.v1.BlockService/Unblock"
	BlockService_IsBlocked_FullMethodName     = "/nexa.relation.v1.BlockService/IsBlocked"
	BlockService_GetBlocked_FullMethodName    = "/nexa.relation.v1.BlockService/GetBlocked"
	BlockService_GetUsersCount_FullMethodName = "/nexa.relation.v1.BlockService/GetUsersCount"
	BlockService_ClearUsers_FullMethodName    = "/nexa.relation.v1.BlockService/ClearUsers"
)

// BlockServiceClient is the client API for BlockService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BlockServiceClient interface {
	// Get user total follower and followee
	Block(ctx context.Context, in *BlockUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Unblock(ctx context.Context, in *UnblockUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	IsBlocked(ctx context.Context, in *IsUserBlockedRequest, opts ...grpc.CallOption) (*IsUserBlockedResponse, error)
	GetBlocked(ctx context.Context, in *common.PagedElementInput, opts ...grpc.CallOption) (*GetBlockedResponse, error)
	// Get user total blocked
	GetUsersCount(ctx context.Context, in *GetUsersBlockCountRequest, opts ...grpc.CallOption) (*GetUsersBlockCountResponse, error)
	// Clear user related data on this service
	ClearUsers(ctx context.Context, in *ClearUsersBlockRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type blockServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBlockServiceClient(cc grpc.ClientConnInterface) BlockServiceClient {
	return &blockServiceClient{cc}
}

func (c *blockServiceClient) Block(ctx context.Context, in *BlockUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, BlockService_Block_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockServiceClient) Unblock(ctx context.Context, in *UnblockUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, BlockService_Unblock_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockServiceClient) IsBlocked(ctx context.Context, in *IsUserBlockedRequest, opts ...grpc.CallOption) (*IsUserBlockedResponse, error) {
	out := new(IsUserBlockedResponse)
	err := c.cc.Invoke(ctx, BlockService_IsBlocked_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockServiceClient) GetBlocked(ctx context.Context, in *common.PagedElementInput, opts ...grpc.CallOption) (*GetBlockedResponse, error) {
	out := new(GetBlockedResponse)
	err := c.cc.Invoke(ctx, BlockService_GetBlocked_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockServiceClient) GetUsersCount(ctx context.Context, in *GetUsersBlockCountRequest, opts ...grpc.CallOption) (*GetUsersBlockCountResponse, error) {
	out := new(GetUsersBlockCountResponse)
	err := c.cc.Invoke(ctx, BlockService_GetUsersCount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockServiceClient) ClearUsers(ctx context.Context, in *ClearUsersBlockRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, BlockService_ClearUsers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BlockServiceServer is the server API for BlockService service.
// All implementations must embed UnimplementedBlockServiceServer
// for forward compatibility
type BlockServiceServer interface {
	// Get user total follower and followee
	Block(context.Context, *BlockUserRequest) (*emptypb.Empty, error)
	Unblock(context.Context, *UnblockUserRequest) (*emptypb.Empty, error)
	IsBlocked(context.Context, *IsUserBlockedRequest) (*IsUserBlockedResponse, error)
	GetBlocked(context.Context, *common.PagedElementInput) (*GetBlockedResponse, error)
	// Get user total blocked
	GetUsersCount(context.Context, *GetUsersBlockCountRequest) (*GetUsersBlockCountResponse, error)
	// Clear user related data on this service
	ClearUsers(context.Context, *ClearUsersBlockRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedBlockServiceServer()
}

// UnimplementedBlockServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBlockServiceServer struct {
}

func (UnimplementedBlockServiceServer) Block(context.Context, *BlockUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Block not implemented")
}
func (UnimplementedBlockServiceServer) Unblock(context.Context, *UnblockUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unblock not implemented")
}
func (UnimplementedBlockServiceServer) IsBlocked(context.Context, *IsUserBlockedRequest) (*IsUserBlockedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsBlocked not implemented")
}
func (UnimplementedBlockServiceServer) GetBlocked(context.Context, *common.PagedElementInput) (*GetBlockedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlocked not implemented")
}
func (UnimplementedBlockServiceServer) GetUsersCount(context.Context, *GetUsersBlockCountRequest) (*GetUsersBlockCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUsersCount not implemented")
}
func (UnimplementedBlockServiceServer) ClearUsers(context.Context, *ClearUsersBlockRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearUsers not implemented")
}
func (UnimplementedBlockServiceServer) mustEmbedUnimplementedBlockServiceServer() {}

// UnsafeBlockServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BlockServiceServer will
// result in compilation errors.
type UnsafeBlockServiceServer interface {
	mustEmbedUnimplementedBlockServiceServer()
}

func RegisterBlockServiceServer(s grpc.ServiceRegistrar, srv BlockServiceServer) {
	s.RegisterService(&BlockService_ServiceDesc, srv)
}

func _BlockService_Block_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BlockUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockServiceServer).Block(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BlockService_Block_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockServiceServer).Block(ctx, req.(*BlockUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockService_Unblock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnblockUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockServiceServer).Unblock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BlockService_Unblock_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockServiceServer).Unblock(ctx, req.(*UnblockUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockService_IsBlocked_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsUserBlockedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockServiceServer).IsBlocked(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BlockService_IsBlocked_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockServiceServer).IsBlocked(ctx, req.(*IsUserBlockedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockService_GetBlocked_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.PagedElementInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockServiceServer).GetBlocked(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BlockService_GetBlocked_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockServiceServer).GetBlocked(ctx, req.(*common.PagedElementInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockService_GetUsersCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsersBlockCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockServiceServer).GetUsersCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BlockService_GetUsersCount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockServiceServer).GetUsersCount(ctx, req.(*GetUsersBlockCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockService_ClearUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearUsersBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockServiceServer).ClearUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BlockService_ClearUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockServiceServer).ClearUsers(ctx, req.(*ClearUsersBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BlockService_ServiceDesc is the grpc.ServiceDesc for BlockService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BlockService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nexa.relation.v1.BlockService",
	HandlerType: (*BlockServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Block",
			Handler:    _BlockService_Block_Handler,
		},
		{
			MethodName: "Unblock",
			Handler:    _BlockService_Unblock_Handler,
		},
		{
			MethodName: "IsBlocked",
			Handler:    _BlockService_IsBlocked_Handler,
		},
		{
			MethodName: "GetBlocked",
			Handler:    _BlockService_GetBlocked_Handler,
		},
		{
			MethodName: "GetUsersCount",
			Handler:    _BlockService_GetUsersCount_Handler,
		},
		{
			MethodName: "ClearUsers",
			Handler:    _BlockService_ClearUsers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "relation/v1/block.proto",
}
