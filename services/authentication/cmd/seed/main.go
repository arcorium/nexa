package main

import (
  "context"
  "fmt"
  authZv1 "github.com/arcorium/nexa/proto/gen/go/authorization/v1"
  sharedConst "github.com/arcorium/nexa/shared/constant"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  "google.golang.org/grpc/metadata"
  "log"
  "nexa/services/authentication/constant"
  "os"
)

func main() {
  var err error
  var conn *grpc.ClientConn

  authz, ok := os.LookupEnv("AUTHZ_SERVICE_ADDR")
  if !ok {
    log.Fatalln("AUTHZ_SERVICE_ADDR environment variable not set")
  }

  token, ok := os.LookupEnv("TEMP_TOKEN")
  if !ok {
    log.Fatalln("TEMP_TOKEN environment variable not set")
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

  //ctx := context.Background()
  md := metadata.New(map[string]string{
    "authorization": fmt.Sprintf("%s %s", sharedConst.DEFAULT_ACCESS_TOKEN_SCHEME, token),
  })
  client := authZv1.NewPermissionServiceClient(conn)
  roleClient := authZv1.NewRoleServiceClient(conn)

  var permIds []string
  for {
    // Try until success
    mdCtx := metadata.NewOutgoingContext(context.Background(), md)
    resp, err := client.Seed(mdCtx, &authZv1.SeedPermissionRequest{Permissions: permissions})
    if err != nil {
      log.Printf("failed to create permission: %s", err)
      continue
    }

    permIds = resp.PermissionIds
    break
  }

  log.Println("Succeed seed permissions: ", authz)

  // Append it to super roles
  for {
    mdCtx := metadata.NewOutgoingContext(context.Background(), md)
    _, err := roleClient.AppendSuperRolePermissions(mdCtx, &authZv1.AppendSuperRolePermissionsRequest{
      PermissionIds: permIds,
    })
    if err != nil {
      log.Printf("failed to append super admin role permission: %s", err)
      continue
    }

    break
  }

  log.Println("Succeed append super role permissions: ", authz)

}
