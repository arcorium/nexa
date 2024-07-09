package main

import (
  "context"
  "github.com/uptrace/bun"
  "log"
  "nexa/services/authorization/config"
  "nexa/services/authorization/constant"
  "nexa/services/authorization/internal/infra/repository/model"
  sharedConf "nexa/shared/config"
  "nexa/shared/database"
  "nexa/shared/env"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
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

  // Seed role
  seedPerms := sharedUtil.MapToSlice(constant.AUTHZ_PERMISSIONS, func(action string, code string) model.Permission {
    return model.Permission{
      Id:        types.MustCreateId().String(),
      Resource:  constant.SERVICE_RESOURCE,
      Action:    action,
      CreatedAt: time.Now(),
    }
  })

  s := "Default roles that capable to do anything"
  superRole := model.Role{
    Id:          types.MustCreateId().String(),
    Name:        constant.DEFAULT_SUPER_ROLE_NAME,
    Description: &s,
    CreatedAt:   time.Now(),
  }

  rolePerms := sharedUtil.CastSliceP(seedPerms, func(perm *model.Permission) model.RolePermission {
    return model.RolePermission{
      RoleId:       superRole.Id,
      PermissionId: perm.Id,
      CreatedAt:    time.Now(),
    }
  })

  err = db.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
    // Seed role
    err = database.Seed(tx, superRole)
    if err != nil {
      return err
    }

    // Seed permissions
    err = database.Seed(tx, seedPerms)
    if err != nil {
      return err
    }

    // Append permissions to role
    err = database.Seed(tx, rolePerms)
    return err
  })

  if err != nil {
    log.Fatalln(err)
  }
}
