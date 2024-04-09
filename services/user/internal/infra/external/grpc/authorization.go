package grpc

import (
  "context"
  "google.golang.org/grpc"
  authProto "nexa/proto/generated/golang/authorization/v1"
  "nexa/services/user/internal/domain/external"
  "nexa/shared/types"
)

func NewAuthorizationClient(conn grpc.ClientConnInterface) external.IAuthorizationClient {
  return &authorizationClient{
    authClient: authProto.NewAuthorizationClient(conn),
  }
}

type authorizationClient struct {
  authClient authProto.AuthorizationClient
}

func (a *authorizationClient) HasPermission(ctx context.Context, userId types.Id, resource string, actions ...string) error {
  //dto := authProto.CheckUserRequest{
  //  UserId:      userId.Underlying().String(),
  //  Permissions: mapper.ToInternalCheckUserInput(resource, actions...),
  //}
  //_, err := a.authClient.CheckUserPermission(ctx, &dto)
  //return err
  return nil
}
