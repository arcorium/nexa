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
  "reflect"
  "strconv"
  "testing"
  "time"
)

const (
  PROFILE_DB_USERNAME = "user"
  PROFILE_DB_PASSWORD = "password"
  PROFILE_DB          = "nexa"

  SEED_PROFILE_DATA_SIZE = 3
)

type profileTestSuite struct {
  suite.Suite
  container *postgres.PostgresContainer
  db        bun.IDB
  tracer    trace.Tracer // Mock

  profileSeed []entity.Profile
  userSeed    []entity.User
}

func (f *profileTestSuite) SetupSuite() {
  // Prepare data
  for i := 0; i < SEED_USER_DATA_SIZE; i += 1 {
    f.userSeed = append(f.userSeed, generateRandomUser())
  }
  f.userSeed[3].BannedUntil = time.Now().Add(time.Hour * time.Duration(gofakeit.Hour()+2))
  f.userSeed[4].DeletedAt = time.Now().Add(time.Hour * time.Duration(gofakeit.Hour()+2) * -1)

  for i := 0; i < SEED_PROFILE_DATA_SIZE; i += 1 {
    f.profileSeed = append(f.profileSeed, generateRandomProfile(f.userSeed[i].Id))
  }

  ctx := context.Background()

  container, err := postgres.RunContainer(ctx,
    postgres.WithUsername(PROFILE_DB_USERNAME),
    postgres.WithPassword(PROFILE_DB_PASSWORD),
    postgres.WithDatabase(PROFILE_DB),
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
    Username: PROFILE_DB_USERNAME,
    Password: PROFILE_DB_PASSWORD,
    Name:     PROFILE_DB,
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
  users := util.CastSliceP(f.userSeed, func(from *entity.User) model.User {
    return model.FromUserDomain(from, func(ent *entity.User, profile *model.User) {
    })
  })
  profiles := util.CastSliceP(f.profileSeed, func(from *entity.Profile) model.Profile {
    return model.FromProfileDomain(from, func(ent *entity.Profile, profile *model.Profile) {
    })
  })
  err = database.Seed(f.db, users...)
  f.Require().NoError(err)
  err = database.Seed(f.db, profiles...)
  f.Require().NoError(err)
}

func (f *profileTestSuite) TearDownSuite() {
  err := f.container.Terminate(context.Background())
  f.Require().NoError(err)
}

func (f *profileTestSuite) Test_profileRepository_Create() {
  type args struct {
    ctx     context.Context
    profile *entity.Profile
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        ctx:     context.Background(),
        profile: generateRandomProfileP(f.userSeed[3].Id),
      },
      wantErr: false,
    },
    {
      name: "Duplicate User Id",
      args: args{
        ctx: context.Background(),
        profile: util.CopyWithP(generateRandomProfile(types.MustCreateId()), func(e *entity.Profile) {
          e.Id = f.userSeed[0].Id // Use same id
        }),
      },
      wantErr: true,
    },
    {
      name: "User not found",
      args: args{
        ctx:     context.Background(),
        profile: generateRandomProfileP(types.MustCreateId()),
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      t := f.T()

      p := profileRepository{
        db:     tx,
        tracer: f.tracer,
      }

      err = p.Create(tt.args.ctx, tt.args.profile)
      if (err != nil) != tt.wantErr {
        t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
      }

      if err != nil {
        return
      }

      profiles, err := p.FindByIds(tt.args.ctx, tt.args.profile.Id)
      require.NoError(t, err)

      if !reflect.DeepEqual(profiles[0], *tt.args.profile) {
        t.Errorf("Create() got = %v, want %v", profiles[0], *tt.args.profile)
      }
    })
  }
}

func (f *profileTestSuite) Test_profileRepository_FindByIds() {
  type args struct {
    ctx     context.Context
    userIds []types.Id
  }
  tests := []struct {
    name    string
    args    args
    want    []entity.Profile
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        ctx:     context.Background(),
        userIds: []types.Id{f.profileSeed[0].Id},
      },
      want:    f.profileSeed[:1],
      wantErr: false,
    },
    {
      name: "Some id is invalid",
      args: args{
        ctx:     context.Background(),
        userIds: []types.Id{f.profileSeed[0].Id, types.MustCreateId(), f.profileSeed[1].Id},
      },
      want:    []entity.Profile{f.profileSeed[0], f.profileSeed[1]},
      wantErr: false,
    },
    {
      name: "User Not Found",
      args: args{
        ctx:     context.Background(),
        userIds: []types.Id{types.MustCreateId()},
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

      t := f.T()

      p := profileRepository{
        db:     tx,
        tracer: f.tracer,
      }
      got, err := p.FindByIds(tt.args.ctx, tt.args.userIds...)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("FindByIds() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      comparatorFunc := func(e *entity.Profile, e2 *entity.Profile) bool {
        return e.Id == e2.Id
      }

      if !util.ArbitraryCheck(got, tt.want, comparatorFunc) != tt.wantErr {
        t.Errorf("FindByIds() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *profileTestSuite) Test_profileRepository_Patch() {
  type args struct {
    ctx     context.Context
    profile *entity.PatchedProfile
    baseIdx int
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Update All fields",
      args: args{
        ctx: context.Background(),
        profile: &entity.PatchedProfile{
          Id:        f.profileSeed[0].Id,
          FirstName: gofakeit.FirstName(),
          LastName:  types.SomeNullable(gofakeit.LastName()),
          Bio:       types.SomeNullable(gofakeit.LoremIpsumParagraph(1, 2, 20, ".")),
          PhotoId:   types.SomeNullable(types.MustCreateId()),
          PhotoURL:  types.SomeNullable(types.FilePathFromString(gofakeit.URL())),
        },
        baseIdx: 0,
      },
      wantErr: false,
    },
    {
      name: "Update several data",
      args: args{
        ctx: context.Background(),
        profile: &entity.PatchedProfile{
          Id:        f.profileSeed[0].Id,
          FirstName: "arcorium",
          LastName:  types.SomeNullable("liz"),
        },
        baseIdx: 0,
      },
      wantErr: false,
    },
    {
      name: "Delete bio and last name",
      args: args{
        ctx: context.Background(),
        profile: &entity.PatchedProfile{
          Id:       f.profileSeed[0].Id,
          Bio:      types.SomeNullable(""),
          LastName: types.SomeNullable(""),
        },
        baseIdx: 0,
      },
      wantErr: false,
    },
    {
      name: "Delete bio and last name and change photo data",
      args: args{
        ctx: context.Background(),
        profile: &entity.PatchedProfile{
          Id:       f.profileSeed[0].Id,
          LastName: types.SomeNullable(""),
          Bio:      types.SomeNullable(""),
          PhotoId:  types.SomeNullable(types.MustCreateId()),
          PhotoURL: types.SomeNullable(types.FilePathFromString(gofakeit.URL())),
        },
        baseIdx: 0,
      },
      wantErr: false,
    },
    {
      name: "User not found",
      args: args{
        ctx: context.Background(),
        profile: &entity.PatchedProfile{
          Id:        f.userSeed[SEED_PROFILE_DATA_SIZE].Id,
          FirstName: "arcorium",
        },
        baseIdx: -1,
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      p := profileRepository{
        db:     tx,
        tracer: f.tracer,
      }

      t := f.T()

      err = p.Patch(tt.args.ctx, tt.args.profile)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("Patch() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      profiles, err := p.FindByIds(tt.args.ctx, []types.Id{tt.args.profile.Id}...)
      require.NoError(t, err)
      require.Len(t, profiles, 1)

      comparator := f.profileSeed[tt.args.baseIdx]
      // Set patched data
      if tt.args.profile.FirstName != "" {
        comparator.FirstName = tt.args.profile.FirstName
      }
      types.SetOnNonNull(&comparator.LastName, tt.args.profile.LastName)
      types.SetOnNonNull(&comparator.Bio, tt.args.profile.Bio)
      types.SetOnNonNull(&comparator.PhotoId, tt.args.profile.PhotoId)
      types.SetOnNonNull(&comparator.PhotoURL, tt.args.profile.PhotoURL)

      if !reflect.DeepEqual(profiles[0], comparator) != tt.wantErr {
        t.Errorf("Update() got = %v, want %v", profiles[0], comparator)
      }
    })
  }
}

func (f *profileTestSuite) Test_profileRepository_Update() {
  type args struct {
    ctx     context.Context
    profile *entity.Profile
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
        profile: util.CopyWithP(generateRandomProfile(f.userSeed[0].Id), func(ent *entity.Profile) {
          ent.Id = f.profileSeed[0].Id
        }),
      },
      wantErr: false,
    },
    {
      name: "Profile Not Found",
      args: args{
        ctx:     context.Background(),
        profile: generateRandomProfileP(types.MustCreateId()),
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      p := profileRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      err = p.Update(tt.args.ctx, tt.args.profile)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      profiles, err := p.FindByIds(tt.args.ctx, tt.args.profile.Id)
      require.NoError(t, err)
      require.Len(t, profiles, 1)

      if !reflect.DeepEqual(profiles[0], *tt.args.profile) {
        t.Errorf("Update() got = %v, want %v", profiles[0], *tt.args.profile)
      }
    })
  }
}

func TestProfile(t *testing.T) {
  suite.Run(t, &profileTestSuite{})
}

func generateRandomProfile(userId types.Id) entity.Profile {
  return entity.Profile{
    Id:        types.MustCreateId(),
    UserId:    userId,
    FirstName: gofakeit.FirstName(),
    LastName:  gofakeit.LastName(),
    Bio:       gofakeit.LoremIpsumParagraph(1, 3, gofakeit.Number(20, 40), "."),
    PhotoId:   types.MustCreateId(),
    PhotoURL:  types.FilePathFromString(gofakeit.URL()),
  }
}

func generateRandomProfileP(userId types.Id) *entity.Profile {
  profile := generateRandomProfile(userId)
  return &profile
}
