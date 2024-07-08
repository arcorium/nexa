package service

import (
  "context"
  "database/sql"
  "github.com/brianvoe/gofakeit/v7"
  "github.com/stretchr/testify/mock"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/trace/noop"
  uow2 "nexa/services/mailer/internal/app/uow"
  "nexa/services/mailer/internal/domain/dto"
  "nexa/services/mailer/internal/domain/entity"
  extMock "nexa/services/mailer/internal/domain/external/mocks"
  repoMock "nexa/services/mailer/internal/domain/repository/mocks"
  sharedDto "nexa/shared/dto"
  "nexa/shared/status"
  "nexa/shared/types"
  "nexa/shared/uow"
  uowMock "nexa/shared/uow/mocks"
  sharedUtil "nexa/shared/util"
  "nexa/shared/util/repo"
  "reflect"
  "testing"
)

func newMailMocked(t *testing.T) mailMocked {
  // Tracer
  provider := noop.NewTracerProvider()
  return mailMocked{
    UOW:        uowMock.NewUnitOfWorkMock[uow2.MailStorage](t),
    Mail:       repoMock.NewMailMock(t),
    Tag:        repoMock.NewTagMock(t),
    MailClient: extMock.NewMailMock(t),
    Tracer:     provider.Tracer("MOCK"),
  }
}

type mailMocked struct {
  UOW        *uowMock.UnitOfWorkMock[uow2.MailStorage]
  Mail       *repoMock.MailMock
  Tag        *repoMock.TagMock
  MailClient *extMock.MailMock
  Tracer     trace.Tracer
}

func (m *mailMocked) defaultUOWMock() {
  m.UOW.EXPECT().
    Repositories().
    Return(uow2.NewStorage(m.Mail, m.Tag))
}

func (m *mailMocked) txProxy() {
  m.UOW.On("DoTx", mock.Anything, mock.Anything).
      Return(func(ctx context.Context, f uow.UOWBlock[uow2.MailStorage]) error {
        return f(ctx, uow2.NewStorage(m.Mail, m.Tag))
      })
}

type setupMailTestFunc func(mocked *mailMocked, arg any, want any)

func Test_mailService_GetAll(t *testing.T) {
  type args struct {
    ctx      context.Context
    pagedDTO *sharedDto.PagedElementDTO
  }
  tests := []struct {
    name  string
    setup setupMailTestFunc
    args  args
    want  sharedDto.PagedElementResult[dto.MailResponseDTO]
    want1 status.Object
  }{
    {
      name: "Success get all mails",
      setup: func(mocked *mailMocked, arg any, want any) {
        a := arg.(*args)
        w := want.(*sharedDto.PagedElementResult[dto.MailResponseDTO])

        mocked.defaultUOWMock()

        mocked.Mail.EXPECT().
          Get(mock.Anything, a.pagedDTO.ToQueryParam()).
          Return(repo.PaginatedResult[entity.Mail]{nil, w.TotalElements, w.Element}, nil)
      },
      args: args{
        ctx: context.Background(),
        pagedDTO: &sharedDto.PagedElementDTO{
          Element: 5,
          Page:    2,
        },
      },
      want: sharedDto.PagedElementResult[dto.MailResponseDTO]{
        Data:          nil,
        Element:       0,
        Page:          2,
        TotalElements: 50,
        TotalPages:    10,
      },
      want1: status.Success(),
    },
    {
      name: "Failed to get all tags",
      setup: func(mocked *mailMocked, arg any, want any) {
        a := arg.(*args)
        w := want.(*sharedDto.PagedElementResult[dto.MailResponseDTO])

        mocked.defaultUOWMock()

        mocked.Mail.EXPECT().
          Get(mock.Anything, a.pagedDTO.ToQueryParam()).
          Return(repo.PaginatedResult[entity.Mail]{nil, w.TotalElements, w.Element}, sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        pagedDTO: &sharedDto.PagedElementDTO{
          Element: 5,
          Page:    2,
        },
      },
      want: sharedDto.PagedElementResult[dto.MailResponseDTO]{
        Data:          nil,
        Element:       0,
        Page:          0,
        TotalElements: 0,
        TotalPages:    0,
      },
      want1: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newMailMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, &tt.want)
      }

      m := &mailService{
        mailExt: mocked.MailClient,
        mailUow: mocked.UOW,
        tracer:  mocked.Tracer,
      }
      got, got1 := m.GetAll(tt.args.ctx, tt.args.pagedDTO)
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("GetAll() got = %v, want %v", got, tt.want)
      }
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("GetAll() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_mailService_FindByIds(t *testing.T) {
  type args struct {
    ctx     context.Context
    mailIds []types.Id
  }
  tests := []struct {
    name  string
    setup setupMailTestFunc
    args  args
    want  []dto.MailResponseDTO
    want1 status.Object
  }{
    {
      name: "Success find mails by ids",
      setup: func(mocked *mailMocked, arg any, want any) {
        a := arg.(*args)
        w := want.([]dto.MailResponseDTO)

        mocked.defaultUOWMock()

        mailIds := sharedUtil.CastSlice(a.mailIds, sharedUtil.ToAny[types.Id])

        mails := sharedUtil.CastSliceP(w, func(from *dto.MailResponseDTO) entity.Mail {
          return entity.Mail{
            Id:        from.Id,
            Subject:   from.Subject,
            Recipient: from.Recipient,
            Sender:    from.Sender,
            Status:    from.Status,
          }
        })

        mocked.Mail.EXPECT().
          FindByIds(mock.Anything, mailIds...).
          Return(mails, nil)
      },
      args: args{
        ctx:     context.Background(),
        mailIds: sharedUtil.GenerateMultiple(2, types.MustCreateId),
      },
      want: sharedUtil.GenerateMultiple(2, func() dto.MailResponseDTO {
        return dto.MailResponseDTO{
          Id:        types.MustCreateId(),
          Subject:   gofakeit.AnimalType(),
          Recipient: types.Email(gofakeit.Email()),
          Sender:    types.Email(gofakeit.Email()),
          Status:    entity.Status(gofakeit.UintN(uint(entity.StatusFailed.Underlying()))),
        }
      }),
      want1: status.Success(),
    },
    {
      name: "Failed to find mails by ids",
      setup: func(mocked *mailMocked, arg any, want any) {
        a := arg.(*args)
        mocked.defaultUOWMock()

        mailIds := sharedUtil.CastSlice(a.mailIds, sharedUtil.ToAny[types.Id])

        mocked.Mail.EXPECT().
          FindByIds(mock.Anything, mailIds...).
          Return(nil, sql.ErrNoRows)
      },
      args: args{
        ctx:     context.Background(),
        mailIds: sharedUtil.GenerateMultiple(2, types.MustCreateId),
      },
      want:  nil,
      want1: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newMailMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, tt.want)
      }

      m := &mailService{
        mailExt: mocked.MailClient,
        mailUow: mocked.UOW,
        tracer:  mocked.Tracer,
      }
      got, got1 := m.FindByIds(tt.args.ctx, tt.args.mailIds...)
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindByIds() got = %v, want %v", got, tt.want)
      }
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("FindByIds() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_mailService_FindByTag(t *testing.T) {
  type args struct {
    ctx   context.Context
    tagId types.Id
  }
  tests := []struct {
    name  string
    setup setupMailTestFunc
    args  args
    want  []dto.MailResponseDTO
    want1 status.Object
  }{
    {
      name: "Success find mails by tag",
      setup: func(mocked *mailMocked, arg any, want any) {
        a := arg.(*args)
        w := want.([]dto.MailResponseDTO)

        mocked.defaultUOWMock()

        mails := sharedUtil.CastSliceP(w, func(from *dto.MailResponseDTO) entity.Mail {
          return entity.Mail{
            Id:        from.Id,
            Subject:   from.Subject,
            Recipient: from.Recipient,
            Sender:    from.Sender,
            Status:    from.Status,
          }
        })

        mocked.Mail.EXPECT().
          FindByTag(mock.Anything, a.tagId).
          Return(mails, nil)
      },
      args: args{
        ctx:   context.Background(),
        tagId: dummyId,
      },
      want: sharedUtil.GenerateMultiple(2, func() dto.MailResponseDTO {
        return dto.MailResponseDTO{
          Id:        types.MustCreateId(),
          Subject:   gofakeit.AnimalType(),
          Recipient: types.Email(gofakeit.Email()),
          Sender:    types.Email(gofakeit.Email()),
          Status:    entity.Status(gofakeit.UintN(uint(entity.StatusFailed.Underlying()))),
        }
      }),
      want1: status.Success(),
    },
    {
      name: "Failed to find mails by ids",
      setup: func(mocked *mailMocked, arg any, want any) {
        a := arg.(*args)
        mocked.defaultUOWMock()

        mocked.Mail.EXPECT().
          FindByTag(mock.Anything, a.tagId).
          Return(nil, sql.ErrNoRows)
      },
      args: args{
        ctx:   context.Background(),
        tagId: dummyId,
      },
      want:  nil,
      want1: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newMailMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, tt.want)
      }

      m := &mailService{
        mailExt: mocked.MailClient,
        mailUow: mocked.UOW,
        tracer:  mocked.Tracer,
      }
      got, got1 := m.FindByTag(tt.args.ctx, tt.args.tagId)
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindByTag() got = %v, want %v", got, tt.want)
      }
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("FindByTag() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_mailService_Remove(t *testing.T) {
  type args struct {
    ctx    context.Context
    mailId types.Id
  }
  tests := []struct {
    name  string
    setup setupMailTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success remove mail by id",
      setup: func(mocked *mailMocked, arg any, want any) {
        a := arg.(*args)
        mocked.defaultUOWMock()

        mocked.Mail.EXPECT().
          Remove(mock.Anything, a.mailId).
          Return(nil)
      },
      args: args{
        ctx:    context.Background(),
        mailId: types.MustCreateId(),
      },
      want: status.Deleted(),
    },
    {
      name: "Failed to remove mail by id",
      setup: func(mocked *mailMocked, arg any, want any) {
        a := arg.(*args)
        mocked.defaultUOWMock()

        mocked.Mail.EXPECT().
          Remove(mock.Anything, a.mailId).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx:    context.Background(),
        mailId: types.MustCreateId(),
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newMailMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      m := &mailService{
        mailExt: mocked.MailClient,
        mailUow: mocked.UOW,
        tracer:  mocked.Tracer,
      }
      if got := m.Remove(tt.args.ctx, tt.args.mailId); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("Remove() = %v, want %v", got, tt.want)
      }
    })
  }
}

/*
func Test_mailService_Send(t *testing.T) {
  t.Parallel()
  type args struct {
    ctx     context.Context
    mailDTO *dto.SendMailDTO
  }
  tests := []struct {
    name  string
    setup setupMailTestFunc
    args  args
    want  int // Len of returned ids
    want1 status.Object
  }{
    {
      name: "Success send single mail with tags",
      setup: func(mocked *mailMocked, arg any, want any) {
        //a := arg.(*args)
        // Create proxy
        mocked.defaultUOWMock()
        mocked.txProxy()

        mocked.Mail.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(nil)

        mocked.Mail.EXPECT().
          AppendMultipleTags(mock.Anything, mock.Anything).
          Return(nil)

        mocked.MailClient.EXPECT().
          Send(mock.Anything, mock.Anything, mock.Anything).
          Return(nil)

        mocked.Mail.EXPECT().
          Patch(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        mailDTO: &dto.SendMailDTO{
          Subject:    gofakeit.AnimalType(),
          Recipients: []types.Email{types.Email(gofakeit.Email())},
          BodyType:   entity.MailBodyType(gofakeit.UintN(1)),
          Body:       gofakeit.LoremIpsumSentence(10),
          TagIds:     sharedUtil.GenerateMultiple(3, types.MustCreateId),
        },
      },
      want:  1,
      want1: status.Success(),
    },
    {
      name: "Success send multiple mails with tags",
      setup: func(mocked *mailMocked, arg any, want any) {
        //a := arg.(*args)
        // Create proxy
        mocked.defaultUOWMock()
        mocked.txProxy()

        mocked.Mail.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(nil)

        mocked.Mail.EXPECT().
          AppendMultipleTags(mock.Anything, mock.Anything).
          Return(nil)

        mocked.MailClient.EXPECT().
          Send(mock.Anything, mock.Anything, mock.Anything).
          Return(nil).
          Times(4)

        mocked.Mail.EXPECT().
          Patch(mock.Anything, mock.Anything).
          Return(nil).
          Times(4)
      },
      args: args{
        ctx: context.Background(),
        mailDTO: &dto.SendMailDTO{
          Subject: gofakeit.AnimalType(),
          Recipients: sharedUtil.GenerateMultiple(4, func() types.Email {
            return types.Email(gofakeit.Email())
          }),
          BodyType: entity.MailBodyType(gofakeit.UintN(1)),
          Body:     gofakeit.LoremIpsumSentence(10),
          TagIds:   sharedUtil.GenerateMultiple(3, types.MustCreateId),
        },
      },
      want:  4,
      want1: status.Success(),
    },
    {
      name: "Failed to deliver mails",
      setup: func(mocked *mailMocked, arg any, want any) {
        //a := arg.(*args)
        // Create proxy
        mocked.defaultUOWMock()
        mocked.txProxy()

        mocked.Mail.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(nil)

        mocked.Mail.EXPECT().
          AppendMultipleTags(mock.Anything, mock.Anything).
          Return(nil)

        mocked.MailClient.EXPECT().
          Send(mock.Anything, mock.Anything, mock.Anything).
          Return(nil).
          Times(2)

        mocked.Mail.EXPECT().
          Patch(mock.Anything, mock.Anything).
          Return(nil).
          Times(2)
      },
      args: args{
        ctx: context.Background(),
        mailDTO: &dto.SendMailDTO{
          Subject:    gofakeit.AnimalType(),
          Recipients: []types.Email{types.Email(gofakeit.Email()), types.Email(gofakeit.Email())},
          BodyType:   entity.MailBodyType(gofakeit.UintN(1)),
          Body:       gofakeit.LoremIpsumSentence(10),
          TagIds:     sharedUtil.GenerateMultiple(3, types.MustCreateId),
        },
      },
      want:  0,
      want1: status.Success(),
    },
    {
      name: "Failed to append tags",
      setup: func(mocked *mailMocked, arg any, want any) {
        //a := arg.(*args)
        // Create proxy
        mocked.defaultUOWMock()
        mocked.txProxy()

        mocked.Mail.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(nil)

        mocked.Mail.EXPECT().
          AppendMultipleTags(mock.Anything, mock.Anything).
          Return(dummyErr)
      },
      args: args{
        ctx: context.Background(),
        mailDTO: &dto.SendMailDTO{
          Subject:    gofakeit.AnimalType(),
          Recipients: []types.Email{types.Email(gofakeit.Email())},
          BodyType:   entity.MailBodyType(gofakeit.UintN(1)),
          Body:       gofakeit.LoremIpsumSentence(10),
          TagIds:     sharedUtil.GenerateMultiple(3, types.MustCreateId),
        },
      },
      want:  0,
      want1: status.FromRepository(dummyErr, status.NullCode),
    },
    {
      name: "Failed to save the metadata",
      setup: func(mocked *mailMocked, arg any, want any) {
        //a := arg.(*args)
        // Create proxy
        mocked.defaultUOWMock()
        mocked.txProxy()

        mocked.Mail.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(dummyErr)
      },
      args: args{
        ctx: context.Background(),
        mailDTO: &dto.SendMailDTO{
          Subject:    gofakeit.AnimalType(),
          Recipients: []types.Email{types.Email(gofakeit.Email())},
          BodyType:   entity.MailBodyType(gofakeit.UintN(1)),
          Body:       gofakeit.LoremIpsumSentence(10),
          TagIds:     sharedUtil.GenerateMultiple(3, types.MustCreateId),
        },
      },
      want:  0,
      want1: status.FromRepository(dummyErr, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newMailMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, tt.want)
      }

      m := &mailService{
        mailExt: mocked.MailClient,
        mailUow: mocked.UOW,
        tracer:  mocked.Tracer,
      }
      got, got1 := m.Send(tt.args.ctx, tt.args.mailDTO)
      require.Equal(t, len(got), tt.want)

      for {
        if !m.HasWork() {
          break
        }
        time.Sleep(time.Millisecond * 100)
      }

      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("Send() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}
*/
func Test_mailService_Update(t *testing.T) {
  type args struct {
    ctx     context.Context
    mailDTO *dto.UpdateMailDTO
  }
  tests := []struct {
    name  string
    setup setupMailTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success append tags",
      setup: func(mocked *mailMocked, arg any, want any) {
        a := arg.(*args)

        mocked.txProxy()

        mocked.Mail.EXPECT().
          AppendTags(mock.Anything, a.mailDTO.Id, a.mailDTO.AddedTagIds).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        mailDTO: &dto.UpdateMailDTO{
          Id:          types.MustCreateId(),
          AddedTagIds: []types.Id{types.MustCreateId()},
        },
      },
      want: status.Updated(),
    },
    {
      name: "Success remove tags",
      setup: func(mocked *mailMocked, arg any, want any) {
        a := arg.(*args)

        mocked.txProxy()

        mocked.Mail.EXPECT().
          RemoveTags(mock.Anything, a.mailDTO.Id, a.mailDTO.RemovedTagIds).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        mailDTO: &dto.UpdateMailDTO{
          Id:            types.MustCreateId(),
          RemovedTagIds: []types.Id{types.MustCreateId()},
        },
      },
      want: status.Updated(),
    },
    {
      name: "Success add and remove tags",
      setup: func(mocked *mailMocked, arg any, want any) {
        a := arg.(*args)

        mocked.txProxy()

        mocked.Mail.EXPECT().
          AppendTags(mock.Anything, a.mailDTO.Id, a.mailDTO.AddedTagIds).
          Return(nil)

        mocked.Mail.EXPECT().
          RemoveTags(mock.Anything, a.mailDTO.Id, a.mailDTO.RemovedTagIds).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        mailDTO: &dto.UpdateMailDTO{
          Id:            types.MustCreateId(),
          AddedTagIds:   []types.Id{types.MustCreateId()},
          RemovedTagIds: []types.Id{types.MustCreateId()},
        },
      },
      want: status.Updated(),
    },
    {
      name: "Failed to add and remove tags",
      setup: func(mocked *mailMocked, arg any, want any) {
        a := arg.(*args)

        mocked.txProxy()

        mocked.Mail.EXPECT().
          AppendTags(mock.Anything, a.mailDTO.Id, a.mailDTO.AddedTagIds).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        mailDTO: &dto.UpdateMailDTO{
          Id:            types.MustCreateId(),
          AddedTagIds:   []types.Id{types.MustCreateId()},
          RemovedTagIds: []types.Id{types.MustCreateId()},
        },
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newMailMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, &tt.want)
      }

      m := &mailService{
        mailExt: mocked.MailClient,
        mailUow: mocked.UOW,
        tracer:  mocked.Tracer,
      }

      if got := m.Update(tt.args.ctx, tt.args.mailDTO); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("Update() = %v, want %v", got, tt.want)
      }
    })
  }
}
