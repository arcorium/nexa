package mapper

import (
  authNv1 "github.com/arcorium/nexa/proto/gen/go/authentication/v1"
  "google.golang.org/protobuf/types/known/timestamppb"
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
