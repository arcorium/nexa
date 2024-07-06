package pg

import (
  "context"
  "github.com/brianvoe/gofakeit/v7"
  "github.com/stretchr/testify/require"
  "github.com/stretchr/testify/suite"
  "github.com/testcontainers/testcontainers-go"
  "github.com/testcontainers/testcontainers-go/modules/postgres"
  "github.com/testcontainers/testcontainers-go/wait"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/trace/noop"
  "nexa/services/user/internal/domain/entity"
  "nexa/services/user/internal/infra/repository/model"
  sharedConf "nexa/shared/config"
  "nexa/shared/database"
  "nexa/shared/types"
  "nexa/shared/util"
  "nexa/shared/util/repo"
  "reflect"
  "strconv"
  "testing"
  "time"
)

const (
  USER_DB_USERNAME = "user"
  USER_DB_PASSWORD = "password"
  USER_DB          = "nexa"

  SEED_USER_DATA_SIZE = 5
)

var userSeed []entity.User

func ignoreUserFieldsP(got *entity.User) {
  got.BannedUntil = util.RoundTimeToSecond(got.BannedUntil)
  got.CreatedAt = util.RoundTimeToSecond(got.CreatedAt)
  got.DeletedAt = util.RoundTimeToSecond(got.DeletedAt)
  got.CreatedAt = time.Time{}
  got.Password = "" // Same characters could make different hash
  got.Profile = nil
}

func ignoreUserFields(got ...entity.User) {
  // Ignore time fields
  for i := 0; i < len(got); i += 1 {
    ignoreUserFieldsP(&got[i])
  }
}

type userTestSuite struct {
  suite.Suite
  container *postgres.PostgresContainer
  db        bun.IDB
  tracer    trace.Tracer // Mock
}

func (f *userTestSuite) SetupSuite() {
  ctx := context.Background()

  container, err := postgres.RunContainer(ctx,
    postgres.WithUsername(USER_DB_USERNAME),
    postgres.WithPassword(USER_DB_PASSWORD),
    postgres.WithDatabase(USER_DB),
    testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
      WithOccurrence(2).
      WithStartupTimeout(5*time.Second)),
  )
  f.Require().NoError(err)
  f.container = container

  inspect, err := container.Inspect(ctx)
  f.Require().NoError(err)
  ports := inspect.NetworkSettings.Ports
  mapped := ports["5432/tcp"]

  db, err := database.OpenPostgres(&sharedConf.Database{
    Protocol: "postgres",
    Host:     types.Must(container.Host(ctx)),
    Port:     uint16(types.Must(strconv.Atoi(mapped[0].HostPort))),
    Username: USER_DB_USERNAME,
    Password: USER_DB_PASSWORD,
    Name:     USER_DB,
    IsSecure: false,
    Timeout:  time.Second * 10,
  }, false)
  f.Require().NoError(err)
  f.db = db

  // Tracer
  provider := noop.NewTracerProvider()
  f.tracer = provider.Tracer("MOCK")

  // Seed
  model.RegisterBunModels(db)
  err = model.CreateTables(db)
  f.Require().NoError(err)
  // Seeding
  users := util.CastSliceP(userSeed, func(from *entity.User) model.User {
    return model.FromUserDomain(from, func(ent *entity.User, profile *model.User) {
    })
  })
  profiles := util.CastSliceP(profileSeed, func(from *entity.Profile) model.Profile {
    return model.FromProfileDomain(from, func(ent *entity.Profile, profile *model.Profile) {
    })
  })
  err = database.Seed(f.db, users...)
  f.Require().NoError(err)
  err = database.Seed(f.db, profiles...)
  f.Require().NoError(err)
}

func (f *userTestSuite) TearDownSuite() {
  err := f.container.Terminate(context.Background())
  f.Require().NoError(err)
}

func (f *userTestSuite) Test_userRepository_Create() {
  type args struct {
    ctx  context.Context
    user *entity.User
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        ctx:  context.Background(),
        user: generateRandomUserP(),
      },
      wantErr: false,
    },
    {
      name: "Duplicate Username",
      args: args{
        ctx: context.Background(),
        user: util.CopyWithP(generateRandomUser(), func(ent *entity.User) {
          ent.Username = userSeed[0].Username
        }),
      },
      wantErr: true,
    },
    {
      name: "Duplicate Email",
      args: args{
        ctx: context.Background(),
        user: util.CopyWithP(generateRandomUser(), func(ent *entity.User) {
          ent.Email = userSeed[0].Email
        }),
      },
      wantErr: true,
    },
    {
      name: "Empty Non-null field",
      args: args{
        ctx: context.Background(),
        user: util.CopyWithP(generateRandomUser(), func(ent *entity.User) {
          ent.Email = ""
          ent.Password = ""
        }),
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      u := userRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      if err := u.Create(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
        t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func (f *userTestSuite) Test_userRepository_Delete() {
  type args struct {
    ctx context.Context
    id  types.Id
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        ctx: context.Background(),
        id:  userSeed[0].Id,
      },
      wantErr: false,
    },
    {
      name: "User not found",
      args: args{
        ctx: context.Background(),
        id:  types.MustCreateId(),
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      u := userRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      err = u.Delete(tt.args.ctx, tt.args.id)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      users, err := u.FindByIds(tt.args.ctx, tt.args.id)
      require.NoError(t, err)
      require.Len(t, users, 1)

      if users[0].DeletedAt.IsZero() {
        t.Errorf("Delete() failed to change is_delete field, obj = %v", users[0])
      }
    })
  }
}

func (f *userTestSuite) Test_userRepository_FindAllUsers() {
  type args struct {
    ctx   context.Context
    query repo.QueryParameter
  }
  tests := []struct {
    name    string
    args    args
    want    repo.PaginatedResult[entity.User]
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: 0,
          Limit:  0,
        },
      },
      want: repo.PaginatedResult[entity.User]{
        Data:    userSeed,
        Total:   SEED_USER_DATA_SIZE,
        Element: SEED_USER_DATA_SIZE,
      },
      wantErr: false,
    },
    {
      name: "Use offset and limit",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: 1,
          Limit:  2,
        },
      },
      want: repo.PaginatedResult[entity.User]{
        Data:    userSeed[1:3],
        Total:   SEED_USER_DATA_SIZE,
        Element: 2,
      },
      wantErr: false,
    },
    {
      name: "Outside users count",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: 2,
          Limit:  5,
        },
      },
      want: repo.PaginatedResult[entity.User]{
        Data:    userSeed[2:],
        Total:   SEED_USER_DATA_SIZE,
        Element: 3,
      },
      wantErr: false,
    },
    {
      name: "Outside users offset",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: 5,
          Limit:  1,
        },
      },
      want: repo.PaginatedResult[entity.User]{
        Data:    nil,
        Total:   SEED_USER_DATA_SIZE,
        Element: 0,
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      u := userRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      got, err := u.Get(tt.args.ctx, tt.args.query)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      ignoreUserFields(got.Data...)
      ignoreUserFields(tt.want.Data...)

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("Get() \ngot = %v\nwant %v", got, tt.want)
      }
    })
  }
}

func (f *userTestSuite) Test_userRepository_FindByEmails() {
  type args struct {
    ctx    context.Context
    emails []types.Email
  }
  tests := []struct {
    name    string
    args    args
    want    []entity.User
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        ctx: context.Background(),
        emails: []types.Email{
          userSeed[0].Email,
        },
      },
      want:    userSeed[:1],
      wantErr: false,
    },
    {
      name: "Multiple emails with bad emails on it",
      args: args{
        ctx: context.Background(),
        emails: []types.Email{
          userSeed[0].Email,
          types.Email(gofakeit.Email()),
          userSeed[1].Email,
        },
      },
      want:    userSeed[0:2],
      wantErr: false,
    },
    {
      name: "User not found",
      args: args{
        ctx: context.Background(),
        emails: []types.Email{
          types.Email(gofakeit.Email()),
          types.Email(gofakeit.Email()),
        },
      },
      want:    nil,
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      u := userRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      got, err := u.FindByEmails(tt.args.ctx, tt.args.emails...)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("FindByEmails() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      f.Require().Len(got, len(tt.want))
      ignoreUserFields(got...)
      ignoreUserFields(tt.want...)

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindByEmails() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *userTestSuite) Test_userRepository_FindByIds() {
  type args struct {
    ctx     context.Context
    userIds []types.Id
  }
  tests := []struct {
    name    string
    args    args
    want    []entity.User
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        ctx: context.Background(),
        userIds: []types.Id{
          userSeed[0].Id,
        },
      },
      want:    userSeed[:1],
      wantErr: false,
    },
    {
      name: "Multiple with some bad id on it",
      args: args{
        ctx: context.Background(),
        userIds: []types.Id{
          userSeed[0].Id,
          types.MustCreateId(),
          userSeed[1].Id,
        },
      },
      want:    userSeed[:2],
      wantErr: false,
    },
    {
      name: "User Not Found",
      args: args{
        ctx: context.Background(),
        userIds: []types.Id{
          types.MustCreateId(),
          types.MustCreateId(),
        },
      },
      want:    nil,
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      u := userRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      got, err := u.FindByIds(tt.args.ctx, tt.args.userIds...)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("FindByIds() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      ignoreUserFields(got...)
      ignoreUserFields(tt.want...)

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindByIds() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *userTestSuite) Test_userRepository_Patch() {
  verifiedInv := !userSeed[0].IsVerified

  type args struct {
    ctx     context.Context
    user    *entity.PatchedUser
    BaseIdx int
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Patch all fields but time fields",
      args: args{
        ctx: context.Background(),
        user: &entity.PatchedUser{
          Id:         userSeed[0].Id,
          Username:   gofakeit.Username(),
          Email:      types.Email(gofakeit.Email()),
          Password:   types.Must(types.Password(gofakeit.Username()).Hash()),
          IsVerified: types.SomeNullable(verifiedInv),
        },
        BaseIdx: 0,
      },
      wantErr: false,
    },
    {
      name: "Patch several fields but time fields",
      args: args{
        ctx: context.Background(),
        user: &entity.PatchedUser{
          Id:       userSeed[0].Id,
          Username: gofakeit.Username(),
          Email:    types.Email(gofakeit.Email()),
        },
        BaseIdx: 0,
      },
      wantErr: false,
    },
    {
      name: "Change verify field",
      args: args{
        ctx: context.Background(),
        user: &entity.PatchedUser{
          Id:         userSeed[0].Id,
          IsVerified: types.SomeNullable(verifiedInv),
        },
        BaseIdx: 0,
      },
      wantErr: false,
    },
    {
      name: "Banned User",
      args: args{
        ctx: context.Background(),
        user: &entity.PatchedUser{
          Id:             userSeed[0].Id,
          BannedDuration: types.SomeNullable(time.Hour * 10),
        },
        BaseIdx: 0,
      },
      wantErr: false,
    },
    {
      name: "Unban User",
      args: args{
        ctx: context.Background(),
        user: &entity.PatchedUser{
          Id:             userSeed[3].Id,
          BannedDuration: types.SomeNullable(time.Duration(0)),
        },
        BaseIdx: 3,
      },
      wantErr: false,
    },
    {
      name: "Delete User",
      args: args{
        ctx: context.Background(),
        user: &entity.PatchedUser{
          Id:       userSeed[0].Id,
          IsDelete: types.SomeNullable(true),
        },
        BaseIdx: 0,
      },
      wantErr: false,
    },
    {
      name: "Undelete User",
      args: args{
        ctx: context.Background(),
        user: &entity.PatchedUser{
          Id:       userSeed[4].Id,
          IsDelete: types.SomeNullable(false),
        },
        BaseIdx: 4,
      },
      wantErr: false,
    },
    {
      name: "User not found",
      args: args{
        ctx: context.Background(),
        user: &entity.PatchedUser{
          Id:       types.MustCreateId(),
          Username: "arcorium",
        },
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      u := userRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      err = u.Patch(tt.args.ctx, tt.args.user)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("Patch() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      got, err := u.FindByIds(tt.args.ctx, tt.args.user.Id)
      require.NoError(t, err)
      require.Len(t, got, 1)

      // Use data from base id and change patched field
      comparator := userSeed[tt.args.BaseIdx]
      if tt.args.user.Id != types.NullId() {
        comparator.Id = tt.args.user.Id
      }
      if tt.args.user.Username != "" {
        comparator.Username = tt.args.user.Username
      }
      if tt.args.user.Email != "" {
        comparator.Email = tt.args.user.Email
      }
      if tt.args.user.Password != "" {
        comparator.Password = tt.args.user.Password
      }
      if tt.args.user.IsVerified.HasValue() {
        comparator.IsVerified = tt.args.user.IsVerified.RawValue()
      }

      if tt.args.user.IsDelete.HasValue() {
        if tt.args.user.IsDelete.RawValue() {
          // Delete
          comparator.DeletedAt = time.Now()
          // Check
          if !got[0].DeletedAt.Before(comparator.DeletedAt) {
            t.Errorf("Patch() DeletedAt Expect Before \ngot = %v, \nwant = %v", got[0].DeletedAt, comparator.DeletedAt)
          }
        } else {
          // Undelete
          comparator.DeletedAt = time.Time{}
          // Check
          if comparator.DeletedAt.IsZero() && !got[0].DeletedAt.IsZero() {
            t.Errorf("Patch() DeletedAt Expect Zero \ngot = %v, \nwant = %v", got[0].DeletedAt, comparator.DeletedAt)
          }

        }
      }

      if tt.args.user.BannedDuration.HasValue() {
        if int(tt.args.user.BannedDuration.RawValue().Seconds()) == 0 {
          // Unban
          comparator.BannedUntil = time.Time{}
          // Check time
          if comparator.BannedUntil.IsZero() && !got[0].BannedUntil.IsZero() {
            t.Errorf("Patch() BannedUntil Expect Zero \ngot = %v, \nwant = %v", got[0].BannedUntil, comparator.BannedUntil)
          }
        } else {
          // Banned
          comparator.BannedUntil = time.Now().Add(tt.args.user.BannedDuration.RawValue())
          // Check time
          if !got[0].BannedUntil.Before(comparator.BannedUntil) {
            t.Errorf("Patch() BannedUntil Expect Before \ngot = %v, \nwant = %v", got[0].BannedUntil, comparator.BannedUntil)
          }
        }
      }

      ignoreUserFields(got...)
      ignoreUserFieldsP(&comparator)

      if !reflect.DeepEqual(got[0], comparator) != tt.wantErr {
        t.Errorf("Patch() \ngot = %v, \nwant = %v", got[0], comparator)
      }
    })
  }
}

func (f *userTestSuite) Test_userRepository_Update() {
  type args struct {
    ctx  context.Context
    user *entity.User
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        ctx: context.Background(),
        user: util.CopyWithP(generateRandomUser(), func(ent *entity.User) {
          ent.Id = userSeed[0].Id
        }),
      },
      wantErr: false,
    },
    {
      name: "User not found",
      args: args{
        ctx: context.Background(),
        user: util.CopyWithP(generateRandomUser(), func(ent *entity.User) {
          ent.Id = types.MustCreateId()
        }),
      },
      wantErr: true,
    },
    {
      name: "Non-null field is empty",
      args: args{
        ctx: context.Background(),
        user: util.CopyWithP(generateRandomUser(), func(ent *entity.User) {
          ent.Username = ""
          ent.Email = ""
          ent.Password = ""
        }),
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      u := userRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      err = u.Update(tt.args.ctx, tt.args.user)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      users, err := u.FindByIds(tt.args.ctx, tt.args.user.Id)
      require.NoError(t, err)
      require.Len(t, users, 1)

      ignoreUserFields(users...)
      ignoreUserFieldsP(tt.args.user)

      if !reflect.DeepEqual(users[0], *tt.args.user) {
        t.Errorf("FindByIds() got = %v, want %v", users[0], tt.args.user)
      }
    })
  }
}

func TestUser(t *testing.T) {
  seedUserData()
  seedProfileData()

  suite.Run(t, &userTestSuite{})
}

func seedUserData() {
  for i := 0; i < SEED_USER_DATA_SIZE; i += 1 {
    userSeed = append(userSeed, generateRandomUser())
  }
  userSeed[3].BannedUntil = time.Now().Add(time.Hour * time.Duration(gofakeit.Hour()+2))
  userSeed[4].DeletedAt = time.Now().Add(time.Hour * time.Duration(gofakeit.Hour()+2) * -1)

}

func generateRandomUser() entity.User {
  return entity.User{
    Id:          types.MustCreateId(),
    Username:    gofakeit.Username(),
    Email:       types.Must(types.EmailFromString(gofakeit.Email())),
    Password:    types.Must(types.PasswordFromString(gofakeit.Username()).Hash()),
    IsVerified:  gofakeit.Bool(),
    DeletedAt:   time.Time{},
    BannedUntil: time.Time{},
    CreatedAt:   time.Now(),
  }
}

func generateRandomUserP() *entity.User {
  user := generateRandomUser()
  return &user
}
