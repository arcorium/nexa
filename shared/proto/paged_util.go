package proto

import (
	"nexa/shared/dto"
)

func (x *PagedElementInput) ToDTO() dto.PagedElementDTO {
	return dto.PagedElementDTO{
		Element: x.Element,
		Page:    x.Page,
	}
}
