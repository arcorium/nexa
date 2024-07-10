package main

import (
  "context"
  sharedConf "github.com/arcorium/nexa/shared/config"
  "github.com/arcorium/nexa/shared/database"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/uptrace/bun"
  "log"
  "nexa/services/authorization/constant"
  "nexa/services/authorization/internal/infra/repository/model"
  "time"
)

func main() {
  dbConfig, err := sharedConf.Load[sharedConf.PostgresDatabase]()
  if err != nil {
    log.Fatalln(err)
  }

  db, err := database.OpenPostgresWithConfig(dbConfig, true)
  //db, err := database.OpenPostgres("postgres://postgres:password@localhost:5432/postgres", false, 0, true)
  //db, err := database.OpenPostgres("postgres://nexa:nexa@db/nexa-authz", false, 0, true)
  if err != nil {
    log.Fatalln(err)
  }
  defer db.Close()

  model.RegisterBunModels(db)

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
    err = database.Seed(tx, seedPerms...)
    if err != nil {
      return err
    }

    // Append permissions to role
    err = database.Seed(tx, rolePerms...)
    return err
  })

  if err != nil {
    log.Fatalln(err)
  }

  log.Println("Success seed database: ", dbConfig.DSN())
}
