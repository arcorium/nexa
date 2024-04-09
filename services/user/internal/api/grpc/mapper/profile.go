package mapper

import (
  proto "nexa/proto/generated/golang/user/v1"
  "nexa/services/user/internal/domain/dto"
  "nexa/shared/wrapper"
)

func ToDTOProfileUpdateInput(request *proto.UpdateProfileRequest) dto.ProfileUpdateDTO {
  return dto.ProfileUpdateDTO{
    UserId:    request.UserId,
    FirstName: wrapper.NewNullable(request.FirstName),
    LastName:  wrapper.NewNullable(request.LastName),
    Bio:       wrapper.NewNullable(request.Bio),
  }
}

func ToDTOProfilePictureUpdateInput(request *proto.UpdateProfileAvatarRequest) dto.ProfilePictureUpdateDTO {
  return dto.ProfilePictureUpdateDTO{
    UserId:   request.UserId,
    Filename: request.Filename,
    Bytes:    request.Chunk,
  }
}

func ToProtoProfile(response *dto.ProfileResponseDTO) *proto.Profile {
  return &proto.Profile{
    //UserId:    response.UserId,
    FirstName: response.FirstName,
    LastName:  response.LastName,
    Bio:       response.Bio,
    ImagePath: response.PhotoURL,
  }
}
