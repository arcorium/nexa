package grpc

import (
	"context"
	"google.golang.org/grpc"
	"nexa/services/authorization/shared/domain/entity"
	"nexa/services/authorization/shared/proto"
	"nexa/services/user/internal/domain/external"
	"nexa/services/user/internal/infra/external/grpc/mapper"
	"nexa/shared/types"
	"nexa/shared/util"
)

func NewAuthorizationClient(conn grpc.ClientConnInterface) external.IAuthorizationClient {
	return &authorizationClient{
		roleClient:       proto.NewRoleServiceClient(conn),
		permissionClient: proto.NewPermissionServiceClient(conn),
	}
}

type authorizationClient struct {
	roleClient       proto.RoleServiceClient
	permissionClient proto.PermissionServiceClient
}

func (a *authorizationClient) FindUserRoles(ctx context.Context, userId types.Id) ([]entity.Role, error) {
	roles, err := a.roleClient.FindByUserId(ctx, &proto.RoleFindByUserIdInput{UserId: userId.Underlying().String()})

	return util.CastSlice2(roles.Roles, mapper.ToRoleEntity), err
}

func (a *authorizationClient) HasPermission(ctx context.Context, userId types.Id, resource string, actions ...string) error {
	dto := proto.CheckUserInput{
		UserId:      userId.Underlying().String(),
		Permissions: mapper.ToInternalCheckUserInput(resource, actions...),
	}
	_, err := a.permissionClient.CheckUser(ctx, &dto)
	return err
}
