package main

import (
  "context"
  "errors"
  "fmt"
  authZv1 "github.com/arcorium/nexa/proto/gen/go/authorization/v1"
  sharedConst "github.com/arcorium/nexa/shared/constant"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  "google.golang.org/grpc/metadata"
  "log"
  "nexa/services/relation/constant"
  "os"
  "sync"
  "time"
)

func getData() ([]*authZv1.CreatePermissionRequest, []*authZv1.CreatePermissionRequest) {
  // Seed role
  superPerms := sharedUtil.CastSlice(RELATION_SUPER_PERMS, func(action types.Action) *authZv1.CreatePermissionRequest {
    return &authZv1.CreatePermissionRequest{
      Resource: constant.SERVICE_RESOURCE,
      Action:   action.String(),
    }
  })

  defaultPerms := sharedUtil.CastSlice(RELATION_DEFAULT_PERMS, func(action types.Action) *authZv1.CreatePermissionRequest {
    return &authZv1.CreatePermissionRequest{
      Resource: constant.SERVICE_RESOURCE,
      Action:   action.String(),
    }
  })

  return superPerms, defaultPerms
}

func getConfig() (string, string, error) {
  authz, ok := os.LookupEnv("AUTHZ_SERVICE_ADDR")
  if !ok {
    return "", "", errors.New("AUTHZ_SERVICE_ADDR environment variable not set")
  }

  token, ok := os.LookupEnv("TEMP_TOKEN")
  if !ok {
    return "", "", errors.New("TEMP_TOKEN environment variable not set")
  }

  return authz, token, nil
}

func seedPerms(permClient authZv1.PermissionServiceClient, roleClient authZv1.RoleServiceClient, md metadata.MD, role authZv1.DefaultRole, perms []*authZv1.CreatePermissionRequest) {
  // Seed permissions
  var permIds []string
  for {
    // Try until success
    mdCtx := metadata.NewOutgoingContext(context.Background(), md)
    resp, err := permClient.Seed(mdCtx, &authZv1.SeedPermissionRequest{Permissions: perms})
    if err != nil {
      log.Printf("failed to create permission: %s", err)
      time.Sleep(1000)
      continue
    }

    permIds = resp.PermissionIds
    break
  }
  log.Println("Succeed seed super role permissions")

  for {
    mdCtx := metadata.NewOutgoingContext(context.Background(), md)
    _, err := roleClient.AppendDefaultRolePermissions(mdCtx, &authZv1.AppendDefaultRolePermissionsRequest{
      Role:          role,
      PermissionIds: permIds,
    })
    if err != nil {
      log.Printf("failed to append super admin role permission: %s", err)
      time.Sleep(1000)
      continue
    }
    break
  }

  log.Println("Succeed append super role permissions")
}

func main() {
  var err error
  var conn *grpc.ClientConn

  authz, token, err := getConfig()
  if err != nil {
    log.Fatalln(err)
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
  md := metadata.New(map[string]string{
    "authorization": fmt.Sprintf("%s %s", sharedConst.DEFAULT_ACCESS_TOKEN_SCHEME, token),
  })
  permClient := authZv1.NewPermissionServiceClient(conn)
  roleClient := authZv1.NewRoleServiceClient(conn)

  superPerms, defaultPerms := getData()

  wg := sync.WaitGroup{}
  wg.Add(2)
  go func() {
    defer wg.Done()
    seedPerms(permClient, roleClient, md, authZv1.DefaultRole_DEFAULT_ROLE, defaultPerms)
  }()
  go func() {
    defer wg.Done()
    seedPerms(permClient, roleClient, md, authZv1.DefaultRole_SUPER_ROLE, superPerms)
  }()
  wg.Wait()
}
