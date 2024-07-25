package dto

import (
  "github.com/arcorium/nexa/shared/types"
  entity "nexa/services/authentication/internal/domain/entity"
  "time"
)

type UserResponseDTO struct {
  Id         types.Id
  Username   string
  Email      types.Email
  IsVerified bool
  IsBanned   bool
  Profile    *ProfileResponseDTO
}

type UserCreateDTO struct {
  Username  string `validate:"required,gte=6"`
  Email     types.Email
  Password  types.Password
  FirstName string `validate:"required"`
  LastName  types.NullableString
  Bio       types.NullableString
  RoleIds   []types.Id
}

func (d *UserCreateDTO) ToDomain() (entity.User, entity.Profile, error) {
  user, err := entity.NewUser(d.Username, d.Email, d.Password)
  if err != nil {
    return entity.User{}, entity.Profile{}, err
  }

  profile, err := entity.NewProfile(user.Id, d.FirstName)
  if err != nil {
    return entity.User{}, entity.Profile{}, err
  }

  types.SetOnNonNull(&profile.LastName, d.LastName)
  types.SetOnNonNull(&profile.Bio, d.Bio)

  return user, profile, nil
}

type UserUpdateDTO struct {
  Id        types.Id
  Username  types.NullableString
  Email     types.NullableEmail
  FirstName types.NullableString
  LastName  types.NullableString
  Bio       types.NullableString
}

func (d *UserUpdateDTO) ToDomain() (entity.PatchedUser, entity.PatchedProfile) {
  user := entity.PatchedUser{
    Id: d.Id,
  }

  types.SetOnNonNull(&user.Username, d.Username)
  types.SetOnNonNull(&user.Email, d.Email)

  profile := entity.PatchedProfile{
    //Id:       d.Id,
    UserId:   d.Id,
    LastName: d.LastName,
    Bio:      d.Bio,
  }

  types.SetOnNonNull(&profile.FirstName, d.FirstName)
  return user, profile
}

type UserUpdatePasswordDTO struct {
  Id           types.Id
  LastPassword types.Password `validate:"required"`
  NewPassword  types.Password `validate:"required,gte=6"`
}

func (p *UserUpdatePasswordDTO) ToDomain() (entity.PatchedUser, error) {
  hashedPassword, err := p.NewPassword.Hash()
  if err != nil {
    return entity.PatchedUser{}, err
  }

  return entity.PatchedUser{
    Id:       p.Id,
    Password: hashedPassword,
  }, nil
}

type UserBannedDTO struct {
  Id       types.Id
  Duration time.Duration
}

func (u *UserBannedDTO) ToDomain() entity.PatchedUser {
  return entity.PatchedUser{
    Id:             u.Id,
    BannedDuration: types.SomeNullable(u.Duration),
  }
}

type ResetUserPasswordDTO struct {
  UserId      types.Id
  LogoutAll   bool
  NewPassword types.Password `validate:"required,gt=6"`
}

type ResetPasswordWithTokenDTO struct {
  Token       string
  LogoutAll   bool
  NewPassword types.Password `validate:"required,gt=6"`
}

type UpdateUserAvatarDTO struct {
  UserId   types.Id
  Filename string `validate:"required"`
  Bytes    []byte `validate:"required"`
}
