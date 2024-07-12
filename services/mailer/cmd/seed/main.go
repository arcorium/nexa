package main

import (
  "context"
  "fmt"
  authZv1 "github.com/arcorium/nexa/proto/gen/go/authorization/v1"
  sharedConf "github.com/arcorium/nexa/shared/config"
  sharedConst "github.com/arcorium/nexa/shared/constant"
  "github.com/arcorium/nexa/shared/database"
  "github.com/arcorium/nexa/shared/env"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  "google.golang.org/grpc/metadata"
  "log"
  "nexa/services/mailer/config"
  "nexa/services/mailer/constant"
  "nexa/services/mailer/internal/infra/repository/model"
  "os"
  "time"
)

func main() {
  var err error
  var conn *grpc.ClientConn

  envName := ".env"
  if config.IsDebug() {
    envName = "dev.env"
  }
  _ = env.LoadEnvs(envName)

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
      time.Sleep(100)
      continue
    }
    break
  }

  permissions := sharedUtil.MapToSlice(constant.MAILER_PERMISSIONS, func(action string, perms string) *authZv1.CreatePermissionRequest {
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
      time.Sleep(100)
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
