package mapper

import (
  "nexa/services/authorization/internal/domain/dto"
  "nexa/services/authorization/shared/proto"
  sharedDto "nexa/shared/dto"
  sharedProto "nexa/shared/proto"
  "nexa/shared/util"
  "nexa/shared/wrapper"
)

func ToRoleCreateDTO(input *proto.RoleCreateInput) dto.RoleCreateDTO {
  return dto.RoleCreateDTO{
    Name:        input.Name,
    Description: wrapper.NewNullable(input.Description),
  }
}

func ToRoleUpdateDTO(input *proto.RoleUpdateInput) dto.RoleUpdateDTO {
  return dto.RoleUpdateDTO{
    Id:          input.Id,
    Name:        wrapper.NewNullable(input.Name),
    Description: wrapper.NewNullable(input.Description),
  }
}

func ToAddPermissionsDTO(input *proto.RoleAddPermissionsInput) dto.RoleAddPermissionsDTO {
  return dto.RoleAddPermissionsDTO{
    RoleId:        input.RoleId,
    PermissionIds: input.PermissionIds,
  }
}

func ToRemovePermissionsDTO(input *proto.RoleRemovePermissionsInput) dto.RoleRemovePermissionsDTO {
  return dto.RoleRemovePermissionsDTO{
    RoleId:        input.RoleId,
    PermissionIds: input.PermissionIds,
  }
}

func ToAddUsersDTO(input *proto.RoleAddUserInput) dto.RoleAddUsersDTO {
  return dto.RoleAddUsersDTO{
    UserId:  input.UserId,
    RoleIds: input.RoleIds,
  }
}

func ToRemoveUsersDTO(input *proto.RoleRemoveUserInput) dto.RoleRemoveUsersDTO {
  return dto.RoleRemoveUsersDTO{
    UserId:  input.UserId,
    RoleIds: input.RoleIds,
  }
}

func ToRoleResponse(resp *dto.RoleResponseDTO) *proto.RoleResponse {
  return &proto.RoleResponse{
    Id:          resp.Id,
    Name:        resp.Name,
    Description: resp.Description,
  }
}

func ToRoleResponses(result *sharedDto.PagedElementResult[dto.RoleResponseDTO]) *proto.RoleFindAllOutput {
  return &proto.RoleFindAllOutput{
    Details: &sharedProto.PagedElementOutput{
      Element:       result.Element,
      Page:          result.Page,
      TotalElements: result.TotalElements,
      TotalPages:    result.TotalPages,
    },
    Roles: util.CastSlice(result.Data, ToRoleResponse),
  }
}
