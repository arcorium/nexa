package pg

import (
  "context"
  "github.com/stretchr/testify/require"
  "nexa/services/user/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util"
  "nexa/shared/wrapper"
  "reflect"
  "testing"
)

func Test_profileRepository_Create(t *testing.T) {
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
        profile: &entity.Profile{
          Id:        Users[PROFILE_SIZE].Id,
          FirstName: util.RandomString(8),
          LastName:  util.RandomString(12),
          Bio:       util.RandomString(20),
          PhotoURL:  types.FilePathFromString(util.RandomString(12)),
        },
      },
      wantErr: false,
    },
    {
      name: "Duplicate User Id",
      args: args{
        ctx: context.Background(),
        profile: &entity.Profile{
          Id:        Users[0].Id,
          FirstName: util.RandomString(8),
          LastName:  util.RandomString(12),
          Bio:       util.RandomString(20),
          PhotoURL:  types.FilePathFromString(util.RandomString(12)),
        },
      },
      wantErr: true,
    },
    {
      name: "User not found",
      args: args{
        ctx: context.Background(),
        profile: &entity.Profile{
          Id:        wrapper.DropError(types.NewId()),
          FirstName: util.RandomString(8),
          LastName:  util.RandomString(12),
          Bio:       util.RandomString(20),
          PhotoURL:  types.FilePathFromString(util.RandomString(12)),
        },
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tx, err := Db.BeginTx(tt.args.ctx, nil)
      require.NoError(t, err)
      defer tx.Rollback()

      p := profileRepository{
        db: tx,
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

func Test_profileRepository_FindByIds(t *testing.T) {
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
        userIds: []types.Id{Profiles[0].Id},
      },
      want:    Profiles[:1],
      wantErr: false,
    },
    {
      name: "Some",
      args: args{
        ctx:     context.Background(),
        userIds: []types.Id{Profiles[0].Id, wrapper.DropError(types.NewId()), Profiles[1].Id},
      },
      want:    []entity.Profile{Profiles[0], Profiles[1]},
      wantErr: false,
    },
    {
      name: "Profile Not Found",
      args: args{
        ctx:     context.Background(),
        userIds: []types.Id{wrapper.DropError(types.NewId())},
      },
      want:    nil,
      wantErr: true,
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tx, err := Db.BeginTx(tt.args.ctx, nil)
      require.NoError(t, err)
      defer tx.Rollback()

      p := profileRepository{
        db: tx,
      }
      got, err := p.FindByIds(tt.args.ctx, tt.args.userIds...)
      if (err != nil) != tt.wantErr {
        t.Errorf("FindByIds() error = %v, wantErr %v", err, tt.wantErr)
        return
      }
      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindByIds() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_profileRepository_Patch(t *testing.T) {
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
        ctx: nil,
        profile: &entity.Profile{
          Id:        Profiles[0].Id,
          FirstName: "arcorium",
          LastName:  "liz",
        },
      },
      wantErr: false,
    },
    {
      name: "Profile Not Found",
      args: args{
        ctx: nil,
        profile: &entity.Profile{
          Id:        Users[PROFILE_SIZE].Id,
          FirstName: "arcorium",
          LastName:  "liz",
        },
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tx, err := Db.BeginTx(tt.args.ctx, nil)
      require.NoError(t, err)
      defer tx.Rollback()

      p := profileRepository{
        db: tx,
      }
      err = p.Patch(tt.args.ctx, tt.args.profile)
      if (err != nil) != tt.wantErr {
        t.Errorf("Patch() error = %v, wantErr %v", err, tt.wantErr)
      }

      if err != nil {
        return
      }

      profiles, err := p.FindByIds(tt.args.ctx, []types.Id{tt.args.profile.Id}...)
      require.NoError(t, err)

      if profiles[0].FirstName == tt.args.profile.FirstName && profiles[0].LastName == tt.args.profile.LastName {
        return
      }

      t.Errorf("Patch() failed to update fields, got = %v, want = %v", profiles[0], *tt.args.profile)
    })
  }
}

func Test_profileRepository_Update(t *testing.T) {
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
        profile: &entity.Profile{
          Id:        Profiles[0].Id,
          FirstName: "arcorium",
          LastName:  "liz",
          Bio:       Profiles[0].Bio,
          PhotoURL:  Profiles[0].PhotoURL,
        },
      },
      wantErr: false,
    },
    {
      name: "Profile Not Found",
      args: args{
        ctx: context.Background(),
        profile: &entity.Profile{
          Id:        Users[1].Id,
          FirstName: "arcorium",
          LastName:  "liz",
          Bio:       Profiles[0].Bio,
          PhotoURL:  Profiles[0].PhotoURL,
        },
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tx, err := Db.BeginTx(tt.args.ctx, nil)
      require.NoError(t, err)
      defer tx.Rollback()

      p := profileRepository{
        db: tx,
      }
      err = p.Update(tt.args.ctx, tt.args.profile)
      if (err != nil) != tt.wantErr {
        t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
      }

      if err != nil {
        return
      }

      profiles, err := p.FindByIds(tt.args.ctx, []types.Id{tt.args.profile.Id}...)
      require.NoError(t, err)

      if !reflect.DeepEqual(profiles, *tt.args.profile) {
        t.Errorf("Update() got = %v, want %v", profiles, *tt.args.profile)
      }
    })
  }
}
