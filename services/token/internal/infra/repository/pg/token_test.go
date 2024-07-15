package pg

import (
  "context"
  "fmt"
  sharedConf "github.com/arcorium/nexa/shared/config"
  "github.com/arcorium/nexa/shared/database"
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  "github.com/brianvoe/gofakeit/v7"
  "github.com/stretchr/testify/suite"
  "github.com/testcontainers/testcontainers-go"
  "github.com/testcontainers/testcontainers-go/modules/postgres"
  "github.com/testcontainers/testcontainers-go/wait"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/trace/noop"
  "golang.org/x/crypto/sha3"
  "nexa/services/token/internal/domain/entity"
  "nexa/services/token/internal/infra/repository/model"
  util2 "nexa/services/token/util"
  "reflect"
  "testing"
  "time"
)

const (
  TOKEN_DB_USERNAME = "user"
  TOKEN_DB_PASSWORD = "password"
  TOKEN_DB          = "nexa"

  SEED_TOKEN_DATA_SIZE = 8
)

var tokenSeed []entity.Token

type tokenTestSuite struct {
  suite.Suite
  container *postgres.PostgresContainer
  db        bun.IDB
  tracer    trace.Tracer // Mock
}

func (f *tokenTestSuite) SetupSuite() {
  ctx := context.Background()

  container, err := postgres.RunContainer(ctx,
    postgres.WithUsername(TOKEN_DB_USERNAME),
    postgres.WithPassword(TOKEN_DB_PASSWORD),
    postgres.WithDatabase(TOKEN_DB),
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

  db, err := database.OpenPostgresWithConfig(&sharedConf.PostgresDatabase{
    Address:  fmt.Sprintf("%s:%s", types.Must(container.Host(ctx)), mapped[0].HostPort),
    Username: TOKEN_DB_USERNAME,
    Password: TOKEN_DB_PASSWORD,
    Name:     TOKEN_DB,
    IsSecure: false,
    Timeout:  time.Second * 10,
  }, true)
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
  tokenCount := 0
  token := util.CastSliceP(tokenSeed, func(from *entity.Token) model.Token {
    tokenCount++
    return model.FromTokenDomain(from, func(token *entity.Token, models *model.Token) {
      models.CreatedAt = time.Now().Add(time.Duration(tokenCount) * time.Hour).UTC()
    })
  })

  err = database.Seed(f.db, token...)
  f.Require().NoError(err)
}

func (f *tokenTestSuite) TearDownSuite() {
  err := f.container.Terminate(context.Background())
  f.Require().NoError(err)
}

func (f *tokenTestSuite) Test_tokenRepository_Create() {
  type args struct {
    ctx   context.Context
    token *entity.Token
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Create valid token",
      args: args{
        ctx:   context.Background(),
        token: generateTokenP(),
      },
      wantErr: false,
    },
    {
      name: "Create different token usage for same user",
      args: args{
        ctx: context.Background(),
        token: util.CopyWithP(generateToken(), func(e *entity.Token) {
          e.UserId = tokenSeed[0].UserId
          e.Usage = entity.TokenUsageResetPassword
        }),
      },
      wantErr: false,
    },
    {
      name: "Create duplicated token",
      args: args{
        ctx:   context.Background(),
        token: &tokenSeed[0],
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      t := &tokenRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t1 := f.T()

      err = t.Create(tt.args.ctx, tt.args.token)
      if res := err != nil; res {
        if res != tt.wantErr {
          t1.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      got, err := t.Find(tt.args.ctx, tt.args.token.Token)
      f.Require().NoError(err)

      ignoreTokenField(tt.args.token)
      ignoreTokenField(&got)

      if !reflect.DeepEqual(got, *tt.args.token) != tt.wantErr {
        t1.Errorf("Find() got = %v, want %v", got, *tt.args.token)
      }
    })
  }
}

func (f *tokenTestSuite) Test_tokenRepository_Upsert() {
  type args struct {
    ctx   context.Context
    token *entity.Token
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Create valid new token",
      args: args{
        ctx:   context.Background(),
        token: generateTokenP(),
      },
      wantErr: false,
    },
    {
      name: "Create same token usage and user",
      args: args{
        ctx: context.Background(),
        token: util.CopyWithP(tokenSeed[0], func(e *entity.Token) {
          e.Token = util2.RandomString(0)
          e.ExpiredAt = time.Now().UTC().Add(time.Hour * 2)
        }),
      },
      wantErr: false,
    },
    {
      name: "Create different token usage for same user",
      args: args{
        ctx: context.Background(),
        token: util.CopyWithP(generateToken(), func(e *entity.Token) {
          e.UserId = tokenSeed[0].UserId
          e.Usage = entity.TokenUsageResetPassword
        }),
      },
      wantErr: false,
    },
    {
      // NOTE: Weird result, it is return nil on error and update other fields.
      name: "Create full duplicated token",
      args: args{
        ctx:   context.Background(),
        token: &tokenSeed[0],
      },
      wantErr: false,
    },
    {
      name: "Create duplicated token",
      args: args{
        ctx: context.Background(),
        token: &entity.Token{
          Token:     tokenSeed[0].Token,
          UserId:    types.MustCreateId(),
          Usage:     entity.TokenUsageLogin,
          ExpiredAt: time.Now().Add(time.Hour * 2).UTC(),
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

      t := &tokenRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t1 := f.T()

      err = t.Upsert(tt.args.ctx, tt.args.token)
      if res := err != nil; res {
        if res != tt.wantErr {
          t1.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      got, err := t.Find(tt.args.ctx, tt.args.token.Token)
      f.Require().NoError(err)

      ignoreTokenField(tt.args.token)
      ignoreTokenField(&got)

      if !reflect.DeepEqual(got, *tt.args.token) != tt.wantErr {
        t1.Errorf("Find() got = %v, want %v", got, *tt.args.token)
      }
    })
  }
}

func (f *tokenTestSuite) Test_tokenRepository_Delete() {
  type args struct {
    ctx   context.Context
    token string
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Valid token",
      args: args{
        ctx:   context.Background(),
        token: tokenSeed[0].Token,
      },
      wantErr: false,
    },
    {
      name: "Token not found",
      args: args{
        ctx:   context.Background(),
        token: util.RandomString(32),
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      t := &tokenRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t1 := f.T()

      err = t.Delete(tt.args.ctx, tt.args.token)
      if res := err != nil; res {
        if res != tt.wantErr {
          t1.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      _, err = t.Find(tt.args.ctx, tt.args.token)
      f.Require().Error(err)
    })
  }
}

func (f *tokenTestSuite) Test_tokenRepository_DeleteByUserId() {
  type args struct {
    ctx     context.Context
    userId  types.Id
    baseIdx []int
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Delete single token",
      args: args{
        ctx:     context.Background(),
        userId:  tokenSeed[0].UserId,
        baseIdx: []int{0},
      },
      wantErr: false,
    },
    {
      name: "Delete multiple token",
      args: args{
        ctx:     context.Background(),
        userId:  tokenSeed[1].UserId,
        baseIdx: []int{1, 2},
      },
      wantErr: false,
    },
    {
      name: "Token not found",
      args: args{
        ctx:     context.Background(),
        userId:  types.MustCreateId(),
        baseIdx: nil,
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      t := &tokenRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t1 := f.T()

      err = t.DeleteByUserId(tt.args.ctx, tt.args.userId)
      if res := err != nil; res {
        if res != tt.wantErr {
          t1.Errorf("DeleteByUserId() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      for _, id := range tt.args.baseIdx {
        _, err := t.Find(tt.args.ctx, tokenSeed[id].Token)
        f.Require().Error(err)
      }
    })
  }
}

func (f *tokenTestSuite) Test_tokenRepository_Find() {
  type args struct {
    ctx   context.Context
    token string
  }
  tests := []struct {
    name    string
    args    args
    want    entity.Token
    wantErr bool
  }{
    {
      name: "Valid token",
      args: args{
        ctx:   context.Background(),
        token: tokenSeed[0].Token,
      },
      want:    tokenSeed[0],
      wantErr: false,
    },
    {
      name: "Invalid token",
      args: args{
        ctx:   context.Background(),
        token: util.RandomString(32),
      },
      want:    entity.Token{},
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      t := &tokenRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t1 := f.T()

      got, err := t.Find(tt.args.ctx, tt.args.token)
      if res := err != nil; res {
        if res != tt.wantErr {
          t1.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      ignoreTokenField(&tt.want)
      ignoreTokenField(&got)

      if !reflect.DeepEqual(got, tt.want) != tt.wantErr {
        t1.Errorf("Find() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *tokenTestSuite) Test_tokenRepository_Get() {
  type args struct {
    ctx       context.Context
    parameter repo.QueryParameter
  }
  tests := []struct {
    name    string
    args    args
    want    repo.PaginatedResult[entity.Token]
    wantErr bool
  }{
    // TODO: Add test cases.
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      t := &tokenRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t1 := f.T()

      got, err := t.Get(tt.args.ctx, tt.args.parameter)
      if (err != nil) != tt.wantErr {
        t1.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      ignoreTokensFields(got.Data...)
      ignoreTokensFields(tt.want.Data...)

      comparatorFunc := func(e *entity.Token, e2 *entity.Token) bool {
        return e.Token == e2.Token
      }
      if !util.ArbitraryCheck(got.Data, tt.want.Data, comparatorFunc) != tt.wantErr {
        t1.Errorf("Get() got = %v, want = %v", got, tt.want)
      }
    })
  }
}

func TestToken(t *testing.T) {
  seedTokenData()

  suite.Run(t, &tokenTestSuite{})
}

func seedTokenData() {
  for i := 0; i < SEED_TOKEN_DATA_SIZE; i += 1 {
    tokenSeed = append(tokenSeed, generateToken())
  }
  tokenSeed[0].Usage = entity.TokenUsageEmailVerification
  // Token index 1 and 2 has same user id but different usage
  tokenSeed[1].Usage = entity.TokenUsageResetPassword
  tokenSeed[2].UserId = tokenSeed[1].UserId
  tokenSeed[2].Usage = entity.TokenUsageEmailVerification

}

func generateToken() entity.Token {
  return entity.Token{
    Token:     types.MustCreateId().Hash(sha3.NewShake256()),
    UserId:    types.MustCreateId(),
    Usage:     entity.TokenUsage(gofakeit.UintN(uint(entity.TokenUsageUnknown.Underlying()))),
    ExpiredAt: time.Now().Add(time.Hour*3 + time.Duration(gofakeit.Hour())),
  }
}

func generateTokenP() *entity.Token {
  temp := generateToken()
  return &temp
}

func ignoreTokenField(token *entity.Token) {
  token.ExpiredAt = time.Time{}
}

func ignoreTokensFields(role ...entity.Token) {
  for i := 0; i < len(role); i++ {
    ignoreTokenField(&role[i])
  }
}
