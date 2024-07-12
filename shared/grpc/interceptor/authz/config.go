package authz

import (
  "context"
  "crypto/rsa"
  "errors"
  "fmt"
  sharedConst "github.com/arcorium/nexa/shared/constant"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/golang-jwt/jwt/v5"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
  "google.golang.org/grpc/metadata"
  "strings"
)

type UserOption func(config *UserConfig)

type PrivateOption func(config *PrivateConfig)

func NewUserConfig(pubkey *rsa.PublicKey, checkFunc PermCheckFunc[sharedJwt.UserClaims], opts ...UserOption) UserConfig {
  def := UserConfig{
    config: config{
      SigningMethod:    sharedJwt.DefaultSigningMethod,
      AuthorizationKey: sharedConst.DEFAULT_METADATA_AUTHZ_KEY,
      Scheme:           sharedConst.DEFAULT_ACCESS_TOKEN_SCHEME,
      ClaimsKey:        sharedConst.USER_CLAIMS_CONTEXT_KEY,
      KeyFunc: func(token *jwt.Token) (interface{}, error) {
        return pubkey, nil
      },
    },
    CheckFunc: checkFunc,
  }

  for _, opt := range opts {
    opt(&def)
  }
  return def
}

func NewPrivateConfig(pubkey *rsa.PublicKey, checkFunc PermCheckFunc[sharedJwt.PrivateClaims], opts ...PrivateOption) PrivateConfig {
  def := PrivateConfig{
    config: config{
      SigningMethod:    sharedJwt.DefaultSigningMethod,
      AuthorizationKey: sharedConst.DEFAULT_METADATA_AUTHZ_KEY,
      Scheme:           sharedConst.DEFAULT_ACCESS_TOKEN_SCHEME,
      ClaimsKey:        sharedConst.TEMP_CLAIMS_CONTEXT_KEY,
      KeyFunc: func(token *jwt.Token) (interface{}, error) {
        return pubkey, nil
      },
    },
    CheckFunc: checkFunc,
  }

  for _, opt := range opts {
    opt(&def)
  }
  return def
}

type config struct {
  SigningMethod    jwt.SigningMethod
  AuthorizationKey string // Metadata key
  Scheme           string
  ClaimsKey        string
  KeyFunc          jwt.Keyfunc
}

type UserConfig struct {
  config
  CheckFunc PermCheckFunc[sharedJwt.UserClaims]
}

type PrivateConfig struct {
  config
  CheckFunc PermCheckFunc[sharedJwt.PrivateClaims]
}

func (a *config) extractClaims(ctx context.Context, claims jwt.Claims) error {
  // Parse token
  md, found := metadata.FromIncomingContext(ctx)
  if !found {
    return errors.New("no metadata found")
  }

  tokenStr, err := a.getTokenFromMD(md)
  if err != nil {
    return err
  }

  err = a.parseToken(tokenStr, claims)
  if err != nil {
    return err
  }

  return nil
}

func (a *config) parseToken(tokenStr string, claims jwt.Claims) error {
  _, err := jwt.ParseWithClaims(tokenStr, claims, a.KeyFunc)
  if err != nil {
    return fmt.Errorf("invalid auth token: %v", err)
  }
  return nil
}

func (a *config) getTokenFromMD(md metadata.MD) (string, error) {
  vals := md.Get(a.AuthorizationKey)
  if len(vals) != 1 {
    return "", errors.New("metadata has no authorization, expected metadata key: " + a.AuthorizationKey)
  }
  val := vals[0]
  scheme, token, found := strings.Cut(val, " ")
  if !found {
    return "", errors.New("bad authorization string")
  }
  if !strings.EqualFold(scheme, a.Scheme) {
    return "", errors.New("different token scheme, expected: " + a.Scheme)
  }
  return token, nil
}

type CombinationConfig struct {
  Selector     CombinationMatchFunc
  SkipServices []string
  User         UserConfig
  Private      PrivateConfig
}

// CombinationType enum for CombinationMatchFunc which will decide either it should be processed by
// private or user authorization
type CombinationType uint8

const (
  Private CombinationType = iota
  UserAuth
  Public
)

type CombinationMatchFunc func(ctx context.Context, meta interceptors.CallMeta) CombinationType
