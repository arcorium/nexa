package mapper

import (
  authZv1 "github.com/arcorium/nexa/proto/gen/go/authorization/v1"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "nexa/services/authorization/internal/domain/dto"
)

func ToCreatePermissionDTO(input *authZv1.CreatePermissionRequest) (dto.PermissionCreateDTO, error) {
  dtos := dto.PermissionCreateDTO{
    Resource: input.Resource,
    Action:   input.Action,
  }

  err := sharedUtil.ValidateStruct(&dtos)
  return dtos, err
}

func ToProtoPermission(responseDTO *dto.PermissionResponseDTO) *authZv1.Permission {
  return &authZv1.Permission{
    Id:   responseDTO.Id.String(),
    Code: responseDTO.Code,
  }
}
