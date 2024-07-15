package mapper

import (
  "nexa/services/token/internal/domain/dto"
  "nexa/services/token/internal/domain/entity"
)

func ToTokenResponseDTO(token *entity.Token) dto.TokenResponseDTO {
  return dto.TokenResponseDTO{
    Token:     token.Token,
    UserId:    token.UserId,
    Usage:     token.Usage,
    ExpiredAt: token.ExpiredAt,
  }
}
