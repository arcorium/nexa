package mapper

import (
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/entity"
)

func ToUserResponse(user *entity.User) dto.UserResponseDTO {
  return dto.UserResponseDTO{
    Id:         user.Id,
    Username:   user.Username,
    Email:      user.Email,
    IsVerified: user.IsVerified,
    Profile: types.NilOrElse[dto.ProfileResponseDTO](user.Profile, func(obj *entity.Profile) *dto.ProfileResponseDTO {
      tmp := ToProfileResponse(obj)
      return &tmp
    }),
  }
}
