package main

import (
  "context"
  "errors"
  "fmt"
  authZv1 "github.com/arcorium/nexa/proto/gen/go/authorization/v1"
  sharedConf "github.com/arcorium/nexa/shared/config"
  sharedConst "github.com/arcorium/nexa/shared/constant"
  "github.com/arcorium/nexa/shared/database"
  "github.com/arcorium/nexa/shared/env"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  "google.golang.org/grpc/metadata"
  "log"
  "nexa/services/mailer/config"
  "nexa/services/mailer/constant"
  "nexa/services/mailer/internal/infra/repository/model"
  "os"
  "sync"
  "time"
)

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

func seedDatabase() {
  dbConfig, err := sharedConf.Load[sharedConf.PostgresDatabase]()
  if err != nil {
    log.Fatalln(err)
  }

  db, err := database.OpenPostgresWithConfig(dbConfig, true)
  if err != nil {
    log.Fatalln(err)
  }
  defer db.Close()

  model.RegisterBunModels(db)

  // Seed data
  if err = model.SeedDefaultTags(db); err != nil {
    log.Fatalln(err)
  }

  log.Println("Succeed seed database: ", dbConfig.DSN())
}

func getData() ([]*authZv1.CreatePermissionRequest, []*authZv1.CreatePermissionRequest) {
  // Seed role
  superPerms := sharedUtil.CastSlice(MAILER_SUPER_PERMS, func(action types.Action) *authZv1.CreatePermissionRequest {
    return &authZv1.CreatePermissionRequest{
      Resource: constant.SERVICE_RESOURCE,
      Action:   action.String(),
    }
  })

  defaultPerms := sharedUtil.CastSlice(MAILER_DEFAULT_PERMS, func(action types.Action) *authZv1.CreatePermissionRequest {
    return &authZv1.CreatePermissionRequest{
      Resource: constant.SERVICE_RESOURCE,
      Action:   action.String(),
    }
  })

  return superPerms, defaultPerms
}

func seedPermissions(permClient authZv1.PermissionServiceClient, md metadata.MD, perms []*authZv1.CreatePermissionRequest, role string) []string {
  var permIds []string
  for {
    // Try until success
    mdCtx := metadata.NewOutgoingContext(context.Background(), md)
    resp, err := permClient.Seed(mdCtx, &authZv1.SeedPermissionRequest{Permissions: perms})
    if err != nil {
      log.Printf("failed to create permission into %s role: %s", role, err)
      time.Sleep(1000)
      continue
    }

    permIds = resp.PermissionIds
    break
  }
  log.Printf("Succeed seed %s role permissions", role)
  return permIds
}

func appendPermission(roleClient authZv1.RoleServiceClient, md metadata.MD, role authZv1.DefaultRole, permIds []string) {
  for {
    mdCtx := metadata.NewOutgoingContext(context.Background(), md)
    _, err := roleClient.AppendDefaultRolePermissions(mdCtx, &authZv1.AppendDefaultRolePermissionsRequest{
      Role:          role,
      PermissionIds: permIds,
    })
    if err != nil {
      log.Printf("failed to append %s role permission: %s", role.String(), err)
      time.Sleep(1000)
      continue
    }
    break
  }

  log.Printf("Succeed append %s role permissions", role.String())
}

func seed(conn grpc.ClientConnInterface, token string) {
  md := metadata.New(map[string]string{
    "authorization": fmt.Sprintf("%s %s", sharedConst.DEFAULT_ACCESS_TOKEN_SCHEME, token),
  })
  permClient := authZv1.NewPermissionServiceClient(conn)
  roleClient := authZv1.NewRoleServiceClient(conn)

  superPerms, defaultPerms := getData()
  // Seed permissions
  var defaultPermIds []string
  var superPermIds []string

  wg := sync.WaitGroup{}
  wg.Add(2)
  go func() {
    defer wg.Done()
    defaultPermIds = seedPermissions(permClient, md, defaultPerms, "DEFAULT")
  }()
  go func() {
    defer wg.Done()
    superPermIds = seedPermissions(permClient, md, superPerms, "SUPER")
  }()
  wg.Wait()

  // Append permissions
  wg.Add(2)
  go func() {
    defer wg.Done()
    appendPermission(roleClient, md, authZv1.DefaultRole_DEFAULT_ROLE, defaultPermIds)
  }()
  go func() {
    defer wg.Done()
    ids := append(defaultPermIds, superPermIds...)
    appendPermission(roleClient, md, authZv1.DefaultRole_SUPER_ROLE, ids)
  }()

  wg.Wait()
}

func main() {
  var err error
  var conn *grpc.ClientConn

  envName := ".env"
  if config.IsDebug() {
    envName = "dev.env"
  }
  _ = env.LoadEnvs(envName)

  authz, token, err := getConfig()
  if err != nil {
    log.Fatalln(err)
  }

  seedDatabase()

  for {
    option := grpc.WithTransportCredentials(insecure.NewCredentials())
    conn, err = grpc.NewClient(authz, option)
    if err != nil {
      log.Printf("failed to connect to grpc server: %v", err)
      log.Printf("Trying to connect again")
      time.Sleep(100)
      continue
    }
    break
  }

  seed(conn, token)
}
