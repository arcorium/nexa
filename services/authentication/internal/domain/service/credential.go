package service

import (
  "context"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/shared/status"
  "nexa/shared/types"
)

type ICredential interface {
  Login(ctx context.Context, dto *dto.LoginDTO) (dto.LoginResponseDTO, status.Object)
  Register(ctx context.Context, dto *dto.RegisterDTO) status.Object
  RefreshToken(ctx context.Context, dto *dto.RefreshTokenDTO) (dto.RefreshTokenResponseDTO, status.Object)
  GetCredentials(ctx context.Context, userId types.Id) ([]dto.CredentialResponseDTO, status.Object)
  Logout(ctx context.Context, logoutDTO *dto.LogoutDTO) status.Object
  LogoutAll(ctx context.Context, userId types.Id) status.Object
}
