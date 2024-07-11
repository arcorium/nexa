package main

import (
  "context"
  authZv1 "github.com/arcorium/nexa/proto/gen/go/authorization/v1"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  "log"
  "nexa/services/authentication/constant"
  "os"
  "sync"
)

func main() {
  var err error
  var conn *grpc.ClientConn

  authz, ok := os.LookupEnv("AUTHZ_SERVICE_ADDR")
  if !ok {
    log.Fatalln("AUTHZ_SERVICE_ADDR environment variable not set")
  }

  for {
    option := grpc.WithTransportCredentials(insecure.NewCredentials())
    conn, err = grpc.NewClient(authz, option)
    if err != nil {
      log.Printf("failed to connect to grpc server: %v", err)
      log.Printf("Trying to connect again")
      continue
    }
    break
  }

  permissions := sharedUtil.MapToSlice(constant.AUTHN_PERMISSIONS, func(action, perms string) *authZv1.CreatePermissionRequest {
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
          log.Printf("failed to create permission: %s", err)
          continue
        }
        break
      }
    }()
  }

  wg.Wait()

  log.Println("Succeed seed permissions: ", authz)
}
