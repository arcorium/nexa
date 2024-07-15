package external

import (
  "context"
  authZv1 "github.com/arcorium/nexa/proto/gen/go/authorization/v1"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/external"
  "nexa/services/authentication/util"
)

func NewRoleClient(conn grpc.ClientConnInterface) external.IRoleClient {
  return &roleClient{
    client: authZv1.NewRoleServiceClient(conn),
    tracer: util.GetTracer(),
  }
}

type roleClient struct {
  client authZv1.RoleServiceClient

  tracer trace.Tracer
}

func (r roleClient) GetUserRoles(ctx context.Context, userId types.Id) ([]dto.RoleResponseDTO, error) {
  ctx, span := r.tracer.Start(ctx, "RoleClient.GetUserRoles")
  defer span.End()

  req := authZv1.GetUserRolesRequest{
    UserId:            userId.String(),
    IncludePermission: true,
  }

  roles, err := r.client.GetUsers(ctx, &req)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  // Mapping
  rolePerms, ierr := sharedUtil.CastSliceErrs(roles.RolePermissions, func(rolePrem *authZv1.RolePermission) (dto.RoleResponseDTO, error) {
    roleId, err := types.IdFromString(rolePrem.Role.Id)
    if err != nil {
      return dto.RoleResponseDTO{}, err
    }
    permissions, ierr := sharedUtil.CastSliceErrs(rolePrem.Permissions, func(perm *authZv1.Permission) (dto.Permission, error) {
      permId, err := types.IdFromString(perm.Id)
      if err != nil {
        return dto.Permission{}, err
      }

      return dto.Permission{
        Id:   permId,
        Code: perm.Code,
      }, nil
    })

    if !ierr.IsNil() {
      return dto.RoleResponseDTO{}, ierr
    }

    return dto.RoleResponseDTO{
      Id:          roleId,
      Role:        rolePrem.Role.Name,
      Permissions: permissions,
    }, nil
  })

  if !ierr.IsNil() {
    spanUtil.RecordError(err, span)
    return nil, ierr
  }
  return rolePerms, nil
}

func (r roleClient) RemoveUserRoles(ctx context.Context, userId types.Id) error {
  ctx, span := r.tracer.Start(ctx, "RoleClient.RemoveUserRoles")
  defer span.End()

  req := authZv1.RemoveUserRolesRequest{
    UserId:  userId.String(),
    RoleIds: nil, // Nil means delete all
  }

  _, err := r.client.RemoveUser(ctx, &req)
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }
 
  return err
}
