package main

import (
  "context"
  "database/sql"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  "log"
  authZv1 "nexa/proto/gen/go/authorization/v1"
  "nexa/services/user/config"
  "nexa/services/user/internal/infra/repository/model"
  sharedConf "nexa/shared/config"
  "nexa/shared/database"
  "nexa/shared/env"
  "nexa/shared/logger"
  "nexa/shared/types"
  "time"
)

func main() {
  envName := ".env"
  if config.IsDebug() {
    envName = "dev.env"
  }

  if err := env.LoadEnvs(envName); err != nil {
    log.Println(err)
  }

  dbConfig, err := sharedConf.LoadDatabase()
  if err != nil {
    log.Fatalln(err)
  }

  db, err := database.OpenPostgres(dbConfig, true)
  if err != nil {
    log.Fatalln(err)
  }
  defer db.Close()

  model.RegisterBunModels(db)

  if err = model.CreateTables(db); err != nil {
    log.Fatalln(err)
  }

  // Seed super user
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

  // Set super role
  var conn grpc.ClientConnInterface
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

  client := authZv1.NewRoleServiceClient(conn)
  for {
    _, err = client.SetAsSuper(context.Background(), &authZv1.SetAsSuperRequest{
      UserId: user.Id,
    })
    if err != nil {
      logger.Warnf("failed to set as super: %v", err)
      logger.Info("Trying again")
      continue
    }
    break
  }
}
