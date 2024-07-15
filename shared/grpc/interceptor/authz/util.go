package authz

import (
  "context"
  "errors"
  middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  "google.golang.org/grpc"
  "slices"
)

func GetWrappedServerStream(server grpc.ServerStream) (*middleware.WrappedServerStream, error) {
  stream, ok := server.(*middleware.WrappedServerStream)
  if !ok {
    return nil, errors.New("server stream is not convertable to WrappedServerStream")
  }
  return stream, nil
}

func GetWrappedContext(server grpc.ServerStream) context.Context {
  ctx := server.Context()
  srv, err := GetWrappedServerStream(server)
  if err != nil {
    return ctx
  }
  return srv.WrappedContext
}

// SkipSelector will negate the returned condition and add  HealthCheckSkipSelector to skip healthcheck endpoint
func SkipSelector(matchFunc SelectorMatchFunc) SelectorMatchFunc {
  return func(ctx context.Context, callMeta interceptors.CallMeta) bool {
    return !matchFunc(ctx, callMeta)
  }
}

type SkipServiceMatcher struct {
  SkipServices []string
  Chain        SelectorMatchFunc
}

func (s *SkipServiceMatcher) Match(ctx context.Context, callMeta interceptors.CallMeta) bool {
  if slices.Contains(s.SkipServices, callMeta.Service) {
    return false
  }
  return s.Chain(ctx, callMeta)
}
