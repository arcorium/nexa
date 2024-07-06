package mapper

import (
  "nexa/services/user/internal/domain/dto"
  "nexa/services/user/internal/domain/entity"
)

func ToProfileResponse(profile *entity.Profile) dto.ProfileResponseDTO {
  return dto.ProfileResponseDTO{
    Id:        profile.Id,
    FirstName: profile.FirstName,
    LastName:  profile.LastName,
    Bio:       profile.Bio,
    PhotoURL:  profile.PhotoURL,
  }
}
