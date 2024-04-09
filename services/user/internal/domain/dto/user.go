package dto

import (
  "nexa/shared/wrapper"
  "time"
)

type UserResponseDTO struct {
  Id         string              `json:"id"`
  Username   string              `json:"username"`
  Email      string              `json:"email"`
  IsVerified bool                `json:"is_verified"`
  Profile    *ProfileResponseDTO `json:"profile"`
}

type UserCreateDTO struct {
  Username  string `validate:"required,gte=6"`
  Email     string `validate:"required,email"`
  Password  string `validate:"required,gte=6"`
  FirstName string `validate:"required"`
  LastName  wrapper.Nullable[string]
  Bio       wrapper.Nullable[string]
}

type UserUpdateDTO struct {
  Id       string `validate:"required,uuid4"`
  Username wrapper.Nullable[string]
  Email    wrapper.Nullable[string] `validate:"email"`
}

type UserUpdatePasswordDTO struct {
  Id           string `validate:"required,uuid4"`
  LastPassword string `validate:"required"`
  NewPassword  string `validate:"required"`
}

type UserBannedDTO struct {
  Id       string        `validate:"required,uuid4"`
  Duration time.Duration `validate:"required"`
}

type UserResetPasswordDTO struct {
  Id          string `validate:"required,uuid4"`
  NewPassword string `validate:"required"`
}
