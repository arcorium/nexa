package token

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	UserId string   `json:"user_id"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}
