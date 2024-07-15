package mapper

import (
  tokenv1 "github.com/arcorium/nexa/proto/gen/go/token/v1"
  "google.golang.org/protobuf/types/known/timestamppb"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/util"
)

func ToProtoTokenResponse(respDTO *dto.TokenResponseDTO) *tokenv1.Token {
  return &tokenv1.Token{
    Token:     respDTO.Token,
    UserId:    respDTO.UserId.String(),
    Usage:     util.TokenPurposeToUsage(respDTO.Usage),
    ExpiredAt: timestamppb.New(respDTO.ExpiredAt),
  }
}
