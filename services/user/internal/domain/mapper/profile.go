package mapper

import (
  "nexa/services/user/internal/domain/dto"
  "nexa/services/user/internal/domain/entity"
)

func ToProfileResponse(profile *entity.Profile) dto.ProfileResponseDTO {
  return dto.ProfileResponseDTO{
    FirstName: profile.FirstName,
    LastName:  profile.LastName,
    Bio:       profile.LastName,
    PhotoURL:  profile.PhotoURL.Path(),
  }
}
