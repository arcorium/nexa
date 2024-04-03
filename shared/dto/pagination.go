package dto

import (
	"math"
	"nexa/shared/util/repo"
)

type PagedElementDTO struct {
	Element uint64
	Page    uint64
}

func (p *PagedElementDTO) ToQueryParam() repo.QueryParameter {
	return repo.QueryParameter{
		Offset: p.Offset(),
		Limit:  p.Element,
	}
}

func (p *PagedElementDTO) Offset() uint64 {
	return (p.Page - 1) * p.Element
}

func NewPagedElementOutput[T any](data []T, currentPage, totalElements, totalPages uint64) PagedElementResult[T] {

	return PagedElementResult[T]{
		Data:          data,
		Element:       uint64(len(data)),
		Page:          currentPage,
		TotalElements: totalElements,
		TotalPages:    totalPages,
	}
}

func NewPagedElementOutput2[T any](data []T, input *PagedElementDTO, totalElements uint64) PagedElementResult[T] {
	return NewPagedElementOutput(data, input.Page, totalElements, uint64(math.Ceil(float64(totalElements)/float64(input.Element))))
}

type PagedElementResult[T any] struct {
	Data []T

	Element       uint64
	Page          uint64
	TotalElements uint64
	TotalPages    uint64
}
