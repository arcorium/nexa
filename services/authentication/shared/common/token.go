package common

import (
	"github.com/golang-jwt/jwt/v5"
	"nexa/services/authentication/constant"
	"nexa/shared/status"
	"nexa/shared/types"
	"time"
)

func GenerateAccessToken(signingMethod jwt.SigningMethod, exp time.Duration, secretKey []byte, userId types.Id, username string) (types.Id, string, status.Object) {
	ct := time.Now()
	expAt := jwt.NewNumericDate(ct.Add(exp))

	id := types.NewId()
	accessClaims := AccessTokenClaims{
		UserId:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    constant.CLAIMS_ISSUER,
			ExpiresAt: expAt,
			NotBefore: expAt,
			IssuedAt:  jwt.NewNumericDate(ct),
			ID:        id.Underlying().String(),
		},
	}
	accessToken := jwt.NewWithClaims(signingMethod, accessClaims)
	accessSignedString, err := accessToken.SignedString(secretKey)
	if err != nil {
		return types.NullId(), "", status.ErrInternal(err)
	}
	return id, accessSignedString, status.SuccessInternal()
}

func GenerateRefreshToken(signingMethod jwt.SigningMethod, exp time.Duration, secretKey []byte) (types.Id, string, status.Object) {

	ct := time.Now()
	expAt := jwt.NewNumericDate(ct.Add(exp))
	id := types.NewId()

	refreshClaims := RefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    constant.CLAIMS_ISSUER,
			ExpiresAt: expAt,
			NotBefore: expAt,
			IssuedAt:  jwt.NewNumericDate(ct),
			ID:        id.Underlying().String(),
		},
	}
	refreshToken := jwt.NewWithClaims(signingMethod, refreshClaims)
	refreshSignedString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return types.NullId(), "", status.ErrInternal(err)
	}
	return id, refreshSignedString, status.SuccessInternal()
}

type PairTokens struct {
	RefreshToken struct {
		Id     types.Id
		String string
	}
	AccessToken struct {
		Id     types.Id
		String string
	}
}

func GeneratePairTokens(signingMethod jwt.SigningMethod, exp time.Duration, secretKey []byte, userId types.Id, username string) (PairTokens, status.Object) {

	accessId, accessToken, stats := GenerateAccessToken(signingMethod, exp, secretKey, userId, username)
	if stats.IsError() {
		return PairTokens{}, stats
	}
	refreshId, refreshToken, stats := GenerateRefreshToken(signingMethod, exp, secretKey)
	if stats.IsError() {
		return PairTokens{}, stats
	}

	return PairTokens{
		RefreshToken: struct {
			Id     types.Id
			String string
		}{
			Id:     refreshId,
			String: refreshToken,
		},
		AccessToken: struct {
			Id     types.Id
			String string
		}{
			Id:     accessId,
			String: accessToken,
		},
	}, status.SuccessInternal()
}
