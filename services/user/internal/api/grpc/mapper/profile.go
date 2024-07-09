package mapper

import (
  "github.com/arcorium/nexa/proto/gen/go/user/v1"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/user/internal/domain/dto"
)

func ToProfileUpdateDTO(request *userv1.UpdateProfileRequest) (dto.ProfileUpdateDTO, error) {
  id, err := types.IdFromString(request.Id)
  if err != nil {
    return dto.ProfileUpdateDTO{}, sharedErr.NewFieldError("id", err).ToGrpcError()
  }

  return dto.ProfileUpdateDTO{
    Id:        id,
    FirstName: types.NewNullable(request.FirstName),
    LastName:  types.NewNullable(request.LastName),
    Bio:       types.NewNullable(request.Bio),
  }, nil
}

func ToProtoProfile(response *dto.ProfileResponseDTO) *userv1.Profile {
  return &userv1.Profile{
    Id:        response.Id.String(),
    FirstName: response.FirstName,
    LastName:  response.LastName,
    Bio:       response.Bio,
    ImagePath: response.PhotoURL.Path(),
  }
}
