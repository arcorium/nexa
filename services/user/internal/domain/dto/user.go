package dto

import (
  "nexa/services/user/internal/domain/entity"
  "nexa/shared/types"
  "time"
)

type UserResponseDTO struct {
  Id         types.Id
  Username   string
  Email      types.Email
  IsVerified bool
  Profile    *ProfileResponseDTO
}

type UserCreateDTO struct {
  Username  string `validate:"required,gte=6"`
  Email     types.Email
  Password  types.Password
  FirstName string `validate:"required"`
  LastName  types.NullableString
  Bio       types.NullableString
}

func (d *UserCreateDTO) ToDomain() (*entity.User, *entity.Profile, error) {
  userId, err := types.NewId()
  if err != nil {
    return nil, nil, err
  }

  profileId, err := types.NewId()
  if err != nil {
    return nil, nil, err
  }

  password, err := d.Password.Hash()
  if err != nil {
    return nil, nil, err
  }

  user := &entity.User{
    Id:        userId,
    Username:  d.Username,
    Email:     d.Email,
    Password:  password, // hashed
    CreatedAt: time.Now(),
  }

  profile := &entity.Profile{
    Id:        profileId,
    UserId:    user.Id,
    FirstName: d.FirstName,
  }

  types.SetOnNonNull(&profile.LastName, d.LastName)
  types.SetOnNonNull(&profile.Bio, d.Bio)

  return user, profile, nil
}

type UserUpdateDTO struct {
  Id       types.Id
  Username types.NullableString
  Email    types.NullableEmail
}

func (d *UserUpdateDTO) ToDomain() entity.PatchedUser {
  user := entity.PatchedUser{
    Id: d.Id,
  }

  types.SetOnNonNull(&user.Username, d.Username)
  types.SetOnNonNull(&user.Email, d.Email)

  return user
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
  Token       types.NullableString
  UserId      types.NullableId
  LogoutAll   bool           // Need to logout all devices?
  NewPassword types.Password `validate:"required,gt=6"`
}
