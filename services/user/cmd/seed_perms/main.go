package main

import (
  "context"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  authZv1 "nexa/proto/gen/go/authorization/v1"
  "nexa/services/user/constant"
  "nexa/shared/logger"
  sharedUtil "nexa/shared/util"
  "sync"
)

func main() {
  var err error
  var conn *grpc.ClientConn

  for {
    option := grpc.WithTransportCredentials(insecure.NewCredentials())
    conn, err = grpc.NewClient("", option)
    if err != nil {
      logger.Warnf("failed to connect to grpc server: %v", err)
      logger.Info("Trying to connect again")
      continue
    }
    break
  }

  permissions := sharedUtil.MapToSlice(constant.USER_PERMISSIONS, func(action, perm string) *authZv1.CreatePermissionRequest {
    return &authZv1.CreatePermissionRequest{
      Resource: constant.USER_SERVICE_RESOURCE,
      Action:   action,
    }
  })

  ctx := context.Background()
  permClient := authZv1.NewPermissionServiceClient(conn)
  roleClient := authZv1.NewRoleServiceClient(conn)

  wg := sync.WaitGroup{}

  var permIds []string
  for i := 0; i < len(permissions); i++ {
    wg.Add(1)
    go func() {
      defer wg.Done()

      // Create permission
      for {
        // Try until success
        resp, err := permClient.Create(ctx, permissions[i])
        if err != nil {
          logger.Warnf("failed to create permission: %s", err)
          continue
        }
        permIds = append(permIds, resp.PermissionId)
        break
      }

    }()
  }

  wg.Wait()

  // Append it to super roles
  for {
    _, err := roleClient.AppendSuperRolePermissions(ctx, &authZv1.AppendSuperRolePermissionsRequest{
      PermissionIds: permIds,
    })
    if err != nil {
      logger.Warnf("failed to append super admin role permission: %s", err)
      continue
    }

    break
  }
}
