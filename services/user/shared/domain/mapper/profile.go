package mapper

import (
	"nexa/services/user/shared/domain/dto"
	"nexa/services/user/shared/domain/entity"
	"nexa/shared/types"
)

//func MapProfileCreateInput(input *dto.ProfileCreateInput) entity.Profile {
//	return entity.Profile{
//		Id:        types.IdFromString(input.UserId),
//		FirstName: input.FirstName,
//		LastName:  input.LastName,
//		Bio:       input.Bio,
//	}
//}

func MapProfileUpdateInput(input *dto.ProfileUpdateInput) entity.Profile {
	return entity.Profile{
		Id:        types.IdFromString(input.UserId),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Bio:       input.Bio,
	}
}

func MapProfilePictureUpdateInput(input *dto.ProfilePictureUpdateInput) entity.Profile {
	return entity.Profile{
		Id: types.IdFromString(input.UserId),
	}
}

func ToProfileResponse(profile *entity.Profile) dto.ProfileResponse {
	return dto.ProfileResponse{
		UserId:    profile.Id.Underlying().String(),
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Bio:       profile.LastName,
		PhotoURL:  profile.PhotoURL.FullPath(""), // TODO: Implement it, take parameter for file service
	}
}
