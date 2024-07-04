package interceptor

import (
  "context"
  "errors"
  middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  "google.golang.org/grpc"
)

func GetWrappedServerStream(server grpc.ServerStream) (*middleware.WrappedServerStream, error) {
  stream, ok := server.(*middleware.WrappedServerStream)
  if !ok {
    return nil, errors.New("server stream is not convertable to WrappedServerStream")
  }
  return stream, nil
}

// SkipSelector will negate the returned condition
func SkipSelector(matchFunc SelectorMatchFunc) SelectorMatchFunc {
  return func(ctx context.Context, callMeta interceptors.CallMeta) bool {
    return !matchFunc(ctx, callMeta)
  }
}
