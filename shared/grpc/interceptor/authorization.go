package interceptor

import (
  "context"
  "errors"
  "fmt"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
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
  CheckFunc     PermCheckFunc[T]
}

func (a *AuthorizationConfig[T]) Valid() bool {
  return a.SigningMethod != nil && len(a.Scheme) != 0 && len(a.ClaimsKey) != 0 && a.KeyFunc != nil
}

func (a *AuthorizationConfig[T]) extract(ctx context.Context) (*T, error) {
  // Parse token
  md, found := metadata.FromIncomingContext(ctx)
  if !found {
    return nil, errors.New("no metadata found")
  }

  tokenStr, err := a.getTokenFromMD(md)
  if err != nil {
    return nil, err
  }

  claims, err := a.parseToken(tokenStr)
  if err != nil {
    return nil, err
  }

  return claims, nil
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

// unaryAuthorization extract claims from context and check the permission based on argument. If the permission check function is nil
// it will bypass it and only do the claim extraction.
func unaryAuthorization[T jwt.Claims](conf AuthorizationConfig[T]) grpc.UnaryServerInterceptor {
  return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    // extract claims from token in context
    claims, err := conf.extract(ctx)
    if err != nil {
      return nil, status.Error(codes.Unauthenticated, err.Error())
    }

    // Check permission
    meta := interceptors.NewServerCallMeta(info.FullMethod, nil, req)
    if conf.CheckFunc != nil && !conf.CheckFunc(claims, meta) {
      return nil, status.New(codes.PermissionDenied, sharedErr.ErrUnauthorized.Error()).Err()
    }

    ctx = context.WithValue(ctx, conf.ClaimsKey, claims)
    return handler(ctx, req)
  }
}

// streamAuthorization works like unaryAuthorization but for stream
func streamAuthorization[T jwt.Claims](conf AuthorizationConfig[T]) grpc.StreamServerInterceptor {
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

    claims, err := conf.extract(currCtx)
    if err != nil {
      return status.Error(codes.Unauthenticated, err.Error())
    }

    // Check permission
    meta := interceptors.NewServerCallMeta(info.FullMethod, info, nil)
    if conf.CheckFunc != nil && !conf.CheckFunc(claims, meta) {
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

// UnaryServerAuth create unary server interceptor for authorization which will extract token into claims
// and doing the permission check based on the config. Selector used to determine if the request should be forwarded
// to authorization or not. type T is used to determine the structure of the claims
func UnaryServerAuth[T jwt.Claims](config AuthorizationConfig[T], sf SelectorMatchFunc) grpc.UnaryServerInterceptor {
  return selector.UnaryServerInterceptor(
    unaryAuthorization[T](config),
    selector.MatchFunc(sf),
  )
}

// StreamServerAuth works the same as UnaryServerAuth, but it works for stream.
func StreamServerAuth[T jwt.Claims](config AuthorizationConfig[T], sf SelectorMatchFunc) grpc.StreamServerInterceptor {
  return selector.StreamServerInterceptor(
    streamAuthorization[T](config),
    selector.MatchFunc(sf),
  )
}

type CombinationAuthConfig struct {
  AuthSelector      SelectorMatchFunc
  User              AuthorizationConfig[sharedJwt.UserClaims]
  ProtectedSelector SelectorMatchFunc
  Protected         AuthorizationConfig[sharedJwt.TemporaryClaims]
}

func (c *CombinationAuthConfig) Valid() bool {
  return c.User.Valid() && c.Protected.Valid()
}

// UnaryServerCombinationAuth create unary server interceptor for authorization both for user and protected API.
// It is used for minimalize checking redundancy with bypassing the user authorization check when it is
// already handled by protected API authorization. Authorization selector only will be called if only
// the request(rpc) is not handled by protected API authorization.
func UnaryServerCombinationAuth(conf CombinationAuthConfig) grpc.UnaryServerInterceptor {
  if !conf.Valid() {
    panic("Config has empty values on non-nilable field")
  }
  return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
    meta := interceptors.NewServerCallMeta(info.FullMethod, nil, req)
    // Run selector
    if conf.ProtectedSelector(ctx, meta) {
      // Handle for auth
      return unaryAuthorization(conf.Protected)(ctx, req, info, handler) // WARN: Save it as variable outside closure?
    }

    // Doesn't need authorization
    if !conf.AuthSelector(ctx, meta) {
      return handler(ctx, req)
    }

    return unaryAuthorization(conf.User)(ctx, req, info, handler) // WARN: Save it as variable outside closure?
  }
}

// StreamServerCombinationAuth works the same as UnaryServerCombinationAuth, but it works for stream.
func StreamServerCombinationAuth(conf CombinationAuthConfig) grpc.StreamServerInterceptor {
  if !conf.Valid() {
    panic("Config has empty values on non-nilable field")
  }

  return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
    meta := interceptors.NewServerCallMeta(info.FullMethod, info, nil)
    // Run selector
    if conf.ProtectedSelector(ss.Context(), meta) {
      // Handle for auth
      return streamAuthorization(conf.Protected)(srv, ss, info, handler) // WARN: Save it as variable outside closure?
    }

    // Doesn't need authorization
    if !conf.AuthSelector(ss.Context(), meta) {
      return handler(srv, ss)
    }

    return streamAuthorization(conf.User)(srv, ss, info, handler) // WARN: Save it as variable outside closure?
  }
}
