package mapper

import (
  authNv1 "github.com/arcorium/nexa/proto/gen/go/authentication/v1"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  "nexa/services/authentication/internal/domain/dto"
)

func ToProfileUpdateDTO(request *authNv1.UpdateProfileRequest) (dto.ProfileUpdateDTO, error) {
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

func ToProtoProfile(response *dto.ProfileResponseDTO) *authNv1.Profile {
  return &authNv1.Profile{
    Id:        response.Id.String(),
    FirstName: response.FirstName,
    LastName:  response.LastName,
    Bio:       response.Bio,
    ImagePath: response.PhotoURL.Path(),
  }
}
