package authz

import (
  "context"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/variadic"
  "github.com/golang-jwt/jwt/v5"
  middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
)

// NOTE: Differentiation of authorization function is due to the jwt API doesn't support generic as the claims type

// userUnaryAuthorization extract user claims from context and check the permission based on argument. If the permission check function is nil
// it will bypass it and only do the claim extraction.
func userUnaryAuthorization(conf *UserConfig) grpc.UnaryServerInterceptor {
  return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    // extractClaims claims from token in context
    claims := sharedJwt.UserClaims{}
    err := conf.extractClaims(ctx, &claims)
    if err != nil {
      return nil, status.Error(codes.Unauthenticated, err.Error())
    }

    // Check permission
    meta := interceptors.NewServerCallMeta(info.FullMethod, nil, req)
    if conf.CheckFunc != nil && !conf.CheckFunc(&claims, meta) {
      return nil, status.New(codes.PermissionDenied, sharedErr.ErrUnauthorized.Error()).Err()
    }

    ctx = context.WithValue(ctx, conf.ClaimsKey, &claims)
    return handler(ctx, req)
  }
}

func privateUnaryAuthorization(conf *PrivateConfig) grpc.UnaryServerInterceptor {
  return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    // extractClaims claims from token in context
    claims := sharedJwt.PrivateClaims{}
    err := conf.extractClaims(ctx, &claims)
    if err != nil {
      return nil, status.Error(codes.Unauthenticated, err.Error())
    }

    // Check permission
    meta := interceptors.NewServerCallMeta(info.FullMethod, nil, req)
    if conf.CheckFunc != nil && !conf.CheckFunc(&claims, meta) {
      return nil, status.New(codes.PermissionDenied, sharedErr.ErrUnauthorized.Error()).Err()
    }

    ctx = context.WithValue(ctx, conf.ClaimsKey, &claims)
    return handler(ctx, req)
  }
}

// userStreamAuthorization works like unaryAuthorization but for stream
func userStreamAuthorization(conf *UserConfig) grpc.StreamServerInterceptor {
  return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
    var err error
    var newCtx context.Context

    // Check if server stream is wrapped
    var currCtx context.Context
    wrappedStream, ok := ss.(*middleware.WrappedServerStream)
    if !ok {
      currCtx = ss.Context()
    } else {
      currCtx = wrappedStream.WrappedContext
    }

    claims := sharedJwt.UserClaims{}
    err = conf.extractClaims(currCtx, &claims)
    if err != nil {
      return status.Error(codes.Unauthenticated, err.Error())
    }

    // Check permission
    meta := interceptors.NewServerCallMeta(info.FullMethod, info, nil)
    if conf.CheckFunc != nil && !conf.CheckFunc(&claims, meta) {
      return status.New(codes.PermissionDenied, sharedErr.ErrUnauthorized.Error()).Err()
    }

    // Set new context
    newCtx = context.WithValue(currCtx, conf.ClaimsKey, &claims)
    wrappedServerStream := middleware.WrapServerStream(ss)
    wrappedServerStream.WrappedContext = newCtx

    return handler(srv, wrappedServerStream)
  }
}

func privateStreamAuthorization(conf *PrivateConfig) grpc.StreamServerInterceptor {
  return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
    var err error
    var newCtx context.Context

    // Check if server stream is wrapped
    var currCtx context.Context
    wrappedStream, ok := ss.(*middleware.WrappedServerStream)
    if !ok {
      currCtx = ss.Context()
    } else {
      currCtx = wrappedStream.WrappedContext
    }

    claims := sharedJwt.PrivateClaims{}
    err = conf.extractClaims(currCtx, &claims)
    if err != nil {
      return status.Error(codes.Unauthenticated, err.Error())
    }

    // Check permission
    meta := interceptors.NewServerCallMeta(info.FullMethod, info, nil)
    if conf.CheckFunc != nil && !conf.CheckFunc(&claims, meta) {
      return status.New(codes.PermissionDenied, sharedErr.ErrUnauthorized.Error()).Err()
    }

    // Set new context
    newCtx = context.WithValue(currCtx, conf.ClaimsKey, &claims)
    wrappedServerStream := middleware.WrapServerStream(ss)
    wrappedServerStream.WrappedContext = newCtx

    return handler(srv, wrappedServerStream)
  }
}

type PermCheckFunc[T jwt.Claims] func(claims *T, meta interceptors.CallMeta) bool

type SelectorMatchFunc func(ctx context.Context, callMeta interceptors.CallMeta) bool

// UserUnaryServerInterceptor create unary server interceptor for authorization which will extractClaims token into claims
// and doing the permission check based on the config. Selector used to determine if the request should be forwarded
// to authorization or not.
func UserUnaryServerInterceptor(config *UserConfig, sf SelectorMatchFunc, skipServices ...string) grpc.UnaryServerInterceptor {
  if variadic.New(skipServices...).HasAtLeast(1) {
    return selector.UnaryServerInterceptor(
      userUnaryAuthorization(config),
      &SkipServiceMatcher{
        SkipServices: skipServices,
        Chain:        sf,
      },
    )
  }
  return selector.UnaryServerInterceptor(
    userUnaryAuthorization(config),
    selector.MatchFunc(sf),
  )
}

// PrivateUnaryServerInterceptor works the same as UserUnaryServerInterceptor, but it using private claims
func PrivateUnaryServerInterceptor(config *PrivateConfig, sf SelectorMatchFunc, skipServices ...string) grpc.UnaryServerInterceptor {
  if variadic.New(skipServices...).HasAtLeast(1) {
    return selector.UnaryServerInterceptor(
      privateUnaryAuthorization(config),
      &SkipServiceMatcher{
        SkipServices: skipServices,
        Chain:        sf,
      },
    )
  }
  return selector.UnaryServerInterceptor(
    privateUnaryAuthorization(config),
    selector.MatchFunc(sf),
  )
}

// UserStreamServerInterceptor works the same as UserUnaryServerInterceptor, but it works for stream.
func UserStreamServerInterceptor(config *UserConfig, sf SelectorMatchFunc, skipServices ...string) grpc.StreamServerInterceptor {
  if variadic.New(skipServices...).HasAtLeast(1) {
    return selector.StreamServerInterceptor(
      userStreamAuthorization(config),
      &SkipServiceMatcher{
        SkipServices: skipServices,
        Chain:        sf,
      },
    )
  }
  return selector.StreamServerInterceptor(
    userStreamAuthorization(config),
    selector.MatchFunc(sf),
  )
}

// PrivateStreamServerInterceptor works the same as PrivateUnaryServerInterceptor, but it works for stream
func PrivateStreamServerInterceptor(config *PrivateConfig, sf SelectorMatchFunc, skipServices ...string) grpc.StreamServerInterceptor {
  if variadic.New(skipServices...).HasAtLeast(1) {
    return selector.StreamServerInterceptor(
      privateStreamAuthorization(config),
      &SkipServiceMatcher{
        SkipServices: skipServices,
        Chain:        sf,
      },
    )
  }
  return selector.StreamServerInterceptor(
    privateStreamAuthorization(config),
    selector.MatchFunc(sf),
  )
}
