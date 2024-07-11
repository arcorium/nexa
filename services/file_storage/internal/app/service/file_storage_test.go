package service

import (
  "context"
  "database/sql"
  "errors"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/uow"
  uowMock "github.com/arcorium/nexa/shared/uow/mocks"
  "github.com/brianvoe/gofakeit/v7"
  "github.com/stretchr/testify/mock"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/trace/noop"
  uow2 "nexa/services/file_storage/internal/app/uow"
  "nexa/services/file_storage/internal/domain/dto"
  "nexa/services/file_storage/internal/domain/entity"
  extMock "nexa/services/file_storage/internal/domain/external/mocks"
  repoMock "nexa/services/file_storage/internal/domain/repository/mocks"
  "reflect"
  "testing"
)

var dummyErr = errors.New("dummy error")

var dummyId = types.MustCreateId()

func newStorageMocked(t *testing.T) storageMocked {
  // Tracer
  provider := noop.NewTracerProvider()
  return storageMocked{
    UOW:      uowMock.NewUnitOfWorkMock[uow2.FileMetadataStorage](t),
    Storage:  extMock.NewStorageMock(t),
    Metadata: repoMock.NewFileMetadataMock(t),
    Tracer:   provider.Tracer("MOCK"),
  }
}

type storageMocked struct {
  UOW      *uowMock.UnitOfWorkMock[uow2.FileMetadataStorage]
  Storage  *extMock.StorageMock
  Metadata *repoMock.FileMetadataMock
  Tracer   trace.Tracer
}

func (m *storageMocked) defaultUOWMock() {
  m.UOW.EXPECT().
    Repositories().
    Return(uow2.NewStorage(m.Metadata))
}

func (m *storageMocked) txProxy() {
  m.UOW.On("DoTx", mock.Anything, mock.Anything).
      Return(func(ctx context.Context, f uow.UOWBlock[uow2.FileMetadataStorage]) error {
        return f(ctx, uow2.NewStorage(m.Metadata))
      })
}

type setupStorageTestFunc func(mocked *storageMocked, arg any, want any)

func Test_fileStorage_Delete(t *testing.T) {
  type args struct {
    ctx context.Context
    id  types.Id
  }
  tests := []struct {
    name  string
    setup setupStorageTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success delete file",
      setup: func(mocked *storageMocked, arg any, want any) {
        a := arg.(*args)

        mocked.defaultUOWMock()
        mocked.txProxy()

        metadata := entity.FileMetadata{
          Id:           types.MustCreateId(),
          Name:         gofakeit.AppName(),
          MimeType:     gofakeit.FileMimeType(),
          Size:         gofakeit.Uint64(),
          IsPublic:     gofakeit.Bool(),
          Provider:     entity.StorageProvider(gofakeit.UintN(1)),
          ProviderPath: gofakeit.URL(),
        }

        mocked.Metadata.EXPECT().
          FindByIds(mock.Anything, a.id).
          Return([]entity.FileMetadata{metadata}, nil)

        mocked.Metadata.EXPECT().
          DeleteById(mock.Anything, a.id).
          Return(nil)

        mocked.Storage.EXPECT().
          Delete(mock.Anything, metadata.ProviderPath).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        id:  types.MustCreateId(),
      },
      want: status.Deleted(),
    },
    {
      name: "File not found",
      setup: func(mocked *storageMocked, arg any, want any) {
        a := arg.(*args)

        mocked.defaultUOWMock()

        metadata := entity.FileMetadata{
          Id:           types.MustCreateId(),
          Name:         gofakeit.AppName(),
          MimeType:     gofakeit.FileMimeType(),
          Size:         gofakeit.Uint64(),
          IsPublic:     gofakeit.Bool(),
          Provider:     entity.StorageProvider(gofakeit.UintN(1)),
          ProviderPath: gofakeit.URL(),
        }

        mocked.Metadata.EXPECT().
          FindByIds(mock.Anything, a.id).
          Return([]entity.FileMetadata{metadata}, sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        id:  types.MustCreateId(),
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
    {
      name: "Failed to delete file in the storage",
      setup: func(mocked *storageMocked, arg any, want any) {
        a := arg.(*args)

        mocked.defaultUOWMock()
        mocked.txProxy()

        metadata := entity.FileMetadata{
          Id:           types.MustCreateId(),
          Name:         gofakeit.AppName(),
          MimeType:     gofakeit.FileMimeType(),
          Size:         gofakeit.Uint64(),
          IsPublic:     gofakeit.Bool(),
          Provider:     entity.StorageProvider(gofakeit.UintN(1)),
          ProviderPath: gofakeit.URL(),
        }

        mocked.Metadata.EXPECT().
          FindByIds(mock.Anything, a.id).
          Return([]entity.FileMetadata{metadata}, nil)

        mocked.Metadata.EXPECT().
          DeleteById(mock.Anything, a.id).
          Return(nil)

        mocked.Storage.EXPECT().
          Delete(mock.Anything, metadata.ProviderPath).
          Return(dummyErr)
      },
      args: args{
        ctx: context.Background(),
        id:  types.MustCreateId(),
      },
      want: status.ErrExternal(dummyErr),
    },
    {
      name: "Failed to delete metadata file",
      setup: func(mocked *storageMocked, arg any, want any) {
        a := arg.(*args)

        mocked.defaultUOWMock()
        mocked.txProxy()

        metadata := entity.FileMetadata{
          Id:           types.MustCreateId(),
          Name:         gofakeit.AppName(),
          MimeType:     gofakeit.FileMimeType(),
          Size:         gofakeit.Uint64(),
          IsPublic:     gofakeit.Bool(),
          Provider:     entity.StorageProvider(gofakeit.UintN(1)),
          ProviderPath: gofakeit.URL(),
        }

        mocked.Metadata.EXPECT().
          FindByIds(mock.Anything, a.id).
          Return([]entity.FileMetadata{metadata}, nil)

        mocked.Metadata.EXPECT().
          DeleteById(mock.Anything, a.id).
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
    t.Run(tt.name, func(t *testing.T) {
      mocked := newStorageMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, &tt.want)
      }

      f := &fileStorage{
        unit:       mocked.UOW,
        storageExt: mocked.Storage,
        tracer:     mocked.Tracer,
      }
      if got := f.Delete(tt.args.ctx, tt.args.id); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("Delete() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_fileStorage_Find(t *testing.T) {
  type args struct {
    ctx context.Context
    id  types.Id
  }
  tests := []struct {
    name  string
    setup setupStorageTestFunc
    args  args
    want  dto.FileResponseDTO
    want1 status.Object
  }{
    {
      name: "File does exists",
      setup: func(mocked *storageMocked, arg any, want any) {
        a := arg.(*args)
        w := want.(*dto.FileResponseDTO)

        mocked.defaultUOWMock()

        metadata := entity.FileMetadata{
          Id:           types.MustCreateId(),
          Name:         w.Name,
          MimeType:     gofakeit.FileMimeType(),
          Size:         w.Size,
          IsPublic:     gofakeit.Bool(),
          Provider:     entity.StorageProvider(gofakeit.UintN(1)),
          ProviderPath: gofakeit.URL(),
        }

        file := entity.File{
          Name:     metadata.Name,
          Bytes:    nil,
          Size:     metadata.Size,
          IsPublic: metadata.IsPublic,
        }

        mocked.Metadata.EXPECT().
          FindByIds(mock.Anything, a.id).
          Return([]entity.FileMetadata{metadata}, nil)

        mocked.Storage.EXPECT().
          Find(mock.Anything, metadata.ProviderPath).
          Return(file, nil)
      },
      args: args{
        ctx: context.Background(),
        id:  types.MustCreateId(),
      },
      want: dto.FileResponseDTO{
        Name: gofakeit.AppName(),
        Size: gofakeit.Uint64(),
        Data: nil,
      },
      want1: status.Success(),
    },
    {
      name: "File doesn't exists on storage",
      setup: func(mocked *storageMocked, arg any, want any) {
        a := arg.(*args)
        w := want.(*dto.FileResponseDTO)

        mocked.defaultUOWMock()

        metadata := entity.FileMetadata{
          Id:           types.MustCreateId(),
          Name:         w.Name,
          MimeType:     gofakeit.FileMimeType(),
          Size:         w.Size,
          IsPublic:     gofakeit.Bool(),
          Provider:     entity.StorageProvider(gofakeit.UintN(1)),
          ProviderPath: gofakeit.URL(),
        }

        mocked.Metadata.EXPECT().
          FindByIds(mock.Anything, a.id).
          Return([]entity.FileMetadata{metadata}, nil)

        mocked.Storage.EXPECT().
          Find(mock.Anything, metadata.ProviderPath).
          Return(entity.File{}, dummyErr)
      },
      args: args{
        ctx: context.Background(),
        id:  types.MustCreateId(),
      },
      want:  dto.FileResponseDTO{},
      want1: status.ErrExternal(dummyErr),
    },
    {
      name: "File metadata doesn't exist",
      setup: func(mocked *storageMocked, arg any, want any) {
        a := arg.(*args)

        mocked.defaultUOWMock()

        mocked.Metadata.EXPECT().
          FindByIds(mock.Anything, a.id).
          Return(nil, sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        id:  types.MustCreateId(),
      },
      want:  dto.FileResponseDTO{},
      want1: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newStorageMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, &tt.want)
      }

      f := &fileStorage{
        unit:       mocked.UOW,
        storageExt: mocked.Storage,
        tracer:     mocked.Tracer,
      }
      got, got1 := f.Find(tt.args.ctx, tt.args.id)
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("Find() got = %v, want %v", got, tt.want)
      }
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("Find() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_fileStorage_FindMetadata(t *testing.T) {
  type args struct {
    ctx context.Context
    id  types.Id
  }
  tests := []struct {
    name  string
    setup setupStorageTestFunc
    args  args
    want  *dto.FileMetadataResponseDTO
    want1 status.Object
  }{
    {
      name: "File metadata does exists",
      setup: func(mocked *storageMocked, arg any, want any) {
        a := arg.(*args)
        w := want.(*dto.FileMetadataResponseDTO)

        mocked.defaultUOWMock()

        metadata := entity.FileMetadata{
          Id:       w.Id,
          Name:     w.Name,
          MimeType: gofakeit.FileMimeType(),
          Size:     w.Size,
          IsPublic: gofakeit.Bool(),
          Provider: entity.StorageProvider(gofakeit.UintN(1)),
          FullPath: w.Path.Path(),
        }

        mocked.Metadata.EXPECT().
          FindByIds(mock.Anything, a.id).
          Return([]entity.FileMetadata{metadata}, nil)
      },
      args: args{
        ctx: context.Background(),
        id:  dummyId,
      },
      want: &dto.FileMetadataResponseDTO{
        Id:   dummyId,
        Name: gofakeit.AppName(),
        Size: gofakeit.Uint64(),
        Path: types.FilePath(gofakeit.URL()),
      },
      want1: status.Success(),
    },
    {
      name: "File metadata doesn't exists",
      setup: func(mocked *storageMocked, arg any, want any) {
        a := arg.(*args)

        mocked.defaultUOWMock()

        mocked.Metadata.EXPECT().
          FindByIds(mock.Anything, a.id).
          Return(nil, sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        id:  dummyId,
      },
      want:  nil,
      want1: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newStorageMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, tt.want)
      }

      f := &fileStorage{
        unit:       mocked.UOW,
        storageExt: mocked.Storage,
        tracer:     mocked.Tracer,
      }
      got, got1 := f.FindMetadata(tt.args.ctx, tt.args.id)
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindMetadata() got = %v, want %v", got, tt.want)
      }
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("FindMetadata() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_fileStorage_Store(t *testing.T) {
  type args struct {
    ctx     context.Context
    fileDto *dto.FileStoreDTO
  }
  tests := []struct {
    name     string
    setup    setupStorageTestFunc
    args     args
    wantNull bool
    want1    status.Object
  }{
    {
      name: "Success store public file",
      setup: func(mocked *storageMocked, arg any, want any) {
        mocked.defaultUOWMock()
        mocked.txProxy()

        mocked.Storage.EXPECT().
          GetProvider().Return(entity.StorageProviderAWSS3)

        relativePath := gofakeit.URL()
        mocked.Storage.EXPECT().
          Store(mock.Anything, mock.Anything).
          Return(relativePath, nil)

        mocked.Storage.EXPECT().
          GetFullPath(mock.Anything, relativePath).
          Return(types.FilePath(gofakeit.URL()), nil)

        mocked.Metadata.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        fileDto: &dto.FileStoreDTO{
          Name:     gofakeit.AppName(),
          Data:     nil,
          IsPublic: true,
        },
      },
      wantNull: false,
      want1:    status.Success(),
    },
    {
      name: "Success store private file",
      setup: func(mocked *storageMocked, arg any, want any) {
        mocked.defaultUOWMock()
        mocked.txProxy()

        mocked.Storage.EXPECT().
          GetProvider().Return(entity.StorageProviderAWSS3)

        mocked.Storage.EXPECT().
          Store(mock.Anything, mock.Anything).
          Return(gofakeit.URL(), nil)

        mocked.Metadata.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        fileDto: &dto.FileStoreDTO{
          Name:     gofakeit.AppName(),
          Data:     nil,
          IsPublic: false,
        },
      },
      wantNull: false,
      want1:    status.Success(),
    },
    {
      name: "Failed to store file in storage",
      setup: func(mocked *storageMocked, arg any, want any) {
        mocked.txProxy()

        mocked.Storage.EXPECT().
          GetProvider().Return(entity.StorageProviderAWSS3)

        path := gofakeit.URL()
        mocked.Storage.EXPECT().
          Store(mock.Anything, mock.Anything).
          Return(path, dummyErr)
      },
      args: args{
        ctx: context.Background(),
        fileDto: &dto.FileStoreDTO{
          Name:     gofakeit.AppName(),
          Data:     nil,
          IsPublic: false,
        },
      },
      wantNull: true,
      want1:    status.ErrExternal(dummyErr),
    },
    {
      name: "Failed to upload file metadata",
      setup: func(mocked *storageMocked, arg any, want any) {
        mocked.defaultUOWMock()
        mocked.txProxy()

        mocked.Storage.EXPECT().
          GetProvider().Return(entity.StorageProviderAWSS3)

        mocked.Storage.EXPECT().
          Store(mock.Anything, mock.Anything).
          Return(gofakeit.URL(), nil)

        mocked.Metadata.EXPECT().
          Create(mock.Anything, mock.Anything).
          Return(sql.ErrNoRows)

        mocked.Storage.EXPECT().
          Delete(mock.Anything, mock.Anything).
          Return(nil)

      },
      args: args{
        ctx: context.Background(),
        fileDto: &dto.FileStoreDTO{
          Name:     gofakeit.AppName(),
          Data:     nil,
          IsPublic: false,
        },
      },
      wantNull: true,
      want1:    status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newStorageMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      f := &fileStorage{
        unit:       mocked.UOW,
        storageExt: mocked.Storage,
        tracer:     mocked.Tracer,
      }

      got, got1 := f.Store(tt.args.ctx, tt.args.fileDto)
      if (got == types.NullId()) != tt.wantNull {
        t.Errorf("Store() got = %v, want %v", got, tt.wantNull)
      }

      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("Store() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_fileStorage_Move(t *testing.T) {
  type args struct {
    ctx   context.Context
    input *dto.UpdateFileMetadataDTO
  }
  tests := []struct {
    name  string
    setup setupStorageTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success moving file into public path",
      setup: func(mocked *storageMocked, arg any, want any) {
        a := arg.(*args)

        mocked.defaultUOWMock()
        mocked.txProxy()

        metadata := entity.FileMetadata{
          Id:           types.MustCreateId(),
          Name:         gofakeit.AppName(),
          MimeType:     gofakeit.FileMimeType(),
          Size:         gofakeit.Uint64(),
          IsPublic:     false,
          Provider:     entity.StorageProvider(gofakeit.UintN(1)),
          ProviderPath: gofakeit.URL(),
        }

        mocked.Metadata.EXPECT().
          FindByIds(mock.Anything, a.input.Id).
          Return([]entity.FileMetadata{metadata}, nil)

        newPath := gofakeit.URL()

        mocked.Storage.EXPECT().
          Copy(mock.Anything, metadata.ProviderPath, mock.Anything).
          Return(newPath, nil)

        mocked.Storage.EXPECT().
          GetFullPath(mock.Anything, newPath).
          Return(types.FilePath(gofakeit.URL()), nil)

        mocked.Metadata.EXPECT().
          Patch(mock.Anything, mock.Anything).
          Return(nil)

        mocked.Storage.EXPECT().
          Delete(mock.Anything, metadata.ProviderPath).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        input: &dto.UpdateFileMetadataDTO{
          Id:       types.MustCreateId(),
          IsPublic: true,
        },
      },
      want: status.Updated(),
    },
    {
      name: "Success moving file from public path",
      setup: func(mocked *storageMocked, arg any, want any) {
        a := arg.(*args)

        mocked.defaultUOWMock()
        mocked.txProxy()

        metadata := entity.FileMetadata{
          Id:           types.MustCreateId(),
          Name:         gofakeit.AppName(),
          MimeType:     gofakeit.FileMimeType(),
          Size:         gofakeit.Uint64(),
          IsPublic:     true,
          Provider:     entity.StorageProvider(gofakeit.UintN(1)),
          ProviderPath: gofakeit.URL(),
        }

        mocked.Metadata.EXPECT().
          FindByIds(mock.Anything, a.input.Id).
          Return([]entity.FileMetadata{metadata}, nil)

        newPath := gofakeit.URL()

        mocked.Storage.EXPECT().
          Copy(mock.Anything, metadata.ProviderPath, mock.Anything).
          Return(newPath, nil)

        mocked.Metadata.EXPECT().
          Patch(mock.Anything, mock.Anything).
          Return(nil)

        mocked.Storage.EXPECT().
          Delete(mock.Anything, metadata.ProviderPath).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        input: &dto.UpdateFileMetadataDTO{
          Id:       types.MustCreateId(),
          IsPublic: false,
        },
      },
      want: status.Updated(),
    },
    {
      name: "Moving public file to public location",
      setup: func(mocked *storageMocked, arg any, want any) {
        a := arg.(*args)

        mocked.defaultUOWMock()

        metadata := entity.FileMetadata{
          Id:           types.MustCreateId(),
          Name:         gofakeit.AppName(),
          MimeType:     gofakeit.FileMimeType(),
          Size:         gofakeit.Uint64(),
          IsPublic:     true,
          Provider:     entity.StorageProvider(gofakeit.UintN(1)),
          ProviderPath: gofakeit.URL(),
        }

        mocked.Metadata.EXPECT().
          FindByIds(mock.Anything, a.input.Id).
          Return([]entity.FileMetadata{metadata}, nil)

      },
      args: args{
        ctx: context.Background(),
        input: &dto.UpdateFileMetadataDTO{
          Id:       types.MustCreateId(),
          IsPublic: true,
        },
      },
      want: status.New(status.OBJECT_NOT_FOUND, errors.New("nothings to do, the file is already on right location")),
    },
    {
      name: "File not found",
      setup: func(mocked *storageMocked, arg any, want any) {
        a := arg.(*args)

        mocked.defaultUOWMock()

        mocked.Metadata.EXPECT().
          FindByIds(mock.Anything, a.input.Id).
          Return(nil, sql.ErrNoRows)

      },
      args: args{
        ctx: context.Background(),
        input: &dto.UpdateFileMetadataDTO{
          Id:       types.MustCreateId(),
          IsPublic: false,
        },
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
    {
      name: "Failed to copy file in storage",
      setup: func(mocked *storageMocked, arg any, want any) {
        a := arg.(*args)

        mocked.defaultUOWMock()
        mocked.txProxy()

        metadata := entity.FileMetadata{
          Id:           types.MustCreateId(),
          Name:         gofakeit.AppName(),
          MimeType:     gofakeit.FileMimeType(),
          Size:         gofakeit.Uint64(),
          IsPublic:     false,
          Provider:     entity.StorageProvider(gofakeit.UintN(1)),
          ProviderPath: gofakeit.URL(),
        }

        mocked.Metadata.EXPECT().
          FindByIds(mock.Anything, a.input.Id).
          Return([]entity.FileMetadata{metadata}, nil)

        mocked.Storage.EXPECT().
          Copy(mock.Anything, metadata.ProviderPath, mock.Anything).
          Return("", dummyErr)

      },
      args: args{
        ctx: context.Background(),
        input: &dto.UpdateFileMetadataDTO{
          Id:       types.MustCreateId(),
          IsPublic: true,
        },
      },
      want: status.ErrExternal(dummyErr),
    },
    {
      name: "Failed to patch provider path",
      setup: func(mocked *storageMocked, arg any, want any) {
        a := arg.(*args)

        mocked.defaultUOWMock()
        mocked.txProxy()

        metadata := entity.FileMetadata{
          Id:           types.MustCreateId(),
          Name:         gofakeit.AppName(),
          MimeType:     gofakeit.FileMimeType(),
          Size:         gofakeit.Uint64(),
          IsPublic:     true,
          Provider:     entity.StorageProvider(gofakeit.UintN(1)),
          ProviderPath: gofakeit.URL(),
        }

        mocked.Metadata.EXPECT().
          FindByIds(mock.Anything, a.input.Id).
          Return([]entity.FileMetadata{metadata}, nil)

        newPath := gofakeit.URL()

        mocked.Storage.EXPECT().
          Copy(mock.Anything, metadata.ProviderPath, mock.Anything).
          Return(newPath, nil)

        mocked.Metadata.EXPECT().
          Patch(mock.Anything, mock.Anything).
          Return(sql.ErrNoRows)

        mocked.Storage.EXPECT().
          Delete(mock.Anything, newPath).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        input: &dto.UpdateFileMetadataDTO{
          Id:       types.MustCreateId(),
          IsPublic: false,
        },
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
    {
      name: "Failed to delete last file",
      setup: func(mocked *storageMocked, arg any, want any) {
        a := arg.(*args)

        mocked.defaultUOWMock()
        mocked.txProxy()

        metadata := entity.FileMetadata{
          Id:           types.MustCreateId(),
          Name:         gofakeit.AppName(),
          MimeType:     gofakeit.FileMimeType(),
          Size:         gofakeit.Uint64(),
          IsPublic:     true,
          Provider:     entity.StorageProvider(gofakeit.UintN(1)),
          ProviderPath: gofakeit.URL(),
        }

        mocked.Metadata.EXPECT().
          FindByIds(mock.Anything, a.input.Id).
          Return([]entity.FileMetadata{metadata}, nil)

        newPath := gofakeit.URL()

        mocked.Storage.EXPECT().
          Copy(mock.Anything, metadata.ProviderPath, mock.Anything).
          Return(newPath, nil)

        mocked.Metadata.EXPECT().
          Patch(mock.Anything, mock.Anything).
          Return(nil)

        mocked.Storage.EXPECT().
          Delete(mock.Anything, metadata.ProviderPath).
          Return(dummyErr)

        mocked.Storage.EXPECT().
          Delete(mock.Anything, newPath).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        input: &dto.UpdateFileMetadataDTO{
          Id:       types.MustCreateId(),
          IsPublic: false,
        },
      },
      want: status.ErrExternal(dummyErr),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newStorageMocked(t)

      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, nil)
      }

      f := &fileStorage{
        unit:       mocked.UOW,
        storageExt: mocked.Storage,
        tracer:     mocked.Tracer,
      }
      if got := f.Move(tt.args.ctx, tt.args.input); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("Move() = %v, want %v", got, tt.want)
      }
    })
  }
}
