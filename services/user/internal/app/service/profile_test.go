package service

import (
  "context"
  "database/sql"
  "fmt"
  "github.com/brianvoe/gofakeit/v7"
  "github.com/stretchr/testify/mock"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/trace/noop"
  "nexa/services/user/internal/domain/dto"
  "nexa/services/user/internal/domain/entity"
  extMock "nexa/services/user/internal/domain/external/mocks"
  repoMock "nexa/services/user/internal/domain/repository/mocks"
  "nexa/shared/status"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  "reflect"
  "testing"
)

func newProfileMocked(t *testing.T) profileMocked {
  // Tracer
  provider := noop.NewTracerProvider()

  return profileMocked{
    Profile: repoMock.NewProfileMock(t),
    Storage: extMock.NewFileStorageClientMock(t),
    Tracer:  provider.Tracer("MOCK"),
  }
}

type profileMocked struct {
  Profile *repoMock.ProfileMock
  Storage *extMock.FileStorageClientMock
  Tracer  trace.Tracer
}

type setupProfileTestFunc func(mocked *profileMocked, args any, want any)

func Test_profileService_Find(t *testing.T) {
  type args struct {
    ctx     context.Context
    userIds []types.Id
  }
  tests := []struct {
    name  string
    setup setupProfileTestFunc
    args  args
    want  []dto.ProfileResponseDTO
    want1 status.Object
  }{
    {
      name: "Success find single profile",
      setup: func(mocked *profileMocked, arg any, want any) {
        args := arg.(*args)
        expected := want.([]dto.ProfileResponseDTO)

        userIds := sharedUtil.CastSlice(args.userIds, sharedUtil.ToAny[types.Id])

        result := sharedUtil.CastSliceP(expected, func(from *dto.ProfileResponseDTO) entity.Profile {
          return entity.Profile{
            Id:        from.Id,
            FirstName: from.FirstName,
            LastName:  from.LastName,
            Bio:       from.Bio,
            PhotoURL:  from.PhotoURL,
          }
        })

        mocked.Profile.EXPECT().
          FindByIds(mock.Anything, userIds...).
          Return(result, nil)
      },
      args: args{
        ctx:     context.Background(),
        userIds: []types.Id{types.MustCreateId()},
      },
      want: []dto.ProfileResponseDTO{
        {
          Id:        types.MustCreateId(),
          FirstName: gofakeit.FirstName(),
          LastName:  gofakeit.LastName(),
          Bio:       gofakeit.LoremIpsumParagraph(1, 3, 20, "."),
          PhotoURL:  types.FilePathFromString(gofakeit.URL()),
        },
      },
      want1: status.Success(),
    },
    {
      name: "Success find multiple profile",
      setup: func(mocked *profileMocked, arg any, want any) {
        args := arg.(*args)
        expected := want.([]dto.ProfileResponseDTO)

        userIds := sharedUtil.CastSlice(args.userIds, sharedUtil.ToAny[types.Id])
        result := sharedUtil.CastSliceP(expected, func(from *dto.ProfileResponseDTO) entity.Profile {
          return entity.Profile{
            Id:        from.Id,
            FirstName: from.FirstName,
            LastName:  from.LastName,
            Bio:       from.Bio,
            PhotoURL:  from.PhotoURL,
          }
        })

        mocked.Profile.EXPECT().
          FindByIds(mock.Anything, userIds...).
          Return(result, nil)
      },
      args: args{
        ctx:     context.Background(),
        userIds: []types.Id{types.MustCreateId(), types.MustCreateId()},
      },
      want: sharedUtil.GenerateMultiple(2, func() dto.ProfileResponseDTO {
        return dto.ProfileResponseDTO{
          Id:        types.MustCreateId(),
          FirstName: gofakeit.FirstName(),
          LastName:  gofakeit.LastName(),
          Bio:       gofakeit.LoremIpsumParagraph(1, 3, 20, "."),
          PhotoURL:  types.FilePathFromString(gofakeit.URL()),
        }
      }),
      want1: status.Success(),
    },
    {
      name: "Profile not found",
      setup: func(mocked *profileMocked, arg any, want any) {
        args := arg.(*args)
        userIds := sharedUtil.CastSlice(args.userIds, sharedUtil.ToAny[types.Id])

        mocked.Profile.EXPECT().
          FindByIds(mock.Anything, userIds...).
          Return(nil, sql.ErrNoRows)
      },
      args: args{
        ctx:     context.Background(),
        userIds: []types.Id{types.MustCreateId()},
      },
      want:  nil,
      want1: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newProfileMocked(t)
      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, tt.want)
      }

      p := profileService{
        profileRepo: mocked.Profile,
        storageExt:  mocked.Storage,
        tracer:      mocked.Tracer,
      }

      got, got1 := p.Find(tt.args.ctx, tt.args.userIds...)
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("Find() got = %v, want %v", got, tt.want)
      }
      if !reflect.DeepEqual(got1, tt.want1) {
        t.Errorf("Find() got1 = %v, want %v", got1, tt.want1)
      }
    })
  }
}

func Test_profileService_Update(t *testing.T) {
  type args struct {
    ctx   context.Context
    input *dto.ProfileUpdateDTO
  }
  tests := []struct {
    name  string
    setup setupProfileTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success update profile",
      setup: func(mocked *profileMocked, arg any, want any) {
        a := arg.(*args)

        patchedProfile := a.input.ToDomain()

        mocked.Profile.EXPECT().
          Patch(mock.Anything, &patchedProfile).
          Return(nil)
      },
      args: args{
        ctx: context.Background(),
        input: &dto.ProfileUpdateDTO{
          Id:        types.MustCreateId(),
          FirstName: types.Nullable[string]{},
          LastName:  types.Nullable[string]{},
          Bio:       types.Nullable[string]{},
        },
      },
      want: status.Updated(),
    },
    {
      name: "Failed update profile",
      setup: func(mocked *profileMocked, arg any, want any) {
        a := arg.(*args)

        patchedProfile := a.input.ToDomain()

        mocked.Profile.EXPECT().
          Patch(mock.Anything, &patchedProfile).
          Return(sql.ErrNoRows)
      },
      args: args{
        ctx: context.Background(),
        input: &dto.ProfileUpdateDTO{
          Id:        types.MustCreateId(),
          FirstName: types.Nullable[string]{},
          LastName:  types.Nullable[string]{},
          Bio:       types.Nullable[string]{},
        },
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newProfileMocked(t)
      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, tt.want)
      }

      p := profileService{
        profileRepo: mocked.Profile,
        storageExt:  mocked.Storage,
        tracer:      mocked.Tracer,
      }

      if got := p.Update(tt.args.ctx, tt.args.input); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("Update() = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_profileService_UpdateAvatar(t *testing.T) {
  type args struct {
    ctx   context.Context
    input *dto.ProfileAvatarUpdateDTO
  }
  tests := []struct {
    name  string
    setup setupProfileTestFunc
    args  args
    want  status.Object
  }{
    {
      name: "Success update new avatar",
      setup: func(mocked *profileMocked, arg any, want any) {
        a := arg.(*args)

        mocked.Profile.EXPECT().
          FindByIds(mock.Anything, a.input.Id).
            Return([]entity.Profile{
              {
                Id:        a.input.Id,
                UserId:    types.MustCreateId(),
                FirstName: gofakeit.FirstName(),
                LastName:  gofakeit.LastName(),
                Bio:       gofakeit.LoremIpsumParagraph(1, 3, 20, "."),
              },
            }, nil).Once()

        profileId := types.MustCreateId()
        profilePath := types.FilePathFromString(gofakeit.URL())

        mocked.Storage.EXPECT().
            UploadProfileImage(mock.Anything, &dto.UploadImageDTO{
              Filename: a.input.Filename,
              Data:     a.input.Bytes,
            }).
          Return(profileId, profilePath, nil).
          Once()

        mocked.Profile.EXPECT().
            Patch(mock.Anything, &entity.PatchedProfile{
              Id:       a.input.Id,
              PhotoId:  types.SomeNullable(profileId),
              PhotoURL: types.SomeNullable(profilePath),
            }).
          Return(nil).
          Once()

      },
      args: args{
        ctx: context.Background(),
        input: &dto.ProfileAvatarUpdateDTO{
          Id:       types.MustCreateId(),
          Filename: fmt.Sprintf("%s.%s", gofakeit.AppName(), gofakeit.FileExtension()),
          Bytes:    []byte{0x12, 0x12, 0x12},
        },
      },
      want: status.Updated(),
    },
    {
      name: "Success replace avatar",
      setup: func(mocked *profileMocked, arg any, want any) {
        a := arg.(*args)

        result := entity.Profile{
          Id:        a.input.Id,
          UserId:    types.MustCreateId(),
          FirstName: gofakeit.FirstName(),
          LastName:  gofakeit.LastName(),
          Bio:       gofakeit.LoremIpsumParagraph(1, 3, 20, "."),
          PhotoId:   types.MustCreateId(),
          PhotoURL:  types.FilePathFromString(gofakeit.URL()),
        }
        mocked.Profile.EXPECT().
          FindByIds(mock.Anything, a.input.Id).
          Return([]entity.Profile{result}, nil).Once()

        imageId := types.MustCreateId()
        imagePath := types.FilePathFromString(gofakeit.URL())

        mocked.Storage.EXPECT().
            UploadProfileImage(mock.Anything, &dto.UploadImageDTO{
              Filename: a.input.Filename,
              Data:     a.input.Bytes,
            }).
          Return(imageId, imagePath, nil).
          Once()

        mocked.Profile.EXPECT().
            Patch(mock.Anything, &entity.PatchedProfile{
              Id:       a.input.Id,
              PhotoId:  types.SomeNullable(imageId),
              PhotoURL: types.SomeNullable(imagePath),
            }).
          Return(nil).
          Once()

        mocked.Storage.EXPECT().
          DeleteProfileImage(mock.Anything, result.PhotoId).
          Return(nil).
          Once()
      },
      args: args{
        ctx: context.Background(),
        input: &dto.ProfileAvatarUpdateDTO{
          Id:       types.MustCreateId(),
          Filename: fmt.Sprintf("%s.%s", gofakeit.AppName(), gofakeit.FileExtension()),
          Bytes:    []byte{0x12, 0x12, 0x12},
        },
      },
      want: status.Updated(),
    },
    {
      name: "Failed to set photo data into profile",
      setup: func(mocked *profileMocked, arg any, want any) {
        a := arg.(*args)

        result := entity.Profile{
          Id:        a.input.Id,
          UserId:    types.MustCreateId(),
          FirstName: gofakeit.FirstName(),
          LastName:  gofakeit.LastName(),
          Bio:       gofakeit.LoremIpsumParagraph(1, 3, 20, "."),
          PhotoId:   types.MustCreateId(),
          PhotoURL:  types.FilePathFromString(gofakeit.URL()),
        }
        mocked.Profile.EXPECT().
          FindByIds(mock.Anything, a.input.Id).
          Return([]entity.Profile{result}, nil).Once()

        imageId := types.MustCreateId()
        imagePath := types.FilePathFromString(gofakeit.URL())

        mocked.Storage.EXPECT().
            UploadProfileImage(mock.Anything, &dto.UploadImageDTO{
              Filename: a.input.Filename,
              Data:     a.input.Bytes,
            }).
          Return(imageId, imagePath, nil).
          Once()

        mocked.Profile.EXPECT().
            Patch(mock.Anything, &entity.PatchedProfile{
              Id:       a.input.Id,
              PhotoId:  types.SomeNullable(imageId),
              PhotoURL: types.SomeNullable(imagePath),
            }).
          Return(sql.ErrNoRows).
          Once()

        mocked.Storage.EXPECT().
          DeleteProfileImage(mock.Anything, imageId).
          Return(nil).
          Once()
      },
      args: args{
        ctx: context.Background(),
        input: &dto.ProfileAvatarUpdateDTO{
          Id:       types.MustCreateId(),
          Filename: fmt.Sprintf("%s.%s", gofakeit.AppName(), gofakeit.FileExtension()),
          Bytes:    []byte{0x12, 0x12, 0x12},
        },
      },
      want: status.FromRepository(sql.ErrNoRows, status.NullCode),
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mocked := newProfileMocked(t)
      if tt.setup != nil {
        tt.setup(&mocked, &tt.args, tt.want)
      }

      p := profileService{
        profileRepo: mocked.Profile,
        storageExt:  mocked.Storage,
        tracer:      mocked.Tracer,
      }

      if got := p.UpdateAvatar(tt.args.ctx, tt.args.input); !reflect.DeepEqual(got, tt.want) {
        t.Errorf("UpdateAvatar() = %v, want %v", got, tt.want)
      }
    })
  }
}
