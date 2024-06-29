package mapper

import (
  "nexa/services/authentication/internal/domain/dto"
  domain "nexa/services/authentication/internal/domain/entity"
)

func ToTokenResponseDTO(token *domain.Token) dto.TokenResponseDTO {
  return dto.TokenResponseDTO{
    Token:     token.Token,
    Usage:     token.Usage.Underlying(),
    ExpiredAt: token.ExpiredAt,
  }
}
