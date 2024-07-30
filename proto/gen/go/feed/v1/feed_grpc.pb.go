// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: feed/v1/feed.proto

package feedv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	FeedService_GetUserFeed_FullMethodName = "/nexa.feed.v1.FeedService/GetUserFeed"
)

// FeedServiceClient is the client API for FeedService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FeedServiceClient interface {
	GetUserFeed(ctx context.Context, in *GetUserFeedRequest, opts ...grpc.CallOption) (*GetUserFeedResponse, error)
}

type feedServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFeedServiceClient(cc grpc.ClientConnInterface) FeedServiceClient {
	return &feedServiceClient{cc}
}

func (c *feedServiceClient) GetUserFeed(ctx context.Context, in *GetUserFeedRequest, opts ...grpc.CallOption) (*GetUserFeedResponse, error) {
	out := new(GetUserFeedResponse)
	err := c.cc.Invoke(ctx, FeedService_GetUserFeed_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FeedServiceServer is the server API for FeedService service.
// All implementations must embed UnimplementedFeedServiceServer
// for forward compatibility
type FeedServiceServer interface {
	GetUserFeed(context.Context, *GetUserFeedRequest) (*GetUserFeedResponse, error)
	mustEmbedUnimplementedFeedServiceServer()
}

// UnimplementedFeedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFeedServiceServer struct {
}

func (UnimplementedFeedServiceServer) GetUserFeed(context.Context, *GetUserFeedRequest) (*GetUserFeedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserFeed not implemented")
}
func (UnimplementedFeedServiceServer) mustEmbedUnimplementedFeedServiceServer() {}

// UnsafeFeedServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FeedServiceServer will
// result in compilation errors.
type UnsafeFeedServiceServer interface {
	mustEmbedUnimplementedFeedServiceServer()
}

func RegisterFeedServiceServer(s grpc.ServiceRegistrar, srv FeedServiceServer) {
	s.RegisterService(&FeedService_ServiceDesc, srv)
}

func _FeedService_GetUserFeed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserFeedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedServiceServer).GetUserFeed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FeedService_GetUserFeed_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedServiceServer).GetUserFeed(ctx, req.(*GetUserFeedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FeedService_ServiceDesc is the grpc.ServiceDesc for FeedService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FeedService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nexa.feed.v1.FeedService",
	HandlerType: (*FeedServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserFeed",
			Handler:    _FeedService_GetUserFeed_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "feed/v1/feed.proto",
}