package mapper

import (
  "google.golang.org/protobuf/types/known/timestamppb"
  authv1 "nexa/proto/gen/go/authentication/v1"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/entity"
  sharedErr "nexa/shared/errors"
  "nexa/shared/types"
)

func ToCreateTokenDTO(req *authv1.TokenCreateRequest) (dto.TokenCreateDTO, error) {
  var fieldErrors []sharedErr.FieldError

  userId, err := types.IdFromString(req.UserId)
  if err != nil {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("user_id", err))
  }
  usage, err := entity.NewTokenUsage(uint8(req.Usage))
  if err != nil {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("usage", err))
  }

  if len(fieldErrors) > 0 {
    return dto.TokenCreateDTO{}, sharedErr.GrpcFieldErrors2(fieldErrors...)
  }

  return dto.TokenCreateDTO{
    UserId: userId,
    Usage:  usage,
  }, nil
}

func ToTokenVerifyDTO(req *authv1.TokenVerifyRequest) (dto.TokenVerifyDTO, error) {
  var fieldErrors []sharedErr.FieldError
  if len(req.Token) == 0 {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("user_id", sharedErr.ErrFieldEmpty))
  }

  usage, err := entity.NewTokenUsage(uint8(req.Usage))
  if err != nil {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("usage", err))
  }

  if len(fieldErrors) > 0 {
    return dto.TokenVerifyDTO{}, sharedErr.GrpcFieldErrors2(fieldErrors...)
  }

  return dto.TokenVerifyDTO{
    Token: req.Token,
    Usage: usage,
  }, nil
}

func ToProtoToken(resp *dto.TokenResponseDTO) *authv1.Token {
  return &authv1.Token{
    Token:     resp.Token,
    Usage:     authv1.TokenUsage(resp.Usage),
    ExpiredAt: timestamppb.New(resp.ExpiredAt),
  }
}
