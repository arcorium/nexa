package pg

import (
  "context"
  "github.com/google/uuid"
  "github.com/stretchr/testify/require"
  "nexa/services/user/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util"
  "nexa/shared/util/repo"
  "nexa/shared/wrapper"
  "reflect"
  "testing"
)

func ignoreUserFields(t *testing.T, got []entity.User, want []entity.User) {
  // Ignore time fields
  require.Len(t, got, len(want))

  for i := 0; i < len(want); i += 1 {
    got[i].BannedUntil = want[i].BannedUntil
    got[i].Password = want[i].Password // Password can have different hash even when the plain is the same
  }
}

func Test_userRepository_Create(t *testing.T) {
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
        user: &entity.User{
          Id:         types.Id(uuid.New()),
          Username:   util.RandomString(15),
          Email:      wrapper.DropError(types.EmailFromString("something@gmail.com")),
          Password:   wrapper.DropError(types.PasswordFromString(util.RandomString(12))),
          IsVerified: false,
          IsDeleted:  false,
        },
      },
      wantErr: false,
    },
    {
      name: "Duplicate Username",
      args: args{
        ctx: context.Background(),
        user: &entity.User{
          Id:         types.Id(uuid.New()),
          Username:   Users[0].Username,
          Email:      wrapper.DropError(types.EmailFromString(util.RandomString(12))),
          Password:   wrapper.DropError(types.PasswordFromString(util.RandomString(12))),
          IsVerified: false,
          IsDeleted:  false,
        },
      },
      wantErr: true,
    },
    {
      name: "Duplicate Email",
      args: args{
        ctx: context.Background(),
        user: &entity.User{
          Id:         types.Id(uuid.New()),
          Username:   util.RandomString(15),
          Email:      Users[0].Email,
          Password:   wrapper.DropError(types.PasswordFromString(util.RandomString(12))),
          IsVerified: false,
          IsDeleted:  false,
        },
      },
      wantErr: true,
    },
    {
      name: "Bad Id",
      args: args{
        ctx: context.Background(),
        user: &entity.User{
          Id:         wrapper.DropError(types.IdFromString(util.RandomString(24))),
          Username:   util.RandomString(15),
          Email:      wrapper.DropError(types.EmailFromString(util.RandomString(12))),
          Password:   wrapper.DropError(types.PasswordFromString(util.RandomString(12))),
          IsVerified: false,
          IsDeleted:  false,
        },
      },
      wantErr: true,
    },
    {
      name: "Null Id",
      args: args{
        ctx: context.Background(),
        user: &entity.User{
          Id:         types.NullId(),
          Username:   util.RandomString(15),
          Email:      wrapper.DropError(types.EmailFromString(util.RandomString(12))),
          Password:   wrapper.DropError(types.PasswordFromString(util.RandomString(12))),
          IsVerified: false,
          IsDeleted:  false,
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

      u := userRepository{
        db: tx,
      }
      if err := u.Create(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
        t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
      }
    })
  }
}

func Test_userRepository_Delete(t *testing.T) {
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
        id:  Users[0].Id,
      },
      wantErr: false,
    },
    {
      name: "User not found",
      args: args{
        ctx: context.Background(),
        id:  wrapper.DropError(types.NewId()),
      },
      wantErr: true,
    },
    {
      name: "Null id",
      args: args{
        ctx: context.Background(),
        id:  types.NullId(),
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tx, err := Db.BeginTx(tt.args.ctx, nil)
      require.NoError(t, err)
      defer tx.Rollback()

      u := userRepository{
        db: tx,
      }

      err = u.Delete(tt.args.ctx, tt.args.id)

      if (err != nil) != tt.wantErr {
        t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
      }

      if err != nil {
        return
      }

      users, err := u.FindByIds(tt.args.ctx, []types.Id{tt.args.id}...)
      require.NoError(t, err)
      if !users[0].IsDeleted {
        t.Errorf("Delete() failed to change is_delete field, obj = %v", users[0])
      }
    })
  }
}

func Test_userRepository_FindAllUsers(t *testing.T) {
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
        Data:    Users,
        Total:   USER_SIZE,
        Element: USER_SIZE,
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
        Data:    Users[1:3],
        Total:   USER_SIZE,
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
        Data:    Users[2:],
        Total:   USER_SIZE,
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
        Total:   USER_SIZE,
        Element: 0,
      },
      wantErr: true,
    },
  }
  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tx, err := Db.BeginTx(tt.args.ctx, nil)
      require.NoError(t, err)
      defer tx.Rollback()

      u := userRepository{
        db: tx,
      }

      got, err := u.FindAllUsers(tt.args.ctx, tt.args.query)
      if (err != nil) != tt.wantErr {
        t.Errorf("FindAllUsers() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      ignoreUserFields(t, got.Data, tt.want.Data)

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindAllUsers() \ngot = %v\nwant %v", got, tt.want)
      }
    })
  }
}

func Test_userRepository_FindByEmails(t *testing.T) {
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
          Users[0].Email,
        },
      },
      want:    Users[:1],
      wantErr: false,
    },
    {
      name: "Multiple Emails",
      args: args{
        ctx: context.Background(),
        emails: []types.Email{
          Users[0].Email,
          Users[1].Email,
        },
      },
      want:    Users[0:2],
      wantErr: false,
    },
    {
      name: "User not found",
      args: args{
        ctx: context.Background(),
        emails: []types.Email{
          wrapper.DropError(types.EmailFromString(util.RandomString(12) + "@gmail.com")),
        },
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

      u := userRepository{
        db: tx,
      }
      got, err := u.FindByEmails(tt.args.ctx, tt.args.emails...)
      if (err != nil) != tt.wantErr {
        t.Errorf("FindByEmails() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      ignoreUserFields(t, got, tt.want)

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindByEmails() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_userRepository_FindByIds(t *testing.T) {
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
          Users[0].Id,
        },
      },
      want:    Users[:1],
      wantErr: false,
    },
    {
      name: "Multiple",
      args: args{
        ctx: context.Background(),
        userIds: []types.Id{
          Users[0].Id,
          Users[1].Id,
        },
      },
      want:    Users[:2],
      wantErr: false,
    },
    {
      name: "User Not Found",
      args: args{
        ctx: context.Background(),
        userIds: []types.Id{
          types.Id(wrapper.DropError(uuid.NewUUID())),
        },
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

      u := userRepository{
        db: tx,
      }
      got, err := u.FindByIds(tt.args.ctx, tt.args.userIds...)
      if (err != nil) != tt.wantErr {
        t.Errorf("FindByIds() error = %v, wantErr %v", err, tt.wantErr)
        return
      }

      ignoreUserFields(t, got, tt.want)

      if !reflect.DeepEqual(got, tt.want) {
        t.Errorf("FindByIds() got = %v, want %v", got, tt.want)
      }
    })
  }
}

func Test_userRepository_Patch(t *testing.T) {
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
        user: &entity.User{
          Id:       Users[0].Id,
          Username: "arcorium",
          Email:    wrapper.DropError(types.EmailFromString("something@gmail.com")),
        },
      },
      wantErr: false,
    },
    {
      name: "User not found",
      args: args{
        ctx: context.Background(),
        user: &entity.User{
          Id:       wrapper.DropError(types.NewId()),
          Username: "arcorium",
        },
      },
      wantErr: true,
    },
    {
      name: "Bad Id",
      args: args{
        ctx: context.Background(),
        user: &entity.User{
          Id:       types.NullId(),
          Username: "arcorium",
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

      u := userRepository{
        db: tx,
      }
      err = u.Patch(tt.args.ctx, tt.args.user)

      if (err != nil) != tt.wantErr {
        t.Errorf("Patch() error = %v, wantErr %v", err, tt.wantErr)
      }

      if err != nil {
        return
      }

      users, err := u.FindByIds(tt.args.ctx, tt.args.user.Id)
      require.NoError(t, err)

      ignoreUserFields(t, users, []entity.User{*tt.args.user})

      if !reflect.DeepEqual(users[0], *tt.args.user) {
        t.Errorf("Patch() \ngot = %v, \nwant = %v", users[0], tt.args.user)
      }
    })
  }
}

func Test_userRepository_Update(t *testing.T) {
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
        user: &entity.User{
          Id:          Users[0].Id,
          Username:    "arcorium",
          Email:       "arcorium@gmail.com",
          Password:    Users[0].Password,
          IsVerified:  Users[0].IsVerified,
          IsDeleted:   Users[0].IsDeleted,
          BannedUntil: Users[0].BannedUntil,
        },
      },
      wantErr: false,
    },
    {
      name: "Non-Important field is empty",
      args: args{
        ctx: context.Background(),
        user: &entity.User{
          Id:       Users[0].Id,
          Username: "arcorium",
          Email:    "arcorium@gmail.com",
          Password: Users[0].Password,
          //IsVerified:  Users[0].IsVerified,
          IsDeleted:   Users[0].IsDeleted,
          BannedUntil: Users[0].BannedUntil,
        },
      },
      wantErr: false,
    },
    {
      name: "Important field is empty",
      args: args{
        ctx: context.Background(),
        user: &entity.User{
          Id: Users[0].Id,
          //Username: "arcorium",
          //Email:    "arcorium@gmail.com",
          Password:    Users[0].Password,
          IsVerified:  Users[0].IsVerified,
          IsDeleted:   Users[0].IsDeleted,
          BannedUntil: Users[0].BannedUntil,
        },
      },
      wantErr: true,
    },
    {
      name: "User not found",
      args: args{
        ctx: context.Background(),
        user: &entity.User{
          Id:          wrapper.DropError(types.NewId()),
          Username:    "arcorium",
          Email:       "arcorium@gmail.com",
          Password:    Users[0].Password,
          IsVerified:  Users[0].IsVerified,
          IsDeleted:   Users[0].IsDeleted,
          BannedUntil: Users[0].BannedUntil,
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

      u := userRepository{
        db: tx,
      }

      err = u.Update(tt.args.ctx, tt.args.user)

      if (err != nil) != tt.wantErr {
        t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
      }

      if err != nil {
        return
      }

      users, err := u.FindByIds(tt.args.ctx, []types.Id{tt.args.user.Id}...)
      require.NoError(t, err)

      ignoreUserFields(t, users, []entity.User{*tt.args.user})

      if !reflect.DeepEqual(users[0], *tt.args.user) {
        t.Errorf("FindByIds() got = %v, want %v", users[0], *tt.args.user)
      }
    })
  }
}
