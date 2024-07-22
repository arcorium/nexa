// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: reaction/v1/reaction.proto

package reactionv1

import (
	context "context"
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
	ReactionService_Like_FullMethodName        = "/nexa.reaction.v1.ReactionService/Like"
	ReactionService_Dislike_FullMethodName     = "/nexa.reaction.v1.ReactionService/Dislike"
	ReactionService_GetItems_FullMethodName    = "/nexa.reaction.v1.ReactionService/GetItems"
	ReactionService_GetCount_FullMethodName    = "/nexa.reaction.v1.ReactionService/GetCount"
	ReactionService_DeleteItems_FullMethodName = "/nexa.reaction.v1.ReactionService/DeleteItems"
	ReactionService_ClearUsers_FullMethodName  = "/nexa.reaction.v1.ReactionService/ClearUsers"
)

// ReactionServiceClient is the client API for ReactionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReactionServiceClient interface {
	Like(ctx context.Context, in *LikeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Dislike(ctx context.Context, in *DislikeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Get all users liked the item (post/comments)
	GetItems(ctx context.Context, in *GetItemReactionsRequest, opts ...grpc.CallOption) (*GetItemReactionsResponse, error)
	GetCount(ctx context.Context, in *GetCountRequest, opts ...grpc.CallOption) (*GetCountResponse, error)
	// Delete items reactions
	DeleteItems(ctx context.Context, in *DeleteItemsReactionsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Remove all likes created by the user
	ClearUsers(ctx context.Context, in *ClearUsersReactionsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type reactionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewReactionServiceClient(cc grpc.ClientConnInterface) ReactionServiceClient {
	return &reactionServiceClient{cc}
}

func (c *reactionServiceClient) Like(ctx context.Context, in *LikeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ReactionService_Like_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reactionServiceClient) Dislike(ctx context.Context, in *DislikeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ReactionService_Dislike_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reactionServiceClient) GetItems(ctx context.Context, in *GetItemReactionsRequest, opts ...grpc.CallOption) (*GetItemReactionsResponse, error) {
	out := new(GetItemReactionsResponse)
	err := c.cc.Invoke(ctx, ReactionService_GetItems_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reactionServiceClient) GetCount(ctx context.Context, in *GetCountRequest, opts ...grpc.CallOption) (*GetCountResponse, error) {
	out := new(GetCountResponse)
	err := c.cc.Invoke(ctx, ReactionService_GetCount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reactionServiceClient) DeleteItems(ctx context.Context, in *DeleteItemsReactionsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ReactionService_DeleteItems_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reactionServiceClient) ClearUsers(ctx context.Context, in *ClearUsersReactionsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ReactionService_ClearUsers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReactionServiceServer is the server API for ReactionService service.
// All implementations must embed UnimplementedReactionServiceServer
// for forward compatibility
type ReactionServiceServer interface {
	Like(context.Context, *LikeRequest) (*emptypb.Empty, error)
	Dislike(context.Context, *DislikeRequest) (*emptypb.Empty, error)
	// Get all users liked the item (post/comments)
	GetItems(context.Context, *GetItemReactionsRequest) (*GetItemReactionsResponse, error)
	GetCount(context.Context, *GetCountRequest) (*GetCountResponse, error)
	// Delete items reactions
	DeleteItems(context.Context, *DeleteItemsReactionsRequest) (*emptypb.Empty, error)
	// Remove all likes created by the user
	ClearUsers(context.Context, *ClearUsersReactionsRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedReactionServiceServer()
}

// UnimplementedReactionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedReactionServiceServer struct {
}

func (UnimplementedReactionServiceServer) Like(context.Context, *LikeRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Like not implemented")
}
func (UnimplementedReactionServiceServer) Dislike(context.Context, *DislikeRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Dislike not implemented")
}
func (UnimplementedReactionServiceServer) GetItems(context.Context, *GetItemReactionsRequest) (*GetItemReactionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetItems not implemented")
}
func (UnimplementedReactionServiceServer) GetCount(context.Context, *GetCountRequest) (*GetCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCount not implemented")
}
func (UnimplementedReactionServiceServer) DeleteItems(context.Context, *DeleteItemsReactionsRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteItems not implemented")
}
func (UnimplementedReactionServiceServer) ClearUsers(context.Context, *ClearUsersReactionsRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearUsers not implemented")
}
func (UnimplementedReactionServiceServer) mustEmbedUnimplementedReactionServiceServer() {}

// UnsafeReactionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReactionServiceServer will
// result in compilation errors.
type UnsafeReactionServiceServer interface {
	mustEmbedUnimplementedReactionServiceServer()
}

func RegisterReactionServiceServer(s grpc.ServiceRegistrar, srv ReactionServiceServer) {
	s.RegisterService(&ReactionService_ServiceDesc, srv)
}

func _ReactionService_Like_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LikeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReactionServiceServer).Like(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReactionService_Like_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReactionServiceServer).Like(ctx, req.(*LikeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReactionService_Dislike_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DislikeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReactionServiceServer).Dislike(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReactionService_Dislike_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReactionServiceServer).Dislike(ctx, req.(*DislikeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReactionService_GetItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetItemReactionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReactionServiceServer).GetItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReactionService_GetItems_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReactionServiceServer).GetItems(ctx, req.(*GetItemReactionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReactionService_GetCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReactionServiceServer).GetCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReactionService_GetCount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReactionServiceServer).GetCount(ctx, req.(*GetCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReactionService_DeleteItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteItemsReactionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReactionServiceServer).DeleteItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReactionService_DeleteItems_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReactionServiceServer).DeleteItems(ctx, req.(*DeleteItemsReactionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReactionService_ClearUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearUsersReactionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReactionServiceServer).ClearUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReactionService_ClearUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReactionServiceServer).ClearUsers(ctx, req.(*ClearUsersReactionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ReactionService_ServiceDesc is the grpc.ServiceDesc for ReactionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReactionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nexa.reaction.v1.ReactionService",
	HandlerType: (*ReactionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Like",
			Handler:    _ReactionService_Like_Handler,
		},
		{
			MethodName: "Dislike",
			Handler:    _ReactionService_Dislike_Handler,
		},
		{
			MethodName: "GetItems",
			Handler:    _ReactionService_GetItems_Handler,
		},
		{
			MethodName: "GetCount",
			Handler:    _ReactionService_GetCount_Handler,
		},
		{
			MethodName: "DeleteItems",
			Handler:    _ReactionService_DeleteItems_Handler,
		},
		{
			MethodName: "ClearUsers",
			Handler:    _ReactionService_ClearUsers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "reaction/v1/reaction.proto",
}
