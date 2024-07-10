package mapper

import (
  authZv1 "github.com/arcorium/nexa/proto/gen/go/authorization/v1"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "nexa/services/authorization/internal/domain/dto"
)

func ToIsAuthorizationDTO(request *authZv1.CheckUserRequest) (dto.IsAuthorizationDTO, error) {
  id, err := types.IdFromString(request.UserId)
  if err != nil {
    return dto.IsAuthorizationDTO{}, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  dtos := dto.IsAuthorizationDTO{
    UserId:             id,
    ExpectedPermission: request.ExpectedPermission,
  }

  return dtos, sharedUtil.ValidateStruct(&dtos)
}
