package mapper

import (
  "nexa/services/user/internal/domain/dto"
  "nexa/services/user/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/wrapper"
)

//func MapProfileCreateInput(input *dto.ProfileCreateInput) entity.Profile {
//	return entity.Profile{
//		Id:        types.IdFromString(input.UserId),
//		FirstName: input.FirstName,
//		LastName:  input.LastName,
//		Bio:       input.Bio,
//	}
//}

func MapProfileUpdateDTO(input *dto.ProfileUpdateDTO) entity.Profile {
  profile := entity.Profile{
    Id: wrapper.DropError(types.IdFromString(input.UserId)),
  }

  if input.FirstName.HasValue() {
    profile.FirstName = input.FirstName.Value2()
  }

  wrapper.SetOnNonNull(&profile.LastName, input.LastName)
  wrapper.SetOnNonNull(&profile.Bio, input.Bio)
  return profile
}

func MapProfilePictureUpdateDTO(input *dto.ProfilePictureUpdateDTO) entity.Profile {
  return entity.Profile{
    Id: wrapper.DropError(types.IdFromString(input.UserId)),
  }
}

func ToProfileResponse(profile *entity.Profile) dto.ProfileResponseDTO {
  return dto.ProfileResponseDTO{
    //UserId:    profile.Id.Underlying().String(),
    FirstName: profile.FirstName,
    LastName:  profile.LastName,
    Bio:       profile.LastName,
    PhotoURL:  profile.PhotoURL.FullPath(""), // TODO: Implement it, take parameter for file service
  }
}
