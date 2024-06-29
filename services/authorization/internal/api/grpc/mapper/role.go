package mapper

import (
  authZv1 "nexa/proto/gen/go/authorization/v1"
  "nexa/services/authorization/internal/domain/dto"
  sharedUtil "nexa/shared/util"
  "nexa/shared/wrapper"
)

func ToRoleCreateDTO(req *authZv1.RoleCreateRequest) dto.RoleCreateDTO {
  return dto.RoleCreateDTO{
    Name:        req.Name,
    Description: wrapper.NewNullable(req.Description),
  }
}

func ToRoleUpdateDTO(req *authZv1.RoleUpdateRequest) dto.RoleUpdateDTO {
  return dto.RoleUpdateDTO{
    Id:          req.Id,
    Name:        wrapper.NewNullable(req.Name),
    Description: wrapper.NewNullable(req.Description),
  }
}

func ToAddRolePermissionsDTO(req *authZv1.RoleAppendPermissionsRequest) dto.RoleAddPermissionsDTO {
  return dto.RoleAddPermissionsDTO{
    RoleId:        req.RoleId,
    PermissionIds: req.PermissionIds,
  }
}

func ToRemoveRolePermissionsDTO(req *authZv1.RoleRemovePermissionsRequest) dto.RoleRemovePermissionsDTO {
  return dto.RoleRemovePermissionsDTO{
    RoleId:        req.RoleId,
    PermissionIds: req.PermissionIds,
  }
}

func ToAddUsersDTO(input *authZv1.AddUserRolesRequest) dto.RoleAddUsersDTO {
  return dto.RoleAddUsersDTO{
    UserId:  input.UserId,
    RoleIds: input.RoleIds,
  }
}

func ToRemoveUsersDTO(input *authZv1.RemoveUserRolesRequest) dto.RoleRemoveUsersDTO {
  return dto.RoleRemoveUsersDTO{
    UserId:  input.UserId,
    RoleIds: input.RoleIds,
  }
}

func ToRoleResponse(resp *dto.RoleResponseDTO) *authZv1.Role {
  return &authZv1.Role{
    Id:          resp.Id,
    Name:        resp.Name,
    Description: resp.Description,
  }
}

func ToProtoRolePermission(responseDTO *dto.RoleResponseDTO, includePerm bool) *authZv1.RolePermission {
  rolePerm := &authZv1.RolePermission{
    Role: &authZv1.Role{
      Id:          responseDTO.Id,
      Name:        responseDTO.Name,
      Description: responseDTO.Description,
    },
  }
  if includePerm {
    rolePerm.Permissions = sharedUtil.CastSliceP(responseDTO.Permissions, func(perm *dto.PermissionResponseDTO) *authZv1.Permission {
      return &authZv1.Permission{
        Id:   perm.Id,
        Code: perm.Code,
      }
    })
  }
  return rolePerm
}
