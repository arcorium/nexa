package interceptor

import (
  "context"
  "errors"
  middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  "google.golang.org/grpc"
  "google.golang.org/grpc/health/grpc_health_v1"
)

func GetWrappedServerStream(server grpc.ServerStream) (*middleware.WrappedServerStream, error) {
  stream, ok := server.(*middleware.WrappedServerStream)
  if !ok {
    return nil, errors.New("server stream is not convertable to WrappedServerStream")
  }
  return stream, nil
}

// SkipSelector will negate the returned condition and add  HealthCheckSkipSelector to skip healthcheck endpoint
func SkipSelector(matchFunc SelectorMatchFunc) SelectorMatchFunc {
  return func(ctx context.Context, callMeta interceptors.CallMeta) bool {
    return HealthCheckSelector(ctx, callMeta) && !matchFunc(ctx, callMeta)
  }
}

func HealthCheckSkipSelector(_ context.Context, meta interceptors.CallMeta) bool {
  return meta.Service == grpc_health_v1.Health_ServiceDesc.ServiceName
}

func HealthCheckSelector(_ context.Context, meta interceptors.CallMeta) bool {
  return meta.Service != grpc_health_v1.Health_ServiceDesc.ServiceName
}
