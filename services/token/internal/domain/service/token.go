package service

import (
  "context"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/token/internal/domain/dto"
)

type IToken interface {
  // Request to create a token
  Request(ctx context.Context, dto *dto.TokenCreateDTO) (dto.TokenResponseDTO, status.Object)
  // Verify token and return the user id related to the token.
  // For example, it could be used for reset password
  Verify(ctx context.Context, dto *dto.TokenVerifyDTO) (types.Id, status.Object)
  // AuthVerify used for verifying token that should be related to some user.
  // For example, it could be used for login
  AuthVerify(ctx context.Context, dto *dto.TokenAuthVerifyDTO) status.Object
}
