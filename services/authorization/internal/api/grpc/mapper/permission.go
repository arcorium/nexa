package mapper

import (
	"nexa/services/authorization/internal/domain/dto"
	"nexa/services/authorization/shared/proto"
	sharedDto "nexa/shared/dto"
	sharedProto "nexa/shared/proto"
	"nexa/shared/util"
)

func ToPermissionCreateDTO(input *proto.PermissionCreateInput) dto.PermissionCreateDTO {
	return dto.PermissionCreateDTO{
		ResourceId: input.ResourceId,
		ActionId:   input.ActionId,
	}
}

func ToCheckUserPermissionDTO(input *proto.CheckUserInput) dto.CheckUserPermissionDTO {
	return dto.CheckUserPermissionDTO{
		UserId: input.UserId,
		Permissions: util.CastSlice2(input.Permissions, func(from *proto.InternalCheckUserInput) dto.InternalCheckUserPermissionDTO {
			return dto.InternalCheckUserPermissionDTO{
				Resource: from.ResourceName,
				Action:   from.ActionName,
			}
		}),
	}
}

func ToPermissionResponse(resp *dto.PermissionResponseDTO) *proto.PermissionResponse {
	return &proto.PermissionResponse{
		ResourceId:          resp.Resource.Id,
		ResourceName:        resp.Resource.Name,
		ResourceDescription: resp.Resource.Description,
		ActionId:            resp.Action.Id,
		ActionName:          resp.Action.Name,
		ActionDescription:   resp.Action.Description,
		Code:                resp.Code,
	}
}

func ToPermissionResponses(result *sharedDto.PagedElementResult[dto.PermissionResponseDTO]) *proto.PermissionFindAllOutput {
	return &proto.PermissionFindAllOutput{
		Details: &sharedProto.PagedElementOutput{
			Element:       result.Element,
			Page:          result.Page,
			TotalElements: result.TotalElements,
			TotalPages:    result.TotalPages,
		},
		Permissions: util.CastSlice(result.Data, ToPermissionResponse),
	}
}
