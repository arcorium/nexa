package mapper

import (
  "nexa/proto/gen/go/user/v1"
  "nexa/services/user/internal/domain/dto"
  sharedErr "nexa/shared/errors"
  "nexa/shared/types"
  "nexa/shared/wrapper"
)

func ToProfileUpdateDTO(request *userv1.UpdateProfileRequest) (dto.ProfileUpdateDTO, error) {
  id, err := types.IdFromString(request.UserId)
  if err != nil {
    return dto.ProfileUpdateDTO{},
      sharedErr.GrpcFieldErrors2(sharedErr.NewFieldError("user_id", err))
  }

  return dto.ProfileUpdateDTO{
    UserId:    id,
    FirstName: wrapper.NewNullable(request.FirstName),
    LastName:  wrapper.NewNullable(request.LastName),
    Bio:       wrapper.NewNullable(request.Bio),
  }, nil
}

func ToProtoProfile(response *dto.ProfileResponseDTO) *userv1.Profile {
  return &userv1.Profile{
    //UserId:    response.UserId,
    FirstName: response.FirstName,
    LastName:  response.LastName,
    Bio:       response.Bio,
    ImagePath: response.PhotoURL,
  }
}
