package redis

import (
  "context"
  "fmt"
  "github.com/brianvoe/gofakeit/v7"
  "github.com/redis/go-redis/v9"
  "github.com/stretchr/testify/suite"
  "github.com/testcontainers/testcontainers-go"
  rtest "github.com/testcontainers/testcontainers-go/modules/redis"
  "github.com/testcontainers/testcontainers-go/wait"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/trace/noop"
  domain "nexa/services/authentication/internal/domain/entity"
  sharedJwt "nexa/shared/jwt"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  "nexa/shared/util/repo"
  "reflect"
  "strconv"
  "testing"
  "time"
)

const (
  USERNAME = "user"
  PASSWORD = "password"
  DATABASE = "nexa"
)

var dataSeed []domain.Credential

const METADATA_SEED_SIZE = 10

type credentialTestSuite struct {
  suite.Suite
  container *rtest.RedisContainer
  client    redis.UniversalClient
  tracer    trace.Tracer
}

func (c *credentialTestSuite) SetupSuite() {
  ctx := context.Background()

  container, err := rtest.RunContainer(ctx,
    testcontainers.WithWaitStrategy(wait.ForExposedPort()),
  )
  c.Require().NoError(err)
  c.container = container

  inspect, err := container.Inspect(ctx)
  c.Require().NoError(err)
  ports := inspect.NetworkSettings.Ports
  mapped := ports["6379/tcp"]

  sharedUtil.DoNothing(ports, mapped)

  c.client = redis.NewClient(&redis.Options{
    Addr: fmt.Sprintf("localhost:%d", uint16(types.Must(strconv.Atoi(mapped[0].HostPort)))),
    //Username: "",
    //Password: "",
    //DB:       0,
  })

  stat := c.client.Ping(ctx)
  c.Require().NoError(stat.Err())

  // Tracer
  provider := noop.NewTracerProvider()
  c.tracer = provider.Tracer("MOCK")
}

func (c *credentialTestSuite) TearDownSuite() {
  err := c.container.Terminate(context.Background())
  c.Require().NoError(err)

}

func (c *credentialTestSuite) SetupTest() {
  ctx := context.Background()
  // Seed
  credRepo := NewCredential(c.client, nil)
  for i := 0; i < METADATA_SEED_SIZE; i++ {
    err := credRepo.Create(ctx, &dataSeed[i])
    c.Require().NoError(err)
  }
  fmt.Println("============== SETUP TEST ==============")
}

func (c *credentialTestSuite) TearDownTest() {
  // Remove data
  statCmd := c.client.FlushDB(context.Background())
  c.Require().NoError(statCmd.Err())
  fmt.Println("============== END TEST ==============")
}

func (c *credentialTestSuite) Test_credentialRepository_Create() {
  ctx := context.Background()

  type args struct {
    ctx        context.Context
    credential *domain.Credential
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        ctx:        ctx,
        credential: generateRandomCredentialP(),
      },
      wantErr: false,
    },
    {
      name: "Duplicate",
      args: args{
        ctx:        ctx,
        credential: &dataSeed[METADATA_SEED_SIZE-1],
      },
      wantErr: true,
    },
    {
      name: "Duplicate with different id",
      args: args{
        ctx: ctx,
        credential: sharedUtil.CopyWithP(dataSeed[METADATA_SEED_SIZE-1], func(i *domain.Credential) {
          i.Id = types.MustCreateId()
        }),
      },
      wantErr: false,
    },
  }
  for _, tt := range tests {
    c.Run(tt.name, func() {
      repos := &credentialRepository{
        config: defaultCredential,
        client: c.client,
        tracer: c.tracer,
      }
      if err := repos.Create(tt.args.ctx, tt.args.credential); (err != nil) != tt.wantErr {
        c.T().Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func (c *credentialTestSuite) Test_credentialRepository_Delete() {
  ctx := context.Background()

  type args struct {
    ctx             context.Context
    refreshTokenIds []types.Id
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Single UserId",
      args: args{
        ctx:             ctx,
        refreshTokenIds: []types.Id{dataSeed[0].Id},
      },
      wantErr: false,
    },
    {
      name: "Multiple Ids",
      args: args{
        ctx:             ctx,
        refreshTokenIds: []types.Id{dataSeed[1].Id, dataSeed[2].Id},
      },
      wantErr: false,
    },
    {
      name: "Combination Of Ids",
      args: args{
        ctx:             ctx,
        refreshTokenIds: []types.Id{types.MustCreateId(), dataSeed[3].Id, types.MustCreateId()},
      },
      wantErr: false,
    },
    {
      name: "UserId Not Found",
      args: args{
        ctx:             ctx,
        refreshTokenIds: []types.Id{types.MustCreateId()},
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    c.Run(tt.name, func() {
      repos := &credentialRepository{
        config: defaultCredential,
        client: c.client,
        tracer: c.tracer,
      }
      if err := repos.Delete(tt.args.ctx, tt.args.refreshTokenIds...); (err != nil) != tt.wantErr {
        c.T().Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func (c *credentialTestSuite) Test_credentialRepository_DeleteByUserId() {
  ctx := context.Background()

  type args struct {
    ctx    context.Context
    userId types.Id
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        ctx:    ctx,
        userId: dataSeed[0].UserId,
      },
      wantErr: false,
    },
    {
      name: "User id does not exist",
      args: args{
        ctx:    ctx,
        userId: types.MustCreateId(),
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    c.Run(tt.name, func() {
      repos := &credentialRepository{
        config: defaultCredential,
        client: c.client,
        tracer: c.tracer,
      }
      if err := repos.DeleteByUserId(tt.args.ctx, tt.args.userId); (err != nil) != tt.wantErr {
        c.T().Errorf("DeleteByUserId() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func (c *credentialTestSuite) Test_credentialRepository_Find() {
  ctx := context.Background()

  type args struct {
    ctx            context.Context
    refreshTokenId types.Id
  }
  tests := []struct {
    name    string
    args    args
    want    *domain.Credential
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        ctx:            ctx,
        refreshTokenId: dataSeed[0].Id,
      },
      want:    &dataSeed[0],
      wantErr: false,
    },
    {
      name: "Credential Not Found",
      args: args{
        ctx:            ctx,
        refreshTokenId: types.MustCreateId(),
      },
      want:    nil,
      wantErr: true,
    },
  }
  for _, tt := range tests {
    c.Run(tt.name, func() {
      repos := &credentialRepository{
        config: defaultCredential,
        client: c.client,
        tracer: c.tracer,
      }
      got, err := repos.Find(tt.args.ctx, tt.args.refreshTokenId)
      if (err != nil) != tt.wantErr {
        c.T().Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if !reflect.DeepEqual(got, tt.want) {
        c.T().Errorf("Find() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (c *credentialTestSuite) Test_credentialRepository_FindAll() {
  ctx := context.Background()

  type args struct {
    ctx       context.Context
    parameter repo.QueryParameter
  }
  tests := []struct {
    name    string
    args    args
    want    repo.PaginatedResult[domain.Credential]
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        ctx: ctx,
        parameter: repo.QueryParameter{
          Offset: 0,
          Limit:  3,
        },
      },
      want: repo.PaginatedResult[domain.Credential]{
        Data:    nil,
        Total:   METADATA_SEED_SIZE,
        Element: 3,
      },
      wantErr: false,
    },
    {
      name: "Zero Limit",
      args: args{
        ctx: ctx,
        parameter: repo.QueryParameter{
          Offset: 0,
          Limit:  0,
        },
      },
      want: repo.PaginatedResult[domain.Credential]{
        Data:    nil,
        Total:   METADATA_SEED_SIZE,
        Element: METADATA_SEED_SIZE,
      },
      wantErr: false,
    },
  }
  for _, tt := range tests {
    c.Run(tt.name, func() {
      repos := &credentialRepository{
        config: defaultCredential,
        client: c.client,
        tracer: c.tracer,
      }
      got, err := repos.FindAll(tt.args.ctx, tt.args.parameter)
      if isErr := err != nil; isErr {
        if isErr == tt.wantErr {
          return
        }
        c.T().Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      c.Require().Equal(tt.want.Element, got.Element)
      c.Require().Equal(tt.want.Total, got.Total)

      if got.Data == nil {
        if tt.want.Data != nil {
          c.T().Errorf("FindAll() got = %v, want %v", got, tt.want)
        }
      }
    })
  }
}

func (c *credentialTestSuite) Test_credentialRepository_FindByUserId() {
  ctx := context.Background()

  type args struct {
    ctx    context.Context
    userId types.Id
  }
  tests := []struct {
    name    string
    args    args
    want    []domain.Credential
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        ctx:    ctx,
        userId: dataSeed[0].UserId,
      },
      want:    dataSeed[:3],
      wantErr: false,
    },
    {
      name: "User Not Found",
      args: args{
        ctx:    ctx,
        userId: types.MustCreateId(),
      },
      want:    nil,
      wantErr: true,
    },
  }
  for _, tt := range tests {
    c.Run(tt.name, func() {
      repos := &credentialRepository{
        config: defaultCredential,
        client: c.client,
        tracer: c.tracer,
      }
      got, err := repos.FindByUserId(tt.args.ctx, tt.args.userId)
      if isErr := err != nil; isErr {
        if isErr == tt.wantErr {
          return
        }
        c.T().Errorf("FindByUserId() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      c.Require().Equal(len(got), len(tt.want))

      for _, g := range got {
        found := false
        for _, w := range tt.want {
          if g.Id != w.Id {
            continue
          }

          found = true
          if !reflect.DeepEqual(got, tt.want) {
            c.T().Errorf("FindByUserId() got = %v, want %v", got, tt.want)
          }
          break
        }

        if !found {
          panic("Error")
        }
      }

    })
  }
}

func (c *credentialTestSuite) Test_credentialRepository_Patch() {
  ctx := context.Background()
  type args struct {
    basedId    int
    ctx        context.Context
    credential *domain.Credential
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Normal",
      args: args{
        ctx: ctx,
        credential: &domain.Credential{
          Id: dataSeed[0].Id,
          Device: domain.Device{
            Name: "Changed",
          },
          RefreshToken: sharedUtil.RandomString(32),
          ExpiresAt:    time.Now().Add(time.Hour * 1),
        },
        basedId: 0,
      },
      wantErr: false,
    },
    {
      name: "Change Access AccessToken",
      args: args{
        ctx: ctx,
        credential: &domain.Credential{
          Id:            dataSeed[1].Id,
          AccessTokenId: types.MustCreateId(),
          Device: domain.Device{
            Name: "Changed",
          },
        },
        basedId: 1,
      },
      wantErr: false,
    },
    {
      name: "Credential Not found",
      args: args{
        ctx: ctx,
        credential: &domain.Credential{
          Id:            types.MustCreateId(),
          AccessTokenId: types.MustCreateId(),
          Device: domain.Device{
            Name: "Changed",
          },
        },
        basedId: -1,
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    c.Run(tt.name, func() {
      repos := &credentialRepository{
        config: defaultCredential,
        client: c.client,
        tracer: c.tracer,
      }

      err := repos.Patch(tt.args.ctx, tt.args.credential)

      if (err != nil) != tt.wantErr {
        c.T().Errorf("Patch() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func TestCredential(t *testing.T) {
  for i := 0; i < METADATA_SEED_SIZE; i += 1 {
    dataSeed = append(dataSeed, generateRandomCredential())
  }
  // User id index 0 will have 3 credentials
  dataSeed[1].UserId = dataSeed[0].UserId
  dataSeed[2].UserId = dataSeed[0].UserId

  suite.Run(t, &credentialTestSuite{})
}

func generateRandomCredential() domain.Credential {
  times := time.Now().Add(time.Hour * time.Duration(gofakeit.Hour()+2))
  times = time.Unix(times.Unix(), 0) // Round to second
  return domain.Credential{
    Id:            types.MustCreateId(),
    UserId:        types.MustCreateId(),
    AccessTokenId: types.MustCreateId(),
    Device:        domain.Device{Name: gofakeit.Name()},
    RefreshToken:  sharedJwt.GenerateRefreshToken(),
    ExpiresAt:     times,
  }
}

func generateRandomCredentialP() *domain.Credential {
  cred := generateRandomCredential()
  return &cred
}
