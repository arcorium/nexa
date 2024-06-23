package interceptor

import (
  "context"
  "errors"
  "fmt"
  "github.com/golang-jwt/jwt/v5"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/metadata"
  "google.golang.org/grpc/status"
  "nexa/shared/constant"
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
