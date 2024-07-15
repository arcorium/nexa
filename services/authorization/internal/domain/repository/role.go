package repository

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util/repo"
  "nexa/services/authorization/internal/domain/entity"
)

type IRole interface {
  // FindByIds get roles bad on the ids provided
  FindByIds(ctx context.Context, ids ...types.Id) ([]entity.Role, error)
  // FindByUserId get user roles
  FindByUserId(ctx context.Context, userId types.Id) ([]entity.Role, error)
  // FindByName get roles by the name
  FindByName(ctx context.Context, name string) (entity.Role, error)
  // GetAll get all roles
  Get(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Role], error)
  // Create create new role
  Create(ctx context.Context, role *entity.Role) error
  // Patch update role
  Patch(ctx context.Context, role *entity.PatchedRole) error
  // Delete delete role
  Delete(ctx context.Context, id types.Id) error
  // AddPermissions add permission into role
  AddPermissions(ctx context.Context, roleId types.Id, permissionIds ...types.Id) error
  // RemovePermissions remove permission from role
  RemovePermissions(ctx context.Context, roleId types.Id, permissionIds ...types.Id) error
  // ClearPermission remove all role's permissions
  ClearPermission(ctx context.Context, roleId types.Id) error
  // AddUser add roles into user
  AddUser(ctx context.Context, userId types.Id, roleIds ...types.Id) error
  // RemoveUser remove roles from user
  RemoveUser(ctx context.Context, userId types.Id, roleIds ...types.Id) error
  // ClearUser remove all user's roles
  ClearUser(ctx context.Context, userId types.Id) error
}
