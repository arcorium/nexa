package mapper

import (
  tokenv1 "github.com/arcorium/nexa/proto/gen/go/token/v1"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  "google.golang.org/protobuf/types/known/timestamppb"
  "nexa/services/token/internal/domain/dto"
  "nexa/services/token/internal/domain/entity"
)

func ToCreateTokenDTO(req *tokenv1.CreateTokenRequest) (dto.TokenCreateDTO, error) {
  var fieldErrors []sharedErr.FieldError

  userId, err := types.IdFromString(req.UserId)
  if err != nil {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("user_id", err))
  }
  usage, err := entity.NewTokenUsage(uint8(req.Usage))
  if err != nil {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("usage", err))
  }

  tokenType, err := entity.NewTokenType(uint8(req.Type))
  if err != nil {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("type", err))
  }

  if len(fieldErrors) > 0 {
    return dto.TokenCreateDTO{}, sharedErr.GrpcFieldErrors2(fieldErrors...)
  }

  return dto.TokenCreateDTO{
    UserId: userId,
    Type:   tokenType,
    Length: req.TokenLength,
    Usage:  usage,
  }, nil
}

func ToTokenVerifyDTO(req *tokenv1.VerifyTokenRequest) (dto.TokenVerifyDTO, error) {
  var fieldErrors []sharedErr.FieldError
  if len(req.Token) == 0 {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("token", sharedErr.ErrFieldEmpty))
  }

  usage, err := entity.NewTokenUsage(uint8(req.Usage))
  if err != nil {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("usage", err))
  }

  if len(fieldErrors) > 0 {
    return dto.TokenVerifyDTO{}, sharedErr.GrpcFieldErrors2(fieldErrors...)
  }

  return dto.TokenVerifyDTO{
    Token:         req.Token,
    ExpectedUsage: usage,
  }, nil
}

func ToTokenAuthVerifyDTO(req *tokenv1.VerifyAuthTokenRequest) (dto.TokenAuthVerifyDTO, error) {
  var fieldErrors []sharedErr.FieldError
  if len(req.Token) == 0 {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("token", sharedErr.ErrFieldEmpty))
  }

  usage, err := entity.NewTokenUsage(uint8(req.Usage))
  if err != nil {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("usage", err))
  }

  userId, err := types.IdFromString(req.UserId)
  if err != nil {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("user_id", err))
  }

  if len(fieldErrors) > 0 {
    return dto.TokenAuthVerifyDTO{}, sharedErr.GrpcFieldErrors2(fieldErrors...)
  }

  return dto.TokenAuthVerifyDTO{
    TokenVerifyDTO: dto.TokenVerifyDTO{
      Token:         req.Token,
      ExpectedUsage: usage,
    },
    ExpectedUserId: userId,
  }, nil
}

func ToProtoToken(resp *dto.TokenResponseDTO) *tokenv1.Token {
  return &tokenv1.Token{
    Token:     resp.Token,
    UserId:    resp.UserId.String(),
    Usage:     tokenv1.TokenUsage(resp.Usage),
    ExpiredAt: timestamppb.New(resp.ExpiredAt),
  }
}
