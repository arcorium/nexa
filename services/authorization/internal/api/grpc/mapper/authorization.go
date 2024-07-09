package mapper

import (
  authZv1 "nexa/proto/gen/go/authorization/v1"
  "nexa/services/authorization/internal/domain/dto"
  sharedErr "nexa/shared/errors"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
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
