package mapper

import (
	"nexa/services/user/shared/domain/dto"
	"nexa/services/user/shared/proto"
	"nexa/shared/wrapper"
)

func ToDTOProfileUpdateInput(request *proto.UpdateProfileRequest) dto.ProfileUpdateInput {
	return dto.ProfileUpdateInput{
		UserId:    request.UserId,
		FirstName: wrapper.NewNullable(request.FirstName),
		LastName:  wrapper.NewNullable(request.LastName),
		Bio:       wrapper.NewNullable(request.Bio),
	}
}

func ToDTOProfilePictureUpdateInput(request *proto.UpdateProfileAvatarRequest) dto.ProfilePictureUpdateInput {
	return dto.ProfilePictureUpdateInput{
		UserId:   request.UserId,
		Filename: request.Filename,
		Bytes:    request.Chunk,
	}
}

func ToProtoProfile(response *dto.ProfileResponse) proto.Profile {
	return proto.Profile{
		UserId:    response.UserId,
		FirstName: response.FirstName,
		LastName:  response.LastName,
		Bio:       response.Bio,
		ImagePath: response.PhotoURL,
	}
}
