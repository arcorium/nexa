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
      Resource: constant.SERVICE_RESOURCE,
      Action:   action,
    }
  })

  ctx := context.Background()
  client := authZv1.NewPermissionServiceClient(conn)

  wg := sync.WaitGroup{}

  for i := 0; i < len(permissions); i++ {
    wg.Add(1)
    go func() {
      defer wg.Done()

      for {
        // Try until success
        _, err := client.Create(ctx, permissions[i])
        if err != nil {
          logger.Warnf("failed to create permission: %s", err)
          continue
        }
        break
      }
    }()
  }

  wg.Wait()
}
