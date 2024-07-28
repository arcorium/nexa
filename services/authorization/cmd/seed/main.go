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

func seedData() ([]model.Role, []model.Permission, []model.RolePermission) {
  // Seed role
  superPerms := sharedUtil.MapToSlice(AUTHZ_SUPER_PERMS, func(action types.Action, code string) model.Permission {
    return model.Permission{
      Id:        types.MustCreateId().String(),
      Resource:  constant.SERVICE_RESOURCE,
      Action:    action.String(),
      CreatedAt: time.Now(),
    }
  })

  defaultPerms := sharedUtil.MapToSlice(AUTHZ_DEFAULT_PERMS, func(action types.Action, code string) model.Permission {
    return model.Permission{
      Id:        types.MustCreateId().String(),
      Resource:  constant.SERVICE_RESOURCE,
      Action:    action.String(),
      CreatedAt: time.Now(),
    }
  })

  perms := append(superPerms, defaultPerms...)

  superDesc := "Role that capable to do anything"
  defaultDesc := "Default role"

  superRole := model.Role{
    Id:          types.MustCreateId().String(),
    Name:        constant.SUPER_ROLE_NAME,
    Description: &superDesc,
    CreatedAt:   time.Now(),
  }

  defaultRole := model.Role{
    Id:          types.MustCreateId().String(),
    Name:        constant.DEFAULT_ROLE_NAME,
    Description: &defaultDesc,
    CreatedAt:   time.Now(),
  }

  roles := []model.Role{defaultRole, superRole}

  // Default role permissions
  defaultRolePerms := sharedUtil.CastSliceP(defaultPerms, func(perm *model.Permission) model.RolePermission {
    return model.RolePermission{
      RoleId:       defaultRole.Id,
      PermissionId: perm.Id,
      CreatedAt:    time.Now(),
    }
  })

  // Super role permissions
  superRolePerms := sharedUtil.CastSliceP(perms, func(perm *model.Permission) model.RolePermission {
    return model.RolePermission{
      RoleId:       superRole.Id,
      PermissionId: perm.Id,
      CreatedAt:    time.Now(),
    }
  })

  rolePerms := append(superRolePerms, defaultRolePerms...)

  return roles, perms, rolePerms
}

func main() {
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

  roles, perms, rolePerms := seedData()

  err = db.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
    // Seed role
    err = database.Seed(tx, roles...)
    if err != nil {
      return err
    }

    // Seed permissions
    err = database.Seed(tx, perms...)
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
