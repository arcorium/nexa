package dto

import (
  "nexa/services/user/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/wrapper"
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
  LastName  wrapper.NullableString
  Bio       wrapper.NullableString
}

func (d *UserCreateDTO) ToDomain() (*entity.User, *entity.Profile, error) {
  id, err := types.NewId()
  if err != nil {
    return nil, nil, err
  }

  password, err := d.Password.Hash()
  if err != nil {
    return nil, nil, err
  }

  user := &entity.User{
    Id:        id,
    Username:  d.Username,
    Email:     d.Email,
    Password:  password, // hashed
    IsDeleted: false,
  }

  profile := &entity.Profile{
    Id:        user.Id,
    FirstName: d.FirstName,
  }

  wrapper.SetOnNonNull(&profile.LastName, d.LastName)
  wrapper.SetOnNonNull(&profile.Bio, d.Bio)

  return user, profile, nil
}

type UserUpdateDTO struct {
  Id       types.Id
  Username wrapper.NullableString
  Email    wrapper.Nullable[types.Email]
}

func (d *UserUpdateDTO) ToDomain() entity.User {
  user := entity.User{
    Id:        d.Id,
    IsDeleted: false,
  }

  wrapper.SetOnNonNull(&user.Username, d.Username)

  return user
}

type UserUpdatePasswordDTO struct {
  Id           types.Id
  LastPassword types.Password `validate:"required"`
  NewPassword  types.Password `validate:"required,gte=6"`
}

func (p *UserUpdatePasswordDTO) ToDomain() (entity.User, error) {
  hashedPassword, err := p.NewPassword.Hash()
  if err != nil {
    return entity.User{}, err
  }

  return entity.User{
    Id:       p.Id,
    Password: hashedPassword,
  }, nil
}

type UserBannedDTO struct {
  Id       types.Id
  Duration time.Duration `validate:"required,gte=1"`
}

func (u *UserBannedDTO) ToDomain() entity.User {
  return entity.User{
    Id:          u.Id,
    IsDeleted:   true,
    BannedUntil: time.Now().Add(u.Duration),
  }
}

type UserResetPasswordDTO struct {
  Token       wrapper.NullableString
  LogoutAll   bool           // Need to logout all devices?
  NewPassword types.Password `validate:"required,gt=6"`
}
