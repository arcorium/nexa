package mapper

import (
	"nexa/services/authorization/internal/domain/dto"
	"nexa/services/authorization/shared/proto"
	sharedDto "nexa/shared/dto"
	sharedProto "nexa/shared/proto"
	"nexa/shared/util"
	"nexa/shared/wrapper"
)

func ToResourceCreateDTO(input *proto.ResourceCreateInput) dto.ResourceCreateDTO {
	return dto.ResourceCreateDTO{
		Name:        input.Name,
		Description: wrapper.NewNullable(input.Description),
	}
}

func ToResourceUpdateDTO(input *proto.ResourceUpdateInput) dto.ResourceUpdateDTO {
	return dto.ResourceUpdateDTO{
		Id:          input.Id,
		Name:        wrapper.NewNullable(input.Name),
		Description: wrapper.NewNullable(input.Description),
	}
}

func ToResourceResponse(resp *dto.ResourceResponseDTO) *proto.ResourceResponse {
	return &proto.ResourceResponse{
		Id:          resp.Id,
		Name:        resp.Name,
		Description: resp.Description,
	}
}

func ToResourceResponses(result *sharedDto.PagedElementResult[dto.ResourceResponseDTO]) *proto.ResourceFindAllResponse {
	return &proto.ResourceFindAllResponse{
		Details: &sharedProto.PagedElementOutput{
			Element:       result.Element,
			Page:          result.Page,
			TotalElements: result.TotalElements,
			TotalPages:    result.TotalPages,
		},
		Resources: util.CastSlice(result.Data, func(from *dto.ResourceResponseDTO) *proto.ResourceResponse {
			return ToResourceResponse(from)
		}),
	}
}
