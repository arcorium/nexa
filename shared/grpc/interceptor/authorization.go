package interceptor

import (
  "context"
  "errors"
  "fmt"
  "github.com/golang-jwt/jwt/v5"
  middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/metadata"
  "google.golang.org/grpc/status"
  "nexa/shared/constant"
  sharedErr "nexa/shared/errors"
  sharedJwt "nexa/shared/jwt"
  sharedUtil "nexa/shared/util"
  "strings"
)

func parseToken(tokenStr string, config *AuthData) (*sharedJwt.UserClaims, error) {
  var claims *sharedJwt.UserClaims
  token, err := jwt.ParseWithClaims(tokenStr, claims, config.KeyFunc)
  if err != nil {
    return nil, fmt.Errorf("invalid auth token: %v", err)
  }
  sharedUtil.DoNothing(token)
  return claims, nil
}

func getTokenFromMD(md metadata.MD, data *AuthData) (string, error) {
  vals := md.Get(data.Scheme)
  if len(vals) != 1 {
    return "", errors.New("token metadata is malformed")
  }
  val := vals[0]
  scheme, token, found := strings.Cut(val, " ")
  if !found {
    return "", errors.New("bad authorization string")
  }
  if !strings.EqualFold(scheme, data.Scheme) {
    return "", errors.New("token type is different")
  }
  return token, nil
}

func getDataFromMD(md metadata.MD) (*AuthData, error) {
  tokenStr := md.Get(constant.TOKEN_METADATA_KEY)
  if len(tokenStr) != 1 {
    return nil, errors.New("token metadata is malformed")
  }
  scheme := md.Get(constant.TOKEN_SCHEME_METADATA_KEY)
  if len(scheme) != 1 {
    return nil, errors.New("token type is malformed")
  }
  secret := md.Get(constant.JWT_SECRET_METADATA_KEY)
  if len(secret) != 1 {
    return nil, errors.New("secret key is malformed")
  }
  signingMethod := md.Get(constant.JWT_SIGNING_METHOD_METADATA_KEY)
  if len(signingMethod) != 1 {
    return nil, errors.New("signing method is malformed")
  }

  return &AuthData{
    Token:         tokenStr[0],
    SigningMethod: jwt.GetSigningMethod(signingMethod[0]),
    Scheme:        scheme[0],
    KeyFunc: func(token *jwt.Token) (interface{}, error) {
      return []byte(secret[0]), nil
    },
  }, nil
}

type AuthData struct {
  Token         string
  SigningMethod jwt.SigningMethod
  Scheme        string
  KeyFunc       jwt.Keyfunc
}

func Authorization(ctx context.Context) (context.Context, error) {
  md, found := metadata.FromIncomingContext(ctx)
  if !found {
    return nil, status.Error(codes.Unauthenticated, "no metadata found")
  }

  data, err := getDataFromMD(md)
  if err != nil {
    return ctx, status.Error(codes.Unauthenticated, err.Error())
  }

  tokenStr, err := getTokenFromMD(md, data)
  if err != nil {
    return ctx, status.Error(codes.Unauthenticated, err.Error())
  }

  claims, err := parseToken(tokenStr, data)
  if err != nil {
    return ctx, status.Error(codes.Unauthenticated, err.Error())
  }

  return context.WithValue(ctx, constant.CLAIMS_CONTEXT_KEY, claims), nil
}

type PermCheckFunc func(claims *sharedJwt.UserClaims, meta interceptors.CallMeta) bool

type SelectorMatchFunc func(ctx context.Context, callMeta interceptors.CallMeta) bool

func UnaryServerAuth(f PermCheckFunc, sf SelectorMatchFunc) grpc.UnaryServerInterceptor {
  return selector.UnaryServerInterceptor(
    permissionCheck(f),
    selector.MatchFunc(sf),
  )
}

func StreamServerAuth(f PermCheckFunc, sf SelectorMatchFunc) grpc.StreamServerInterceptor {
  return selector.StreamServerInterceptor(
    permissionCheckStream(f),
    selector.MatchFunc(sf),
  )
}

func permissionCheck(f PermCheckFunc) grpc.UnaryServerInterceptor {
  return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    var err error

    // Parse and validate token
    ctx, err = Authorization(ctx)
    if err != nil {
      return nil, err
    }

    claims, err := sharedJwt.GetClaimsFromCtx(ctx)
    if err != nil { // NOTE: error check is not necessary
      return nil, status.New(codes.Unauthenticated, err.Error()).Err()
    }

    // Check permission based on route
    meta := interceptors.NewServerCallMeta(info.FullMethod, nil, req)
    if !f(claims, meta) {
      return nil, status.New(codes.PermissionDenied, sharedErr.ErrUnauthorized.Error()).Err()
    }
    return handler(ctx, req)
  }
}

func permissionCheckStream(f PermCheckFunc) grpc.StreamServerInterceptor {
  return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
    var err error
    var newCtx context.Context

    var currCtx context.Context
    wrappedStream, ok := ss.(*middleware.WrappedServerStream)
    if !ok {
      currCtx = ss.Context()
    } else {
      currCtx = wrappedStream.WrappedContext
    }

    // Parse and validate token
    newCtx, err = Authorization(currCtx)
    if err != nil {
      return err
    }

    wrappedServerStream := middleware.WrapServerStream(ss)
    wrappedServerStream.WrappedContext = newCtx

    claims, err := sharedJwt.GetClaimsFromCtx(wrappedServerStream.WrappedContext)
    if err != nil {
      return status.New(codes.Unauthenticated, err.Error()).Err()
    }

    // Check permission based on route
    meta := interceptors.NewServerCallMeta(info.FullMethod, info, nil)
    if !f(claims, meta) {
      return status.New(codes.PermissionDenied, sharedErr.ErrUnauthorized.Error()).Err()
    }
    return handler(srv, wrappedServerStream)
  }
}
