package service

import (
  "context"
  "database/sql"
  "errors"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  "github.com/brianvoe/gofakeit/v7"
  "github.com/stretchr/testify/mock"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/trace/noop"
  "nexa/services/mailer/internal/domain/dto"
  "nexa/services/mailer/internal/domain/entity"
  "nexa/services/mailer/internal/domain/repository/mocks"
  "reflect"
  "testing"
)

var dummyErr = errors.New("dummy error")

var dummyId = types.MustCreateId()

func newTagMocked(t *testing.T) tagMocked {
  // Tracer
  provider := noop.NewTracerProvider()

  return tagMocked{
    Tag:    mocks.NewTagMock(t),
    Tracer: provider.Tracer("MOCK"),
  }
}

type tagMocked struct {
  Tag    *mocks.TagMock
  Tracer trace.Tracer
}

type setupTagTestFunc func(mocked *tagMocked, arg any, want any)

func Test_tagService_Create(t1 *testing.T) {
  type args struct {
    ctx       context.Context
    createDto *dto.CreateTagDTO
  }
  tests := []struct {
    name     string
    setup    setupTagTestFunc
    args     args
    wantNull bool
    want1    status.Object
  }{
    {
      name: "Success create tag",
      setup: func(mocked *tagMocked, arg any, want any) {
        mocked.Tag.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        createDto: &dto.CreateTagDTO{
          Name:        gofakeit.AnimalType(),
          Description: types.SomeNullable(gofakeit.LoremIpsumSentence(5)),
        },
      },
      wantNull: false,
      want1:    status.Created(),
    },
    {
      name: "Failed create tag",
      setup: func(mocked *tagMocked, arg any, want any) {
        mocked.Tag.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        createDto: &dto.CreateTagDTO{
          Name:        gofakeit.AnimalType(),
          Description: types.SomeNullable(gofakeit.LoremIpsumSentence(5)),
        },
      },
      wantNull: true,
      want1:    status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t1.Run(tt.name, func(t1 *testing.T) {
      mocked := newTagMocked(t1)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      t := &tagService{
        tagRepo: mocked.Tag,
        tracer:  mocked.Tracer,
      }

      got, got1 := t.Create(tt.args.ctx, tt.args.createDto)

      if (got == types.NullId()) != tt.wantNull {
        t1.Errorf("Create() got1 = %v, want %v nulled", got, tt.wantNull)
      }

      if !reflect.DeepEqual(got1, tt.want1) {
        t1.Errorf("Create() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_tagService_GetAll(t1 *testing.T) {
  type args struct {
    ctx        context.Context
    elementDTO *sharedDto.PagedElementDTO
  }
  tests := []struct {
    name  string
    setup setupTagTestFunc
    args  args
    want  sharedDto.PagedElementResult[dto.TagResponseDTO]
    want1 status.Object
  }{
    {
      name: "Success get all tags",
      setup: func(mocked *tagMocked, arg any, want any) {
        a := arg.(*args)
        w := want.(*sharedDto.PagedElementResult[dto.TagResponseDTO])

        mocked.Tag.EXPECT().
          Get(mock.Anything, a.elementDTO.ToQueryParam()).
          Return(repo.PaginatedResult[entity.Tag]{nil, w.TotalElements, w.Element}, nil)
      },
      args: args{
        ctx: context.Background(),
        elementDTO: &sharedDto.PagedElementDTO{
          Element: 5,
          Page:    2,
        },
      },
      want: sharedDto.PagedElementResult[dto.TagResponseDTO]{
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
      setup: func(mocked *tagMocked, arg any, want any) {
        a := arg.(*args)
        w := want.(*sharedDto.PagedElementResult[dto.TagResponseDTO])

        mocked.Tag.EXPECT().
          Get(mock.Anything, a.elementDTO.ToQueryParam()).
          Return(repo.PaginatedResult[entity.Tag]{nil, w.TotalElements, w.Element}, sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        elementDTO: &sharedDto.PagedElementDTO{
          Element: 5,
          Page:    2,
        },
      },
      want: sharedDto.PagedElementResult[dto.TagResponseDTO]{
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
    t1.Run(tt.name, func(t1 *testing.T) {
      mocked := newTagMocked(t1)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, &tt.want)
      }

      t := &tagService{
        tagRepo: mocked.Tag,
        tracer:  mocked.Tracer,
      }
      got, got1 := t.GetAll(tt.args.ctx, tt.args.elementDTO)
      if !reflect.DeepEqual(got, tt.want) {
        t1.Errorf("GetAll() got = %v, want %v", got, tt.want)
      }
      if !reflect.DeepEqual(got1, tt.want1) {
        t1.Errorf("GetAll() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_tagService_FindByIds(t1 *testing.T) {
  type args struct {
    ctx    context.Context
    tagIds []types.Id
  }
  tests := []struct {
    name  string
    setup setupTagTestFunc
    args  args
    want  []dto.TagResponseDTO
    want1 status.Object
  }{
    {
      name: "Success find tag by the ids",
      setup: func(mocked *tagMocked, arg any, want any) {
        a := arg.(*args)
        w := want.([]dto.TagResponseDTO)

        tags := sharedUtil.CastSliceP(w, func(tag *dto.TagResponseDTO) entity.Tag {
          return entity.Tag{
            Id:          tag.Id,
            Name:        tag.Name,
            Description: tag.Description,
          }
        })

        ids := sharedUtil.CastSlice(a.tagIds, sharedUtil.ToAny[types.Id])
        mocked.Tag.EXPECT().
          FindByIds(mock.Anything, ids...).
          Return(tags, nil)
      },
      args: args{
        ctx:    context.Background(),
        tagIds: sharedUtil.GenerateMultiple(2, types.MustCreateId),
      },
      want: sharedUtil.GenerateMultiple(2, func() dto.TagResponseDTO {
        return dto.TagResponseDTO{
          Id:          types.MustCreateId(),
          Name:        gofakeit.Name(),
          Description: gofakeit.LoremIpsumSentence(20),
        }
      }),
      want1: status.Success(),
    },
    {
      name: "Failed to find tag by the ids",
      setup: func(mocked *tagMocked, arg any, want any) {
        a := arg.(*args)

        ids := sharedUtil.CastSlice(a.tagIds, sharedUtil.ToAny[types.Id])
        mocked.Tag.EXPECT().
          FindByIds(mock.Anything, ids...).
          Return(nil, sql.ErrNoRows)
      },
      args: args{
        ctx:    context.Background(),
        tagIds: sharedUtil.GenerateMultiple(2, types.MustCreateId),
      },
      want:  nil,
      want1: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t1.Run(tt.name, func(t1 *testing.T) {
      mocked := newTagMocked(t1)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, tt.want)
      }

      t := &tagService{
        tagRepo: mocked.Tag,
        tracer:  mocked.Tracer,
      }
      got, got1 := t.FindByIds(tt.args.ctx, tt.args.tagIds...)
      if !reflect.DeepEqual(got, tt.want) {
        t1.Errorf("FindByIds() got = %v, want %v", got, tt.want)
      }
      if !reflect.DeepEqual(got1, tt.want1) {
        t1.Errorf("FindByIds() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_tagService_FindByName(t1 *testing.T) {
  type args struct {
    ctx  context.Context
    name string
  }
  tests := []struct {
    name  string
    setup setupTagTestFunc
    args  args
    want  dto.TagResponseDTO
    want1 status.Object
  }{
    {
      name: "Success find tag by the name",
      setup: func(mocked *tagMocked, arg any, want any) {
        a := arg.(*args)
        w := want.(*dto.TagResponseDTO)

        mocked.Tag.EXPECT().
          FindByName(mock.Anything, a.name).
            Return(&entity.Tag{
              Id:          w.Id,
              Name:        w.Name,
              Description: w.Description,
            }, nil)
      },
      args: args{
        ctx:  context.Background(),
        name: gofakeit.AnimalType(),
      },
      want: dto.TagResponseDTO{
        Id:          types.MustCreateId(),
        Name:        gofakeit.Name(),
        Description: gofakeit.LoremIpsumSentence(20),
      },
      want1: status.Success(),
    },
    {
      name: "Failed to find tag by the name",
      setup: func(mocked *tagMocked, arg any, want any) {
        a := arg.(*args)

        mocked.Tag.EXPECT().
          FindByName(mock.Anything, a.name).
          Return(&entity.Tag{}, sql.ErrNoRows)
      },
      args: args{
        ctx:  context.Background(),
        name: gofakeit.AnimalType(),
      },
      want:  dto.TagResponseDTO{},
      want1: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t1.Run(tt.name, func(t1 *testing.T) {
      mocked := newTagMocked(t1)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, &tt.want)
      }

      t := &tagService{
        tagRepo: mocked.Tag,
        tracer:  mocked.Tracer,
      }
      got, got1 := t.FindByName(tt.args.ctx, tt.args.name)
      if !reflect.DeepEqual(got, tt.want) {
        t1.Errorf("FindByName() got = %v, want %v", got, tt.want)
      }
      if !reflect.DeepEqual(got1, tt.want1) {
        t1.Errorf("FindByName() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_tagService_Remove(t1 *testing.T) {
  type args struct {
    ctx context.Context
    id  types.Id
  }
  tests := []struct {
    name  string
    setup setupTagTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success remove tag",
      setup: func(mocked *tagMocked, arg any, want any) {
        a := arg.(*args)

        mocked.Tag.EXPECT().
          Remove(mock.Anything, a.id).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        id:  types.MustCreateId(),
      },
      want: status.Deleted(),
    },
    {
      name: "Failed to remove tag",
      setup: func(mocked *tagMocked, arg any, want any) {
        a := arg.(*args)

        mocked.Tag.EXPECT().
          Remove(mock.Anything, a.id).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        id:  types.MustCreateId(),
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t1.Run(tt.name, func(t1 *testing.T) {
      mocked := newTagMocked(t1)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      t := &tagService{
        tagRepo: mocked.Tag,
        tracer:  mocked.Tracer,
      }
      if got := t.Remove(tt.args.ctx, tt.args.id); !reflect.DeepEqual(got, tt.want) {
        t1.Errorf("Remove() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_tagService_Update(t1 *testing.T) {
  type args struct {
    ctx       context.Context
    updateDto *dto.UpdateTagDTO
  }
  tests := []struct {
    name  string
    setup setupTagTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success update tag",
      setup: func(mocked *tagMocked, arg any, want any) {
        a := arg.(*args)

        mocked.Tag.EXPECT().
            Patch(mock.Anything, &entity.PatchedTag{
              Id:          a.updateDto.Id,
              Name:        a.updateDto.Name.ValueOr(""),
              Description: a.updateDto.Description,
            }).Return(nil)

      },
      args: args{
        ctx: context.Background(),
        updateDto: &dto.UpdateTagDTO{
          Id:          types.MustCreateId(),
          Name:        types.SomeNullable(gofakeit.AnimalType()),
          Description: types.SomeNullable(gofakeit.LoremIpsumSentence(5)),
        },
      },
      want: status.Updated(),
    },
    {
      name: "Failed to update tag",
      setup: func(mocked *tagMocked, arg any, want any) {
        a := arg.(*args)

        mocked.Tag.EXPECT().
            Patch(mock.Anything, &entity.PatchedTag{
              Id:          a.updateDto.Id,
              Name:        a.updateDto.Name.ValueOr(""),
              Description: a.updateDto.Description,
            }).Return(sql.ErrNoRows)

      },
      args: args{
        ctx: context.Background(),
        updateDto: &dto.UpdateTagDTO{
          Id:          types.MustCreateId(),
          Name:        types.SomeNullable(gofakeit.AnimalType()),
          Description: types.SomeNullable(gofakeit.LoremIpsumSentence(5)),
        },
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t1.Run(tt.name, func(t1 *testing.T) {
      mocked := newTagMocked(t1)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      t := &tagService{
        tagRepo: mocked.Tag,
        tracer:  mocked.Tracer,
      }
      if got := t.Update(tt.args.ctx, tt.args.updateDto); !reflect.DeepEqual(got, tt.want) {
        t1.Errorf("Update() = %v, want %v", got, tt.want)
      }
    })
  }
}
