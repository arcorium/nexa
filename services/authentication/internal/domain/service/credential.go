package service

import (
  "context"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/shared/status"
  "nexa/shared/types"
)

type ICredential interface {
  Login(ctx context.Context, dto *dto.LoginDTO) (string, status.Object)
  Register(ctx context.Context, dto *dto.RegisterDTO) status.Object
  RefreshToken(ctx context.Context, dto *dto.RefreshTokenDTO) (string, status.Object)
  GetCurrentCredentials(ctx context.Context) ([]dto.CredentialResponseDTO, status.Object)
  // GetCredentials(ctx context.Context, userId types.Id) ([]dto.CredentialResponseDTO, status.Object)
  Logout(ctx context.Context, credIds ...types.Id) status.Object
  LogoutAll(ctx context.Context) status.Object
  // LogoutAllUser(ctx context.Context, userId types.Id) status.Object
}
