package mapper

import (
  "google.golang.org/protobuf/types/known/timestamppb"
  authNv1 "nexa/proto/gen/go/authentication/v1"
  "nexa/services/user/internal/domain/dto"
  "nexa/services/user/util"
)

func ToProtoTokenResponse(dto *dto.TokenResponseDTO) *authNv1.Token {
  return &authNv1.Token{
    Token:     dto.Token,
    Usage:     util.TokenPurposeToUsage(dto.Purpose),
    ExpiredAt: timestamppb.New(dto.ExpiredAt),
  }
}
