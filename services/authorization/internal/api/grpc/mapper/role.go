package mapper

import (
  authZv1 "github.com/arcorium/nexa/proto/gen/go/authorization/v1"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "nexa/services/authorization/internal/domain/dto"
)

func ToRoleCreateDTO(req *authZv1.CreateRoleRequest) (dto.RoleCreateDTO, error) {
  dtos := dto.RoleCreateDTO{
    Name:        req.Name,
    Description: types.NewNullable(req.Description),
  }

  err := sharedUtil.ValidateStruct(&dtos)
  return dtos, err
}

func ToRoleUpdateDTO(req *authZv1.UpdateRoleRequest) (dto.RoleUpdateDTO, error) {
  id, err := types.IdFromString(req.Id)
  if err != nil {
    return dto.RoleUpdateDTO{}, sharedErr.NewFieldError("id", err).ToGrpcError()
  }

  return dto.RoleUpdateDTO{
    RoleId:      id,
    Name:        types.NewNullable(req.Name),
    Description: types.NewNullable(req.Description),
  }, nil
}

func toModifyRolesPermissionsDTO(roleId string, permIds []string, allowNil bool) (dto.ModifyRolesPermissionsDTO, error) {
  id, err := types.IdFromString(roleId)
  if err != nil {
    return dto.ModifyRolesPermissionsDTO{}, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  if (len(permIds) == 0) && allowNil {
    return dto.ModifyRolesPermissionsDTO{
      RoleId:        id,
      PermissionIds: nil,
    }, nil
  }

  ids, ierr := sharedUtil.CastSliceErrs(permIds, types.IdFromString)
  if !ierr.IsNil() {
    return dto.ModifyRolesPermissionsDTO{}, ierr.ToGRPCError("role_ids")
  }

  return dto.ModifyRolesPermissionsDTO{
    RoleId:        id,
    PermissionIds: ids,
  }, nil
}

func ToAddRolePermissionsDTO(req *authZv1.AppendRolePermissionsRequest) (dto.ModifyRolesPermissionsDTO, error) {
  return toModifyRolesPermissionsDTO(req.RoleId, req.PermissionIds, false)
}

func ToRemoveRolePermissionsDTO(req *authZv1.RemoveRolePermissionsRequest) (dto.ModifyRolesPermissionsDTO, error) {
  return toModifyRolesPermissionsDTO(req.RoleId, req.PermissionIds, true)
}

func toModifyUserRolesDTO(userId string, roleIds []string, allowNil bool) (dto.ModifyUserRolesDTO, error) {
  id, err := types.IdFromString(userId)
  if err != nil {
    return dto.ModifyUserRolesDTO{}, sharedErr.NewFieldError("user_id", err).ToGrpcError()
  }

  if (len(roleIds) == 0) && allowNil {
    return dto.ModifyUserRolesDTO{
      UserId:  id,
      RoleIds: nil,
    }, nil
  }

  ids, ierr := sharedUtil.CastSliceErrs(roleIds, types.IdFromString)
  if !ierr.IsNil() {
    return dto.ModifyUserRolesDTO{}, ierr.ToGRPCError("role_ids")
  }

  return dto.ModifyUserRolesDTO{
    UserId:  id,
    RoleIds: ids,
  }, nil
}

func ToAddUsersDTO(input *authZv1.AddUserRolesRequest) (dto.ModifyUserRolesDTO, error) {
  return toModifyUserRolesDTO(input.UserId, input.RoleIds, false)
}

func ToRemoveUsersDTO(input *authZv1.RemoveUserRolesRequest) (dto.ModifyUserRolesDTO, error) {
  return toModifyUserRolesDTO(input.UserId, input.RoleIds, true)
}

func ToProtoRole(resp *dto.RoleResponseDTO) *authZv1.Role {
  return &authZv1.Role{
    Id:          resp.Id.String(),
    Name:        resp.Name,
    Description: resp.Description,
  }
}

func ToProtoRolePermission(responseDTO *dto.RoleResponseDTO, includePerm bool) *authZv1.RolePermission {
  rolePerm := &authZv1.RolePermission{
    Role: &authZv1.Role{
      Id:          responseDTO.Id.String(),
      Name:        responseDTO.Name,
      Description: responseDTO.Description,
    },
  }
  if includePerm {
    rolePerm.Permissions = sharedUtil.CastSliceP(responseDTO.Permissions, func(perm *dto.PermissionResponseDTO) *authZv1.Permission {
      return &authZv1.Permission{
        Id:   perm.Id.String(),
        Code: perm.Code,
      }
    })
  }
  return rolePerm
}
