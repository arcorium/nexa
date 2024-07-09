package interceptor

import (
  "context"
  "errors"
  "fmt"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/golang-jwt/jwt/v5"
  middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/metadata"
  "google.golang.org/grpc/status"
  "strings"
)

//func getDataFromMD(md metadata.MD) (*AuthorizationConfig, error) {
//  tokenStr := md.Get(constant.TOKEN_METADATA_KEY)
//  if len(tokenStr) != 1 {
//    return nil, errors.New("token metadata is malformed")
//  }
//  scheme := md.Get(constant.TOKEN_SCHEME_METADATA_KEY)
//  if len(scheme) != 1 {
//    return nil, errors.New("token type is malformed")
//  }
//  secret := md.Get(constant.JWT_SECRET_METADATA_KEY)
//  if len(secret) != 1 {
//    return nil, errors.New("secret key is malformed")
//  }
//  signingMethod := md.Get(constant.JWT_SIGNING_METHOD_METADATA_KEY)
//  if len(signingMethod) != 1 {
//    return nil, errors.New("signing method is malformed")
//  }
//
//  return &AuthorizationConfig{
//    Token:         tokenStr[0],
//    SigningMethod: jwt.GetSigningMethod(signingMethod[0]),
//    Scheme:        scheme[0],
//    KeyFunc: func(token *jwt.Token) (interface{}, error) {
//      return []byte(secret[0]), nil
//    },
//  }, nil
//}

type AuthorizationConfig[T jwt.Claims] struct {
  SigningMethod jwt.SigningMethod
  Scheme        string
  ClaimsKey     string
  KeyFunc       jwt.Keyfunc
}

func (a *AuthorizationConfig[T]) parseToken(tokenStr string) (*T, error) {
  var claims T
  _, err := jwt.ParseWithClaims(tokenStr, claims, a.KeyFunc)
  if err != nil {
    return nil, fmt.Errorf("invalid auth token: %v", err)
  }
  return &claims, nil
}

func (a *AuthorizationConfig[T]) getTokenFromMD(md metadata.MD) (string, error) {
  vals := md.Get(a.Scheme)
  if len(vals) != 1 {
    return "", errors.New("token metadata is malformed")
  }
  val := vals[0]
  scheme, token, found := strings.Cut(val, " ")
  if !found {
    return "", errors.New("bad authorization string")
  }
  if !strings.EqualFold(scheme, a.Scheme) {
    return "", errors.New("token type is different")
  }
  return token, nil
}

func unaryAuthorization[T jwt.Claims](conf AuthorizationConfig[T], checkFunc PermCheckFunc[T]) grpc.UnaryServerInterceptor {
  return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    // Parse token
    md, found := metadata.FromIncomingContext(ctx)
    if !found {
      return nil, status.Error(codes.Unauthenticated, "no metadata found")
    }

    tokenStr, err := conf.getTokenFromMD(md)
    if err != nil {
      return nil, status.Error(codes.Unauthenticated, err.Error())
    }

    claims, err := conf.parseToken(tokenStr)
    if err != nil {
      return nil, status.Error(codes.Unauthenticated, err.Error())
    }

    // Check permission
    meta := interceptors.NewServerCallMeta(info.FullMethod, nil, req)
    if !checkFunc(claims, meta) {
      return nil, status.New(codes.PermissionDenied, sharedErr.ErrUnauthorized.Error()).Err()
    }

    ctx = context.WithValue(ctx, conf.ClaimsKey, claims)
    return handler(ctx, req)
  }
}

func streamAuthorization[T jwt.Claims](conf AuthorizationConfig[T], checkFunc PermCheckFunc[T]) grpc.StreamServerInterceptor {
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

    // Parse token
    md, found := metadata.FromIncomingContext(currCtx)
    if !found {
      return status.Error(codes.Unauthenticated, "no metadata found")
    }

    tokenStr, err := conf.getTokenFromMD(md)
    if err != nil {
      return status.Error(codes.Unauthenticated, err.Error())
    }

    claims, err := conf.parseToken(tokenStr)
    if err != nil {
      return status.Error(codes.Unauthenticated, err.Error())
    }

    // Check permission
    meta := interceptors.NewServerCallMeta(info.FullMethod, info, nil)
    if !checkFunc(claims, meta) {
      return status.New(codes.PermissionDenied, sharedErr.ErrUnauthorized.Error()).Err()
    }

    // Set new context
    newCtx = context.WithValue(currCtx, conf.ClaimsKey, claims)
    wrappedServerStream := middleware.WrapServerStream(ss)
    wrappedServerStream.WrappedContext = newCtx

    return handler(srv, wrappedServerStream)
  }
}

type PermCheckFunc[T jwt.Claims] func(claims *T, meta interceptors.CallMeta) bool

type SelectorMatchFunc func(ctx context.Context, callMeta interceptors.CallMeta) bool

func UnaryServerAuth[T jwt.Claims](config AuthorizationConfig[T], f PermCheckFunc[T], sf SelectorMatchFunc) grpc.UnaryServerInterceptor {
  return selector.UnaryServerInterceptor(
    unaryAuthorization[T](config, f),
    selector.MatchFunc(sf),
  )
}

func StreamServerAuth[T jwt.Claims](config AuthorizationConfig[T], f PermCheckFunc[T], sf SelectorMatchFunc) grpc.StreamServerInterceptor {
  return selector.StreamServerInterceptor(
    streamAuthorization[T](config, f),
    selector.MatchFunc(sf),
  )
}
