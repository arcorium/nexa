package mapper

import (
	"nexa/services/authentication/internal/domain/dto"
	"nexa/services/authentication/shared/proto"
	sharedDto "nexa/shared/dto"
	sharedProto "nexa/shared/proto"
	"nexa/shared/util"
	"nexa/shared/wrapper"
)

func ToTokenRequestDTO(input *proto.RequestInput) dto.TokenRequestDTO {
	return dto.TokenRequestDTO{
		UsageId: input.UsageId,
	}
}

func ToTokenVerifyDTO(input *proto.VerifyInput) dto.TokenVerifyDTO {
	return dto.TokenVerifyDTO{
		Token:   input.Token,
		UsageId: input.UsageId,
	}
}

func ToTokenAddUsageDTO(input *proto.AddUsageInput) dto.TokenAddUsageDTO {
	return dto.TokenAddUsageDTO{
		Name:        input.Name,
		Description: wrapper.NewNullable(input.Description),
	}
}

func ToTokenUpdateUsageDTO(input *proto.UpdateUsageInput) dto.TokenUpdateUsageDTO {
	return dto.TokenUpdateUsageDTO{
		Id:          input.Id,
		Name:        wrapper.NewNullable(input.Name),
		Description: wrapper.NewNullable(input.Description),
	}
}

func ToTokenUsageResponse(resp *dto.TokenUsageResponseDTO) *proto.TokenUsageResponse {
	return &proto.TokenUsageResponse{
		Id:          resp.Id,
		Name:        resp.Name,
		Description: resp.Description,
	}
}

func ToTokenUsageResponses(result *sharedDto.PagedElementResult[dto.TokenUsageResponseDTO]) *proto.FindAllUsagesOutput {
	return &proto.FindAllUsagesOutput{
		Details: &sharedProto.PagedElementOutput{
			Element:       result.Element,
			Page:          result.Page,
			TotalElements: result.TotalElements,
			TotalPages:    result.TotalPages,
		},
		Usages: util.CastSlice(result.Data, ToTokenUsageResponse),
	}
}
