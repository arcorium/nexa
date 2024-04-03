package mapper

import (
	"nexa/services/authorization/shared/domain/entity"
	"nexa/services/authorization/shared/proto"
	"nexa/shared/types"
	"nexa/shared/util"
)

func ToRoleEntity(role *proto.RoleResponse) entity.Role {
	return entity.Role{
		Id:          types.IdFromString(role.Id),
		Name:        role.Name,
		Description: role.Description,
	}
}

func ToInternalCheckUserInput(resource string, actions ...string) []*proto.InternalCheckUserInput {
	return util.CastSlice2(actions, func(act string) *proto.InternalCheckUserInput {
		return &proto.InternalCheckUserInput{
			ResourceName: resource,
			ActionName:   act,
		}
	})
}
