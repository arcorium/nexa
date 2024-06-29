package mapper

import (
  "google.golang.org/protobuf/types/known/timestamppb"
  authv1 "nexa/proto/gen/go/authentication/v1"
  "nexa/services/authentication/internal/domain/dto"
)

func ToCreateTokenDTO(input *authv1.TokenCreateRequest) dto.TokenCreateDTO {
  return dto.TokenCreateDTO{
    UserId: input.UserId,
    Usage:  uint8(input.Usage),
  }
}

func ToTokenVerifyDTO(input *authv1.TokenVerifyRequest) dto.TokenVerifyDTO {
  return dto.TokenVerifyDTO{
    Token: input.Token,
    Usage: uint8(input.Usage),
  }
}

func ToProtoToken(resp *dto.TokenResponseDTO) *authv1.Token {
  return &authv1.Token{
    Token:     resp.Token,
    Usage:     authv1.TokenUsage(resp.Usage),
    ExpiredAt: timestamppb.New(resp.ExpiredAt),
  }
}
