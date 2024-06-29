package mapper

import (
  authZv1 "nexa/proto/gen/go/authorization/v1"
  "nexa/services/authorization/internal/domain/dto"
)

func ToPermissionCreateDTO(input *authZv1.PermissionCreateRequest) dto.PermissionCreateDTO {
  return dto.PermissionCreateDTO{
    Code: input.Code,
  }
}

func ToProtoPermission(responseDTO *dto.PermissionResponseDTO) *authZv1.Permission {
  return &authZv1.Permission{
    Id:   responseDTO.Id,
    Code: responseDTO.Code,
  }
}
