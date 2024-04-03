package grpc

import (
	"google.golang.org/grpc"
	"nexa/services/authorization/shared/domain/entity"
	"nexa/services/authorization/shared/proto"
	"nexa/services/user/internal/domain/external"
	"nexa/shared/types"
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

func (a *authorizationClient) FindUserRoles(userId types.Id) ([]entity.Role, error) {
	return nil, nil
}

func (a *authorizationClient) HasPermission(userId types.Id) error {
	return nil
}
