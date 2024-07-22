// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: post/v1/post.proto

package postv1

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
	PostService_Find_FullMethodName             = "/nexa.post.v1.PostService/Find"
	PostService_FindEdited_FullMethodName       = "/nexa.post.v1.PostService/FindEdited"
	PostService_FindById_FullMethodName         = "/nexa.post.v1.PostService/FindById"
	PostService_FindUsers_FullMethodName        = "/nexa.post.v1.PostService/FindUsers"
	PostService_Create_FullMethodName           = "/nexa.post.v1.PostService/Create"
	PostService_UpdateVisibility_FullMethodName = "/nexa.post.v1.PostService/UpdateVisibility"
	PostService_Edit_FullMethodName             = "/nexa.post.v1.PostService/Edit"
	PostService_Delete_FullMethodName           = "/nexa.post.v1.PostService/Delete"
	PostService_Bookmark_FullMethodName         = "/nexa.post.v1.PostService/Bookmark"
	PostService_GetBookmarked_FullMethodName    = "/nexa.post.v1.PostService/GetBookmarked"
	PostService_ClearUsers_FullMethodName       = "/nexa.post.v1.PostService/ClearUsers"
)

// PostServiceClient is the client API for PostService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PostServiceClient interface {
	Find(ctx context.Context, in *common.PagedElementInput, opts ...grpc.CallOption) (*FindPostResponse, error)
	FindEdited(ctx context.Context, in *FindEditedPostRequest, opts ...grpc.CallOption) (*FindEditedPostResponse, error)
	FindById(ctx context.Context, in *FindPostByIdRequest, opts ...grpc.CallOption) (*FindPostByIdResponse, error)
	FindUsers(ctx context.Context, in *FindUserPostRequest, opts ...grpc.CallOption) (*FindUserPostResponse, error)
	Create(ctx context.Context, in *CreatePostRequest, opts ...grpc.CallOption) (*CreatePostResponse, error)
	UpdateVisibility(ctx context.Context, in *UpdatePostVisibilityRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// rpc Edit(EditPostRequest) returns (EditPostResponse);
	Edit(ctx context.Context, in *EditPostRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Delete(ctx context.Context, in *DeletePostRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Bookmark(ctx context.Context, in *BookmarkPostRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetBookmarked(ctx context.Context, in *GetBookmarkedPostRequest, opts ...grpc.CallOption) (*GetBookmarkedPostResponse, error)
	ClearUsers(ctx context.Context, in *ClearUserPostsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type postServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPostServiceClient(cc grpc.ClientConnInterface) PostServiceClient {
	return &postServiceClient{cc}
}

func (c *postServiceClient) Find(ctx context.Context, in *common.PagedElementInput, opts ...grpc.CallOption) (*FindPostResponse, error) {
	out := new(FindPostResponse)
	err := c.cc.Invoke(ctx, PostService_Find_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) FindEdited(ctx context.Context, in *FindEditedPostRequest, opts ...grpc.CallOption) (*FindEditedPostResponse, error) {
	out := new(FindEditedPostResponse)
	err := c.cc.Invoke(ctx, PostService_FindEdited_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) FindById(ctx context.Context, in *FindPostByIdRequest, opts ...grpc.CallOption) (*FindPostByIdResponse, error) {
	out := new(FindPostByIdResponse)
	err := c.cc.Invoke(ctx, PostService_FindById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) FindUsers(ctx context.Context, in *FindUserPostRequest, opts ...grpc.CallOption) (*FindUserPostResponse, error) {
	out := new(FindUserPostResponse)
	err := c.cc.Invoke(ctx, PostService_FindUsers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) Create(ctx context.Context, in *CreatePostRequest, opts ...grpc.CallOption) (*CreatePostResponse, error) {
	out := new(CreatePostResponse)
	err := c.cc.Invoke(ctx, PostService_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) UpdateVisibility(ctx context.Context, in *UpdatePostVisibilityRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PostService_UpdateVisibility_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) Edit(ctx context.Context, in *EditPostRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PostService_Edit_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) Delete(ctx context.Context, in *DeletePostRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PostService_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) Bookmark(ctx context.Context, in *BookmarkPostRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PostService_Bookmark_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) GetBookmarked(ctx context.Context, in *GetBookmarkedPostRequest, opts ...grpc.CallOption) (*GetBookmarkedPostResponse, error) {
	out := new(GetBookmarkedPostResponse)
	err := c.cc.Invoke(ctx, PostService_GetBookmarked_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *postServiceClient) ClearUsers(ctx context.Context, in *ClearUserPostsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, PostService_ClearUsers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PostServiceServer is the server API for PostService service.
// All implementations must embed UnimplementedPostServiceServer
// for forward compatibility
type PostServiceServer interface {
	Find(context.Context, *common.PagedElementInput) (*FindPostResponse, error)
	FindEdited(context.Context, *FindEditedPostRequest) (*FindEditedPostResponse, error)
	FindById(context.Context, *FindPostByIdRequest) (*FindPostByIdResponse, error)
	FindUsers(context.Context, *FindUserPostRequest) (*FindUserPostResponse, error)
	Create(context.Context, *CreatePostRequest) (*CreatePostResponse, error)
	UpdateVisibility(context.Context, *UpdatePostVisibilityRequest) (*emptypb.Empty, error)
	// rpc Edit(EditPostRequest) returns (EditPostResponse);
	Edit(context.Context, *EditPostRequest) (*emptypb.Empty, error)
	Delete(context.Context, *DeletePostRequest) (*emptypb.Empty, error)
	Bookmark(context.Context, *BookmarkPostRequest) (*emptypb.Empty, error)
	GetBookmarked(context.Context, *GetBookmarkedPostRequest) (*GetBookmarkedPostResponse, error)
	ClearUsers(context.Context, *ClearUserPostsRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedPostServiceServer()
}

// UnimplementedPostServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPostServiceServer struct {
}

func (UnimplementedPostServiceServer) Find(context.Context, *common.PagedElementInput) (*FindPostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Find not implemented")
}
func (UnimplementedPostServiceServer) FindEdited(context.Context, *FindEditedPostRequest) (*FindEditedPostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindEdited not implemented")
}
func (UnimplementedPostServiceServer) FindById(context.Context, *FindPostByIdRequest) (*FindPostByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindById not implemented")
}
func (UnimplementedPostServiceServer) FindUsers(context.Context, *FindUserPostRequest) (*FindUserPostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindUsers not implemented")
}
func (UnimplementedPostServiceServer) Create(context.Context, *CreatePostRequest) (*CreatePostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedPostServiceServer) UpdateVisibility(context.Context, *UpdatePostVisibilityRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateVisibility not implemented")
}
func (UnimplementedPostServiceServer) Edit(context.Context, *EditPostRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Edit not implemented")
}
func (UnimplementedPostServiceServer) Delete(context.Context, *DeletePostRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedPostServiceServer) Bookmark(context.Context, *BookmarkPostRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Bookmark not implemented")
}
func (UnimplementedPostServiceServer) GetBookmarked(context.Context, *GetBookmarkedPostRequest) (*GetBookmarkedPostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBookmarked not implemented")
}
func (UnimplementedPostServiceServer) ClearUsers(context.Context, *ClearUserPostsRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearUsers not implemented")
}
func (UnimplementedPostServiceServer) mustEmbedUnimplementedPostServiceServer() {}

// UnsafePostServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PostServiceServer will
// result in compilation errors.
type UnsafePostServiceServer interface {
	mustEmbedUnimplementedPostServiceServer()
}

func RegisterPostServiceServer(s grpc.ServiceRegistrar, srv PostServiceServer) {
	s.RegisterService(&PostService_ServiceDesc, srv)
}

func _PostService_Find_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.PagedElementInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).Find(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_Find_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).Find(ctx, req.(*common.PagedElementInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_FindEdited_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindEditedPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).FindEdited(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_FindEdited_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).FindEdited(ctx, req.(*FindEditedPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_FindById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindPostByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).FindById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_FindById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).FindById(ctx, req.(*FindPostByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_FindUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindUserPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).FindUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_FindUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).FindUsers(ctx, req.(*FindUserPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).Create(ctx, req.(*CreatePostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_UpdateVisibility_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePostVisibilityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).UpdateVisibility(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_UpdateVisibility_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).UpdateVisibility(ctx, req.(*UpdatePostVisibilityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_Edit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).Edit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_Edit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).Edit(ctx, req.(*EditPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).Delete(ctx, req.(*DeletePostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_Bookmark_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BookmarkPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).Bookmark(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_Bookmark_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).Bookmark(ctx, req.(*BookmarkPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_GetBookmarked_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBookmarkedPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).GetBookmarked(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_GetBookmarked_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).GetBookmarked(ctx, req.(*GetBookmarkedPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PostService_ClearUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearUserPostsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PostServiceServer).ClearUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PostService_ClearUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PostServiceServer).ClearUsers(ctx, req.(*ClearUserPostsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PostService_ServiceDesc is the grpc.ServiceDesc for PostService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PostService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nexa.post.v1.PostService",
	HandlerType: (*PostServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Find",
			Handler:    _PostService_Find_Handler,
		},
		{
			MethodName: "FindEdited",
			Handler:    _PostService_FindEdited_Handler,
		},
		{
			MethodName: "FindById",
			Handler:    _PostService_FindById_Handler,
		},
		{
			MethodName: "FindUsers",
			Handler:    _PostService_FindUsers_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _PostService_Create_Handler,
		},
		{
			MethodName: "UpdateVisibility",
			Handler:    _PostService_UpdateVisibility_Handler,
		},
		{
			MethodName: "Edit",
			Handler:    _PostService_Edit_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _PostService_Delete_Handler,
		},
		{
			MethodName: "Bookmark",
			Handler:    _PostService_Bookmark_Handler,
		},
		{
			MethodName: "GetBookmarked",
			Handler:    _PostService_GetBookmarked_Handler,
		},
		{
			MethodName: "ClearUsers",
			Handler:    _PostService_ClearUsers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "post/v1/post.proto",
}