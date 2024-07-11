package main

import (
  "context"
  authZv1 "github.com/arcorium/nexa/proto/gen/go/authorization/v1"
  sharedConf "github.com/arcorium/nexa/shared/config"
  "github.com/arcorium/nexa/shared/database"
  "github.com/arcorium/nexa/shared/env"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  "log"
  "nexa/services/mailer/config"
  "nexa/services/mailer/constant"
  "nexa/services/mailer/internal/infra/repository/model"
  "os"
  "sync"
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

  authz, ok := os.LookupEnv("AUTHZ_SERVICE_ADDR")
  if !ok {
    log.Fatalln("AUTHZ_SERVICE_ADDR environment variable not set")
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

  permissions := sharedUtil.MapToSlice(constant.MAILER_PERMISSIONS, func(action string, perms string) *authZv1.CreatePermissionRequest {
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
