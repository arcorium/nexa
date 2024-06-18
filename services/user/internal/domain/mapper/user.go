package mapper

import (
  "nexa/services/user/internal/domain/dto"
  "nexa/services/user/internal/domain/entity"
  sharedErr "nexa/shared/errors"
  "nexa/shared/status"
  "nexa/shared/types"
  "nexa/shared/util"
  "nexa/shared/wrapper"
  "time"
)

func MapUserCreateDTO(input *dto.UserCreateDTO) (entity.User, entity.Profile, status.Object) {
  // Id
  id, err := types.NewId()
  if err != nil {
    return entity.User{}, entity.Profile{}, status.ErrInternal(sharedErr.ErrIdCreation)
  }
  // Password
  pass, err := types.PasswordFromString(input.Password)
  if err != nil {
    return entity.User{}, entity.Profile{}, status.ErrInternal(err)
  }

  user := entity.User{
    Id:         id,
    Username:   input.Username,
    Email:      wrapper.DropError(types.EmailFromString(input.Email)),
    Password:   pass, // hashed
    IsVerified: false,
    IsDeleted:  false,
  }

  profile := entity.Profile{
    Id:        user.Id,
    FirstName: input.FirstName,
  }

  wrapper.SetOnNonNull(&profile.LastName, input.LastName)
  wrapper.SetOnNonNull(&profile.Bio, input.Bio)

  return user, profile, status.SuccessInternal()
}

func MapUserUpdateDTO(input *dto.UserUpdateDTO) (entity.User, status.Object) {
  user := entity.User{
    Id:         wrapper.DropError(types.IdFromString(input.Id)),
    IsVerified: false,
    IsDeleted:  false,
  }

  wrapper.SetOnNonNull(&user.Username, input.Username)
  if input.Email.HasValue() {
    email, err := types.EmailFromString(input.Email.RawValue())
    if err != nil {
      return entity.User{}, status.ErrBadRequest(err)
    }
    user.Email = email
  }
  return user, status.SuccessInternal()
}

func MapUserUpdatePasswordDTO(input *dto.UserUpdatePasswordDTO) (entity.User, status.Object) {
  pass, err := types.PasswordFromString(input.NewPassword)
  if err != nil {
    return entity.User{}, status.ErrInternal(err)
  }

  return entity.User{
    Id:       wrapper.DropError(types.IdFromString(input.Id)),
    Password: pass,
  }, status.SuccessInternal()
}

func MapUserResetPasswordDTO(input *dto.UserResetPasswordDTO) (entity.User, status.Object) {
  pass, err := types.PasswordFromString(input.NewPassword)
  if err != nil {
    return entity.User{}, status.ErrInternal(err)
  }

  return entity.User{
    Id:       wrapper.DropError(types.IdFromString(input.Id)),
    Password: pass,
  }, status.SuccessInternal()
}

func MapUserBannedDTO(input *dto.UserBannedDTO) entity.User {
  return entity.User{
    Id:          wrapper.DropError(types.IdFromString(input.Id)),
    BannedUntil: time.Now().Add(input.Duration),
  }
}

func ToUserResponse(user *entity.User) dto.UserResponseDTO {
  return dto.UserResponseDTO{
    Id:         user.Id.Underlying().String(),
    Username:   user.Username,
    Email:      user.Email.Underlying(),
    IsVerified: user.IsVerified,
    Profile: util.NilOr[dto.ProfileResponseDTO](user.Profile, func(obj *entity.Profile) *dto.ProfileResponseDTO {
      tmp := ToProfileResponse(obj)
      return &tmp
    }),
  }
}
