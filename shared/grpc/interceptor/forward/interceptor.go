package forward

import (
  "context"
  middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/metadata"
  "google.golang.org/grpc/status"
)

func UnaryServerInterceptor(config *Config) grpc.UnaryServerInterceptor {
  return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok && !config.allowEmpty {
      return nil, status.Errorf(codes.InvalidArgument, "expected metadata")
    }

    newMd := config.getNewMetadata(md)
    if err := config.validate(newMd); err != nil {
      return nil, err
    }

    ctx = metadata.NewOutgoingContext(ctx, newMd)
    return handler(ctx, req)
  }
}

func StreamServerInterceptor(config *Config) grpc.StreamServerInterceptor {
  return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
    // Check if server stream is wrapped
    var currCtx context.Context
    wrappedStream, ok := ss.(*middleware.WrappedServerStream)
    if !ok {
      currCtx = ss.Context()
    } else {
      currCtx = wrappedStream.WrappedContext
    }

    md, ok := metadata.FromIncomingContext(currCtx)
    if !ok && !config.allowEmpty {
      return status.Errorf(codes.InvalidArgument, "expected metadata")
    }

    newMd := config.getNewMetadata(md)
    if err := config.validate(newMd); err != nil {
      return err
    }

    newCtx := metadata.NewOutgoingContext(currCtx, md)

    // Set new context
    wrappedServerStream := middleware.WrapServerStream(ss)
    wrappedServerStream.WrappedContext = newCtx

    return handler(srv, wrappedServerStream)
  }
}