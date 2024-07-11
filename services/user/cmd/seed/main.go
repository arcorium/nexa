package main

import (
  "context"
  "database/sql"
  authZv1 "github.com/arcorium/nexa/proto/gen/go/authorization/v1"
  sharedConf "github.com/arcorium/nexa/shared/config"
  "github.com/arcorium/nexa/shared/database"
  "github.com/arcorium/nexa/shared/env"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  "log"
  "nexa/services/user/config"
  "nexa/services/user/constant"
  "nexa/services/user/internal/infra/repository/model"
  "os"
  "sync"
  "time"
)

func seedDatabase() model.User {
  // Connect database
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

  //Seed super user
  tx, err := db.BeginTx(context.Background(), nil)
  if err != nil {
    log.Fatalln(err)
  }

  password := types.Password(env.GetDefaulted("NEXA_PASSWORD", "super123"))

  user := model.User{
    Id:         types.MustCreateId().String(),
    Username:   env.GetDefaulted("NEXA_USER_NAME", "super"),
    Email:      env.GetDefaulted("NEXA_EMAIL", "super@nexa.com"),
    Password:   types.Must(password.Hash()).String(),
    IsVerified: sql.NullBool{Bool: true, Valid: true},
    CreatedAt:  time.Now(),
  }

  profile := model.Profile{
    Id:        types.MustCreateId().String(),
    UserId:    user.Id,
    FirstName: env.GetDefaulted("NEXA_FIRST_NAME", "nexa"),
    LastName:  env.GetDefaultedP("NEXA_LAST_NAME", "super"),
  }

  // Seed user
  err = database.Seed(tx, user)
  if err != nil {
    err := tx.Rollback()
    if err != nil {
      log.Fatalln(err)
    }
    log.Fatalln(err)
  }

  // Seed profile
  err = database.Seed(tx, profile)
  if err != nil {
    err := tx.Rollback()
    if err != nil {
      log.Fatalln(err)
    }
    log.Fatalln(err)
  }

  err = tx.Commit()
  if err != nil {
    log.Fatalln("Failed to commit transaction:", err)
  }

  log.Println("Succeed seed database: ", dbConfig.DSN())
  return user
}

func main() {
  var err error
  var conn *grpc.ClientConn

  envName := ".env"
  if config.IsDebug() {
    envName = "dev.env"
  }
  _ = env.LoadEnvs(envName)

  addr, ok := os.LookupEnv("NEXA_AUTHZ_SERVICE_ADDRESS")
  if !ok {
    log.Fatalln("NEXA_AUTHZ_SERVICE_ADDRESS environment variable not set")
  }

  // Connect authorization service client
  for {
    option := grpc.WithTransportCredentials(insecure.NewCredentials())
    conn, err = grpc.NewClient(addr, option)
    if err != nil {
      log.Printf("failed to connect to grpc server: %v", err)
      log.Printf("Trying to connect again")
      continue
    }
    break
  }

  user := seedDatabase()

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
          log.Printf("failed to create permission: %s", err)
          continue
        }
        permIds = append(permIds, resp.PermissionId)
        break
      }

    }()
  }

  wg.Wait()

  log.Println("Succeed seed permissions: ", addr)

  // Append it to super roles
  for {
    _, err := roleClient.AppendSuperRolePermissions(ctx, &authZv1.AppendSuperRolePermissionsRequest{
      PermissionIds: permIds,
    })
    if err != nil {
      log.Printf("failed to append super admin role permission: %s", err)
      continue
    }

    break
  }

  log.Println("Succeed append super role permissions: ", addr)

  // Set super role
  for {
    _, err = roleClient.SetAsSuper(context.Background(), &authZv1.SetAsSuperRequest{
      UserId: user.Id,
    })
    if err != nil {
      log.Printf("failed to set as super: %v", err)
      log.Printf("Trying again")
      continue
    }
    break
  }

  log.Printf("Succeed set user %v as super role", user)
}
