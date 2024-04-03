package mapper

import (
	"nexa/services/authorization/internal/domain/dto"
	"nexa/services/authorization/shared/proto"
	sharedDto "nexa/shared/dto"
	sharedProto "nexa/shared/proto"
	"nexa/shared/util"
	"nexa/shared/wrapper"
)

func ToActionCreateDTO(input *proto.ActionCreateInput) dto.ActionCreateDTO {
	return dto.ActionCreateDTO{
		Name:        input.Name,
		Description: wrapper.NewNullable(input.Description),
	}
}

func ToActionUpdateDTO(input *proto.ActionUpdateInput) dto.ActionUpdateDTO {
	return dto.ActionUpdateDTO{
		Id:          input.Id,
		Name:        wrapper.NewNullable(input.Name),
		Description: wrapper.NewNullable(input.Description),
	}
}

func ToActionResponse(responseDTO *dto.ActionResponseDTO) *proto.ActionResponse {
	return &proto.ActionResponse{
		Id:          responseDTO.Id,
		Name:        responseDTO.Name,
		Description: responseDTO.Description,
	}
}

func ToActionResponses(result *sharedDto.PagedElementResult[dto.ActionResponseDTO]) *proto.ActionFindAllResponse {
	return &proto.ActionFindAllResponse{
		Details: &sharedProto.PagedElementOutput{
			Element:       result.Element,
			Page:          result.Page,
			TotalElements: result.TotalElements,
			TotalPages:    result.TotalPages,
		},
		Actions: util.CastSlice(result.Data, func(from *dto.ActionResponseDTO) *proto.ActionResponse {
			return ToActionResponse(from)
		}),
	}
}
