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
  "nexa/services/mailer/constant"
  "nexa/services/mailer/internal/domain/entity"
  "nexa/services/mailer/internal/domain/repository"
  "nexa/services/mailer/internal/infra/repository/model"
  "reflect"
  "slices"
  "testing"
  "time"
)

const (
  MAIL_DB_USERNAME = "user"
  MAIL_DB_PASSWORD = "password"
  MAIL_DB          = "nexa"

  SEED_MAIL_DATA_SIZE = 7
)

type mailTestSuite struct {
  suite.Suite
  container *postgres.PostgresContainer
  db        bun.IDB
  tracer    trace.Tracer // Mock

  tagSeed  []entity.Tag
  mailSeed []entity.Mail
}

func (f *mailTestSuite) SetupSuite() {
  // Create data
  for i := 0; i < SEED_MAIL_DATA_SIZE; i += 1 {
    f.mailSeed = append(f.mailSeed, generateMail())
  }
  f.mailSeed[0].DeliveredAt = time.Time{}
  f.mailSeed[0].Status = entity.StatusPending

  for i := 0; i < SEED_TAG_DATA_SIZE; i += 1 {
    f.tagSeed = append(f.tagSeed, generateTag())
  }

  ctx := context.Background()

  container, err := postgres.RunContainer(ctx,
    postgres.WithUsername(MAIL_DB_USERNAME),
    postgres.WithPassword(MAIL_DB_PASSWORD),
    postgres.WithDatabase(MAIL_DB),
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
    Username: MAIL_DB_USERNAME,
    Password: MAIL_DB_PASSWORD,
    Name:     MAIL_DB,
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
  // Mail
  mails := util.CastSliceP(f.mailSeed, func(from *entity.Mail) model.Mail {
    return model.FromMailDomain(from, func(ent *entity.Mail, mail *model.Mail) {
    })
  })
  // Tag
  tags := util.CastSliceP(f.tagSeed, func(from *entity.Tag) model.Tag {
    return model.FromTagDomain(from, func(ent *entity.Tag, tag *model.Tag) {
      tag.CreatedAt = time.Now()
    })
  })
  // Mail Tag
  mailTags := []model.MailTag{
    {
      MailId: f.mailSeed[3].Id.String(),
      TagId:  f.tagSeed[0].Id.String(),
    },
    {
      MailId: f.mailSeed[4].Id.String(),
      TagId:  f.tagSeed[0].Id.String(),
    },
    {
      MailId: f.mailSeed[5].Id.String(),
      TagId:  f.tagSeed[1].Id.String(),
    },
    {
      MailId: f.mailSeed[6].Id.String(),
      TagId:  f.tagSeed[0].Id.String(),
    },
    {
      MailId: f.mailSeed[6].Id.String(),
      TagId:  f.tagSeed[1].Id.String(),
    },
  }

  err = database.Seed(f.db, mails...)
  f.Require().NoError(err)
  err = database.Seed(f.db, tags...)
  f.Require().NoError(err)
  err = database.Seed(f.db, mailTags...)
  f.Require().NoError(err)
}

func (f *mailTestSuite) TearDownSuite() {
  err := f.container.Terminate(context.Background())
  f.Require().NoError(err)
}

func (f *mailTestSuite) Test_mailRepository_AppendMultipleTags() {
  type args struct {
    ctx      context.Context
    mailTags []repository.MailTags
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Append single mail into multiple tag",
      args: args{
        ctx: context.Background(),
        mailTags: []repository.MailTags{
          {
            First: f.mailSeed[0].Id,
            Second: []types.Id{
              f.tagSeed[0].Id,
              f.tagSeed[1].Id,
            },
          },
        },
      },
      wantErr: false,
    },
    {
      name: "Append multiple mail into multiple tag",
      args: args{
        ctx: context.Background(),
        mailTags: []repository.MailTags{
          {
            First: f.mailSeed[0].Id,
            Second: []types.Id{
              f.tagSeed[1].Id,
            },
          },
          {
            First: f.mailSeed[1].Id,
            Second: []types.Id{
              f.tagSeed[0].Id,
            },
          },
          {
            First: f.mailSeed[2].Id,
            Second: []types.Id{
              f.tagSeed[0].Id,
              f.tagSeed[1].Id,
            },
          },
        },
      },
      wantErr: false,
    },
    {
      name: "Mail not found",
      args: args{
        ctx: context.Background(),
        mailTags: []repository.MailTags{
          {
            First: types.MustCreateId(),
            Second: []types.Id{
              f.tagSeed[0].Id,
            },
          },
        },
      },
      wantErr: true,
    },
    {
      name: "Tag not found",
      args: args{
        ctx: context.Background(),
        mailTags: []repository.MailTags{
          {
            First: f.mailSeed[0].Id,
            Second: []types.Id{
              types.MustCreateId(),
            },
          },
        },
      },
      wantErr: true,
    },
    {
      name: "Some tag is not valid",
      args: args{
        ctx: context.Background(),
        mailTags: []repository.MailTags{
          {
            First: f.mailSeed[0].Id,
            Second: []types.Id{
              types.MustCreateId(),
              f.tagSeed[1].Id,
            },
          },
        },
      },
      wantErr: true,
    },
    {
      name: "Some mail is not valid",
      args: args{
        ctx: context.Background(),
        mailTags: []repository.MailTags{
          {
            First: f.mailSeed[0].Id,
            Second: []types.Id{
              f.tagSeed[0].Id,
            },
          },
          {
            First: types.MustCreateId(),
            Second: []types.Id{
              f.tagSeed[1].Id,
            },
          },
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

      m := mailRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      err = m.AppendMultipleTags(tt.args.ctx, tt.args.mailTags...)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("AppendMultipleTags() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      for _, mailTag := range tt.args.mailTags {
        mails, err := m.FindByIds(tt.args.ctx, mailTag.First)
        f.Require().NoError(err)
        f.Require().Len(mails, 1)
        f.Require().GreaterOrEqual(len(mails[0].Tags), len(mailTag.Second))
        for _, tagId := range mailTag.Second {
          res := slices.ContainsFunc(mails[0].Tags, func(tag entity.Tag) bool {
            return tag.Id == tagId
          })
          f.Require().True(res)
        }
      }
    })
  }
}

func (f *mailTestSuite) Test_mailRepository_AppendTags() {
  type args struct {
    ctx    context.Context
    mailId types.Id
    tagIds []types.Id
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Append single tag",
      args: args{
        ctx:    context.Background(),
        mailId: f.mailSeed[0].Id,
        tagIds: []types.Id{f.tagSeed[0].Id},
      },
      wantErr: false,
    },
    {
      name: "Append multiple tag",
      args: args{
        ctx:    context.Background(),
        mailId: f.mailSeed[0].Id,
        tagIds: []types.Id{f.tagSeed[0].Id, f.tagSeed[1].Id},
      },
      wantErr: false,
    },
    {
      name: "Mail not found",
      args: args{
        ctx:    context.Background(),
        mailId: types.MustCreateId(),
        tagIds: []types.Id{f.tagSeed[0].Id},
      },
      wantErr: true,
    },
    {
      name: "Tag not found",
      args: args{
        ctx:    context.Background(),
        mailId: f.mailSeed[0].Id,
        tagIds: []types.Id{types.MustCreateId()},
      },
      wantErr: true,
    },
    {
      name: "Some tag is not valid",
      args: args{
        ctx:    context.Background(),
        mailId: f.mailSeed[0].Id,
        tagIds: []types.Id{f.tagSeed[0].Id, types.MustCreateId()},
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      m := mailRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      err = m.AppendTags(tt.args.ctx, tt.args.mailId, tt.args.tagIds)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("AppendTags() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      mail, err := m.FindByIds(tt.args.ctx, tt.args.mailId)
      f.Require().NoError(err)
      f.Require().Len(mail, 1)

      for _, tagId := range tt.args.tagIds {
        if !slices.ContainsFunc(mail[0].Tags, func(tag entity.Tag) bool {
          return tag.Id.Eq(tagId)
        }) {
          t.Errorf("AppendTags() error = tags not found on mail fields")
        }
      }
    })
  }
}

func (f *mailTestSuite) Test_mailRepository_Create() {
  type args struct {
    ctx   context.Context
    mails []entity.Mail
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Create single mail",
      args: args{
        ctx:   context.Background(),
        mails: util.GenerateMultiple(1, generateMail),
      },
      wantErr: false,
    },
    {
      name: "Create multiple mail",
      args: args{
        ctx:   context.Background(),
        mails: util.GenerateMultiple(3, generateMail),
      },
      wantErr: false,
    },
    {
      name: "Dupplicate id",
      args: args{
        ctx: context.Background(),
        mails: []entity.Mail{
          util.CopyWith(generateMail(), func(e *entity.Mail) {
            e.Id = f.mailSeed[0].Id
          }),
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

      m := mailRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      err = m.Create(tt.args.ctx, tt.args.mails...)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      for _, mail := range tt.args.mails {
        gotMail, err := m.FindByIds(tt.args.ctx, mail.Id)
        f.Require().NoError(err)
        f.Require().Len(gotMail, 1)

        ignoreMailsFields(gotMail...)
        ignoreMailFields(&mail)

        if !reflect.DeepEqual(gotMail[0], mail) != tt.wantErr {
          t.Errorf("Get() got = %v, want %v", gotMail[0], mail)
        }
      }
    })
  }
}

func (f *mailTestSuite) Test_mailRepository_Get() {
  type args struct {
    ctx   context.Context
    query repo.QueryParameter
  }
  tests := []struct {
    name    string
    args    args
    want    repo.PaginatedResult[entity.Mail]
    wantErr bool
  }{
    {
      name: "Get all",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: 0,
          Limit:  0,
        },
      },
      want: repo.PaginatedResult[entity.Mail]{
        Data:    f.mailSeed,
        Total:   uint64(len(f.mailSeed)),
        Element: uint64(len(f.mailSeed)),
      },
      wantErr: false,
    },
    {
      name: "Get all with offset and limit",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: 1,
          Limit:  2,
        },
      },
      want: repo.PaginatedResult[entity.Mail]{
        Data:    f.mailSeed[1:3],
        Total:   uint64(len(f.mailSeed)),
        Element: 2,
      },
      wantErr: false,
    },
    {
      name: "Get all with offset",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: 2,
          Limit:  0,
        },
      },
      want: repo.PaginatedResult[entity.Mail]{
        Data:    f.mailSeed[2:],
        Total:   uint64(len(f.mailSeed)),
        Element: uint64(len(f.mailSeed)) - 2,
      },
      wantErr: false,
    },
    {
      name: "Get all with limit",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: 0,
          Limit:  3,
        },
      },
      want: repo.PaginatedResult[entity.Mail]{
        Data:    f.mailSeed[:3],
        Total:   uint64(len(f.mailSeed)),
        Element: 3,
      },
      wantErr: false,
    },
    {
      name: "Get out of bound offset",
      args: args{
        ctx: context.Background(),
        query: repo.QueryParameter{
          Offset: uint64(len(f.mailSeed)),
          Limit:  3,
        },
      },
      want: repo.PaginatedResult[entity.Mail]{
        Data:    nil,
        Total:   uint64(len(f.mailSeed)),
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

      m := mailRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      got, err := m.Get(tt.args.ctx, tt.args.query)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      ignoreMailsFields(got.Data...)
      ignoreMailsFields(tt.want.Data...)

      if !reflect.DeepEqual(got, tt.want) != tt.wantErr {
        t.Errorf("Get() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *mailTestSuite) Test_mailRepository_FindByIds() {
  type args struct {
    ctx context.Context
    ids []types.Id
  }
  tests := []struct {
    name    string
    args    args
    want    []entity.Mail
    wantErr bool
  }{
    {
      name: "Get single mail",
      args: args{
        ctx: context.Background(),
        ids: []types.Id{f.mailSeed[0].Id},
      },
      want:    f.mailSeed[:1],
      wantErr: false,
    },
    {
      name: "Get multiple mails",
      args: args{
        ctx: context.Background(),
        ids: []types.Id{f.mailSeed[0].Id, f.mailSeed[1].Id},
      },
      want:    f.mailSeed[:2],
      wantErr: false,
    },
    {
      name: "Some mail is not valid",
      args: args{
        ctx: context.Background(),
        ids: []types.Id{f.mailSeed[2].Id, types.MustCreateId(), f.mailSeed[1].Id},
      },
      want:    []entity.Mail{f.mailSeed[1], f.mailSeed[2]},
      wantErr: false,
    },
    {
      name: "Mail not found",
      args: args{
        ctx: context.Background(),
        ids: []types.Id{types.MustCreateId()},
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "All mail not found",
      args: args{
        ctx: context.Background(),
        ids: []types.Id{types.MustCreateId(), types.MustCreateId(), types.MustCreateId()},
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

      m := mailRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      got, err := m.FindByIds(tt.args.ctx, tt.args.ids...)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("FindByIds() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      ignoreMailsFields(got...)
      ignoreMailsFields(tt.want...)

      if !reflect.DeepEqual(got, tt.want) != tt.wantErr {
        t.Errorf("FindByIds() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *mailTestSuite) Test_mailRepository_FindByTag() {
  type args struct {
    ctx context.Context
    tag types.Id
  }
  tests := []struct {
    name    string
    args    args
    want    []entity.Mail
    wantErr bool
  }{
    {
      name: "Valid tag with multiple mails",
      args: args{
        ctx: context.Background(),
        tag: f.tagSeed[0].Id,
      },
      want:    []entity.Mail{f.mailSeed[3], f.mailSeed[4], f.mailSeed[6]},
      wantErr: false,
    },
    {
      name: "Valid tag without mails",
      args: args{
        ctx: context.Background(),
        tag: f.tagSeed[2].Id,
      },
      want:    nil,
      wantErr: true,
    },
    {
      name: "Invalid tag",
      args: args{
        ctx: context.Background(),
        tag: types.MustCreateId(),
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

      m := mailRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      got, err := m.FindByTag(tt.args.ctx, tt.args.tag)
      if (err != nil) != tt.wantErr {
        t.Errorf("FindByTag() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      ignoreMailsFields(got...)
      ignoreMailsFields(tt.want...)

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindByTag() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func (f *mailTestSuite) Test_mailRepository_Patch() {
  type args struct {
    ctx  context.Context
    mail *entity.Mail
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Set mail is successfully delivered",
      args: args{
        ctx: context.Background(),
        mail: &entity.Mail{
          Id:          f.mailSeed[0].Id,
          Status:      entity.StatusDelivered,
          DeliveredAt: time.Now().UTC(),
        },
      },
      wantErr: false,
    },
    {
      name: "Set mail is failed to deliver",
      args: args{
        ctx: context.Background(),
        mail: &entity.Mail{
          Id:     f.mailSeed[0].Id,
          Status: entity.StatusFailed,
        },
      },
      wantErr: false,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      m := mailRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      err = m.Patch(tt.args.ctx, tt.args.mail)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("Patch() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      // Compare
      mail, err := m.FindByIds(tt.args.ctx, tt.args.mail.Id)
      f.Require().NoError(err)
      f.Require().Len(mail, 1)

      mail[0].DeliveredAt = mail[0].DeliveredAt.Round(time.Minute)
      tt.args.mail.DeliveredAt = tt.args.mail.DeliveredAt.Round(time.Minute)

      if mail[0].Status != tt.args.mail.Status ||
          mail[0].DeliveredAt != tt.args.mail.DeliveredAt {
        t.Errorf("Patch() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func (f *mailTestSuite) Test_mailRepository_Remove() {
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
      name: "Remove valid mail",
      args: args{
        ctx: context.Background(),
        id:  f.mailSeed[0].Id,
      },
      wantErr: false,
    },
    {
      name: "Remove invalid mail",
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

      m := mailRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      err = m.Remove(tt.args.ctx, tt.args.id)
      if res := err != nil; res {
        if res != tt.wantErr {

          t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      mail, err := m.FindByIds(tt.args.ctx, tt.args.id)
      f.Require().Error(err)
      f.Require().Nil(mail)
    })
  }
}

func (f *mailTestSuite) Test_mailRepository_RemoveTags() {
  type args struct {
    ctx    context.Context
    mailId types.Id
    tagIds []types.Id
  }
  tests := []struct {
    name    string
    args    args
    wantErr bool
  }{
    {
      name: "Remove valid mail and single valid tags",
      args: args{
        ctx:    context.Background(),
        mailId: f.mailSeed[3].Id,
        tagIds: []types.Id{f.tagSeed[0].Id},
      },
      wantErr: false,
    },
    {
      name: "Remove valid mail and multiple valid tags",
      args: args{
        ctx:    context.Background(),
        mailId: f.mailSeed[6].Id,
        tagIds: []types.Id{f.tagSeed[0].Id, f.tagSeed[1].Id},
      },
      wantErr: false,
    },
    {
      name: "Remove invalid mail",
      args: args{
        ctx:    context.Background(),
        mailId: types.MustCreateId(),
        tagIds: []types.Id{f.tagSeed[0].Id},
      },
      wantErr: true,
    },
    {
      name: "Remove valid mail and invalid tags",
      args: args{
        ctx:    context.Background(),
        mailId: f.mailSeed[0].Id,
        tagIds: []types.Id{f.tagSeed[2].Id, types.MustCreateId()},
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    f.Run(tt.name, func() {
      tx, err := f.db.BeginTx(tt.args.ctx, nil)
      f.Require().NoError(err)
      defer tx.Rollback()

      m := mailRepository{
        db:     tx,
        tracer: f.tracer,
      }
      t := f.T()

      err = m.RemoveTags(tt.args.ctx, tt.args.mailId, tt.args.tagIds)
      if res := err != nil; res {
        if res != tt.wantErr {
          t.Errorf("RemoveTags() error = %v, wantErr %v", err, tt.wantErr)
        }
        return
      }

      mail, err := m.FindByIds(tt.args.ctx, tt.args.mailId)
      f.Require().NoError(err)
      f.Require().Len(mail, 1)

      for _, tagId := range tt.args.tagIds {
        if slices.ContainsFunc(mail[0].Tags, func(tag entity.Tag) bool {
          return tag.Id.Eq(tagId)
        }) {
          t.Errorf("AppendTags() error = tags not found on mail fields")
        }
      }
    })
  }
}

func TestMail(t *testing.T) {

  suite.Run(t, &mailTestSuite{})
}

var count = 0

func generateMail() entity.Mail {
  count += 1
  return entity.Mail{
    Id:          types.MustCreateId(),
    Subject:     gofakeit.Animal(),
    Recipient:   types.Email(gofakeit.Email()),
    Sender:      constant.SERVICE_MAIL_SENDER,
    BodyType:    entity.MailBodyType(gofakeit.UintN(uint(entity.BodyTypeUnknown.Underlying()))),
    Body:        gofakeit.LoremIpsumSentence(50),
    Status:      entity.Status(gofakeit.UintN(uint(entity.StatusFailed.Underlying()))),
    SentAt:      time.Now().Add(time.Hour * time.Duration(count*-1)).UTC(),
    DeliveredAt: time.Now().Add(time.Hour*time.Duration(count) + time.Hour*time.Duration(gofakeit.Hour())).UTC(),
  }
}

func generateMailP() *entity.Mail {
  user := generateMail()
  return &user
}

func ignoreMailFields(mail *entity.Mail) {
  mail.Body = ""
  mail.BodyType = 0
  mail.SentAt = mail.SentAt.Round(time.Hour)
  mail.DeliveredAt = mail.DeliveredAt.Round(time.Hour)
  mail.Tags = nil
}

func ignoreMailsFields(mail ...entity.Mail) {
  for i := 0; i < len(mail); i++ {
    ignoreMailFields(&mail[i])
  }
}
