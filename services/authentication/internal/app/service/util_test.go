package service

import (
  "context"
  "errors"
  sharedConst "github.com/arcorium/nexa/shared/constant"
  sharedJwt "github.com/arcorium/nexa/shared/jwt"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/brianvoe/gofakeit/v7"
  "github.com/golang-jwt/jwt/v5"
  "nexa/services/authentication/constant"
  "time"
)

var dummyErr = errors.New("dummy error")

var dummyId = types.MustCreateId()

var claimsUserId = types.MustCreateId()

func generateRole(actions ...types.Action) sharedJwt.Role {
  return sharedJwt.Role{
    Id:   types.MustCreateId().String(),
    Role: gofakeit.AnimalType(),
    Permissions: sharedUtil.CastSlice(actions, func(action types.Action) string {
      return constant.AUTHN_PERMISSIONS[action]
    }),
  }
}

func generateUserClaims(roles ...sharedJwt.Role) *sharedJwt.UserClaims {
  return &sharedJwt.UserClaims{
    RegisteredClaims: jwt.RegisteredClaims{
      Issuer:    gofakeit.AppName(),
      Subject:   gofakeit.AppName(),
      ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
      NotBefore: jwt.NewNumericDate(time.Now()),
      IssuedAt:  jwt.NewNumericDate(time.Now()),
      ID:        types.MustCreateId().String(),
    },
    CredentialId: types.MustCreateId().String(),
    UserId:       claimsUserId.String(),
    Username:     gofakeit.Username(),
    Roles:        roles,
  }
}

func generateClaimsCtx(actions ...types.Action) context.Context {
  return context.WithValue(context.Background(), sharedConst.USER_CLAIMS_CONTEXT_KEY, generateUserClaims(generateRole(actions...)))
}
