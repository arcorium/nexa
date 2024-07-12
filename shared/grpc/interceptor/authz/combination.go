package authz

import (
  "context"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  "google.golang.org/grpc"
  "slices"
)

// UnaryServerCombination create unary server interceptor for authorization both for user and protected API.
// It is used for minimalize checking redundancy with bypassing the user authorization check when it is
// already handled by protected API authorization. Authorization selector only will be called if only
// the request(rpc) is not handled by protected API authorization.
func UnaryServerCombination(conf *CombinationConfig) grpc.UnaryServerInterceptor {
  return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
    meta := interceptors.NewServerCallMeta(info.FullMethod, nil, req)
    if slices.Contains(conf.SkipServices, meta.Service) {
      return handler(ctx, req)
    }

    result := conf.Selector(ctx, meta)
    switch result {
    case Private:
      return privateUnaryAuthorization(&conf.Private)(ctx, req, info, handler)
    case UserAuth:
      return userUnaryAuthorization(&conf.User)(ctx, req, info, handler)
    default:
      return handler(ctx, req)
    }
  }
}

// StreamServerCombination works the same as UnaryServerCombination, but it works for stream.
func StreamServerCombination(conf *CombinationConfig) grpc.StreamServerInterceptor {
  return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
    meta := interceptors.NewServerCallMeta(info.FullMethod, info, nil)
    if slices.Contains(conf.SkipServices, meta.Service) {
      return handler(srv, ss)
    }

    result := conf.Selector(ss.Context(), meta)
    switch result {
    case Private:
      return privateStreamAuthorization(&conf.Private)(srv, ss, info, handler)
    case UserAuth:
      return userStreamAuthorization(&conf.User)(srv, ss, info, handler)
    default:
      return handler(srv, ss)
    }
  }
}
