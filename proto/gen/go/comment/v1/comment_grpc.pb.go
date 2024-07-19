// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: comment/v1/comment.proto

package commentv1

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
	CommentService_Create_FullMethodName     = "/nexa.comment.v1.CommentService/Create"
	CommentService_Edit_FullMethodName       = "/nexa.comment.v1.CommentService/Edit"
	CommentService_Delete_FullMethodName     = "/nexa.comment.v1.CommentService/Delete"
	CommentService_GetPosts_FullMethodName   = "/nexa.comment.v1.CommentService/GetPosts"
	CommentService_GetReplies_FullMethodName = "/nexa.comment.v1.CommentService/GetReplies"
	CommentService_GetCounts_FullMethodName  = "/nexa.comment.v1.CommentService/GetCounts"
	CommentService_ClearPosts_FullMethodName = "/nexa.comment.v1.CommentService/ClearPosts"
	CommentService_ClearUsers_FullMethodName = "/nexa.comment.v1.CommentService/ClearUsers"
)

// CommentServiceClient is the client API for CommentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommentServiceClient interface {
	Create(ctx context.Context, in *CreateCommentRequest, opts ...grpc.CallOption) (*CreateCommentResponse, error)
	Edit(ctx context.Context, in *EditCommentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Delete(ctx context.Context, in *DeleteCommentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetPosts(ctx context.Context, in *GetPostCommentsRequest, opts ...grpc.CallOption) (*GetPostCommentsResponse, error)
	GetReplies(ctx context.Context, in *GetCommentRepliesRequest, opts ...grpc.CallOption) (*GetCommentRepliesResponse, error)
	GetCounts(ctx context.Context, in *GetCountsRequest, opts ...grpc.CallOption) (*GetCountsResponse, error)
	// Clear posts comments
	ClearPosts(ctx context.Context, in *ClearPostsCommentsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Clear comments created by the user
	ClearUsers(ctx context.Context, in *ClearUserCommentsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type commentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCommentServiceClient(cc grpc.ClientConnInterface) CommentServiceClient {
	return &commentServiceClient{cc}
}

func (c *commentServiceClient) Create(ctx context.Context, in *CreateCommentRequest, opts ...grpc.CallOption) (*CreateCommentResponse, error) {
	out := new(CreateCommentResponse)
	err := c.cc.Invoke(ctx, CommentService_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentServiceClient) Edit(ctx context.Context, in *EditCommentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, CommentService_Edit_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentServiceClient) Delete(ctx context.Context, in *DeleteCommentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, CommentService_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentServiceClient) GetPosts(ctx context.Context, in *GetPostCommentsRequest, opts ...grpc.CallOption) (*GetPostCommentsResponse, error) {
	out := new(GetPostCommentsResponse)
	err := c.cc.Invoke(ctx, CommentService_GetPosts_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentServiceClient) GetReplies(ctx context.Context, in *GetCommentRepliesRequest, opts ...grpc.CallOption) (*GetCommentRepliesResponse, error) {
	out := new(GetCommentRepliesResponse)
	err := c.cc.Invoke(ctx, CommentService_GetReplies_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentServiceClient) GetCounts(ctx context.Context, in *GetCountsRequest, opts ...grpc.CallOption) (*GetCountsResponse, error) {
	out := new(GetCountsResponse)
	err := c.cc.Invoke(ctx, CommentService_GetCounts_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentServiceClient) ClearPosts(ctx context.Context, in *ClearPostsCommentsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, CommentService_ClearPosts_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentServiceClient) ClearUsers(ctx context.Context, in *ClearUserCommentsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, CommentService_ClearUsers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommentServiceServer is the server API for CommentService service.
// All implementations must embed UnimplementedCommentServiceServer
// for forward compatibility
type CommentServiceServer interface {
	Create(context.Context, *CreateCommentRequest) (*CreateCommentResponse, error)
	Edit(context.Context, *EditCommentRequest) (*emptypb.Empty, error)
	Delete(context.Context, *DeleteCommentRequest) (*emptypb.Empty, error)
	GetPosts(context.Context, *GetPostCommentsRequest) (*GetPostCommentsResponse, error)
	GetReplies(context.Context, *GetCommentRepliesRequest) (*GetCommentRepliesResponse, error)
	GetCounts(context.Context, *GetCountsRequest) (*GetCountsResponse, error)
	// Clear posts comments
	ClearPosts(context.Context, *ClearPostsCommentsRequest) (*emptypb.Empty, error)
	// Clear comments created by the user
	ClearUsers(context.Context, *ClearUserCommentsRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedCommentServiceServer()
}

// UnimplementedCommentServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCommentServiceServer struct {
}

func (UnimplementedCommentServiceServer) Create(context.Context, *CreateCommentRequest) (*CreateCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedCommentServiceServer) Edit(context.Context, *EditCommentRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Edit not implemented")
}
func (UnimplementedCommentServiceServer) Delete(context.Context, *DeleteCommentRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedCommentServiceServer) GetPosts(context.Context, *GetPostCommentsRequest) (*GetPostCommentsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPosts not implemented")
}
func (UnimplementedCommentServiceServer) GetReplies(context.Context, *GetCommentRepliesRequest) (*GetCommentRepliesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReplies not implemented")
}
func (UnimplementedCommentServiceServer) GetCounts(context.Context, *GetCountsRequest) (*GetCountsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCounts not implemented")
}
func (UnimplementedCommentServiceServer) ClearPosts(context.Context, *ClearPostsCommentsRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearPosts not implemented")
}
func (UnimplementedCommentServiceServer) ClearUsers(context.Context, *ClearUserCommentsRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearUsers not implemented")
}
func (UnimplementedCommentServiceServer) mustEmbedUnimplementedCommentServiceServer() {}

// UnsafeCommentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommentServiceServer will
// result in compilation errors.
type UnsafeCommentServiceServer interface {
	mustEmbedUnimplementedCommentServiceServer()
}

func RegisterCommentServiceServer(s grpc.ServiceRegistrar, srv CommentServiceServer) {
	s.RegisterService(&CommentService_ServiceDesc, srv)
}

func _CommentService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CommentService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentServiceServer).Create(ctx, req.(*CreateCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommentService_Edit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentServiceServer).Edit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CommentService_Edit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentServiceServer).Edit(ctx, req.(*EditCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommentService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CommentService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentServiceServer).Delete(ctx, req.(*DeleteCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommentService_GetPosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPostCommentsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentServiceServer).GetPosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CommentService_GetPosts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentServiceServer).GetPosts(ctx, req.(*GetPostCommentsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommentService_GetReplies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCommentRepliesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentServiceServer).GetReplies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CommentService_GetReplies_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentServiceServer).GetReplies(ctx, req.(*GetCommentRepliesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommentService_GetCounts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCountsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentServiceServer).GetCounts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CommentService_GetCounts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentServiceServer).GetCounts(ctx, req.(*GetCountsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommentService_ClearPosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearPostsCommentsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentServiceServer).ClearPosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CommentService_ClearPosts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentServiceServer).ClearPosts(ctx, req.(*ClearPostsCommentsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CommentService_ClearUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearUserCommentsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentServiceServer).ClearUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CommentService_ClearUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentServiceServer).ClearUsers(ctx, req.(*ClearUserCommentsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CommentService_ServiceDesc is the grpc.ServiceDesc for CommentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CommentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nexa.comment.v1.CommentService",
	HandlerType: (*CommentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _CommentService_Create_Handler,
		},
		{
			MethodName: "Edit",
			Handler:    _CommentService_Edit_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _CommentService_Delete_Handler,
		},
		{
			MethodName: "GetPosts",
			Handler:    _CommentService_GetPosts_Handler,
		},
		{
			MethodName: "GetReplies",
			Handler:    _CommentService_GetReplies_Handler,
		},
		{
			MethodName: "GetCounts",
			Handler:    _CommentService_GetCounts_Handler,
		},
		{
			MethodName: "ClearPosts",
			Handler:    _CommentService_ClearPosts_Handler,
		},
		{
			MethodName: "ClearUsers",
			Handler:    _CommentService_ClearUsers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "comment/v1/comment.proto",
}
