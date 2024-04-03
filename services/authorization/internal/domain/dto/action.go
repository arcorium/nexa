package dto

import (
	"nexa/services/authorization/shared/domain/entity"
	"nexa/shared/types"
	"nexa/shared/wrapper"
)

type ActionCreateDTO struct {
	Name        string `validate:"required"`
	Description wrapper.NullableString
}

func (a *ActionCreateDTO) ToDomain() entity.Action {
	action := entity.Action{
		Id:   types.NewId(),
		Name: a.Name,
	}

	wrapper.SetOnNonNull(&action.Description, a.Description)
	return action
}

type ActionUpdateDTO struct {
	Id          string `validate:"required,uuid4"`
	Name        wrapper.NullableString
	Description wrapper.NullableString
}

func (u *ActionUpdateDTO) ToDomain() entity.Action {
	action := entity.Action{
		Id: types.IdFromString(u.Id),
	}

	wrapper.SetOnNonNull(&action.Name, u.Name)
	wrapper.SetOnNonNull(&action.Description, u.Description)
	return action
}

type ActionResponseDTO struct {
	Id          string
	Name        string
	Description string
}
