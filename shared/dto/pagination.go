package dto

import (
  "github.com/arcorium/nexa/shared/util/repo"
  "math"
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
  page := p.Page - 1
  if p.Page == 0 {
    page = 0
  }
  return page * p.Element
}

func NewPagedElementResult[T any](data []T, currentPage, totalElements, totalPages uint64) PagedElementResult[T] {

  return PagedElementResult[T]{
    Data:          data,
    Element:       uint64(len(data)),
    Page:          currentPage,
    TotalElements: totalElements,
    TotalPages:    totalPages,
  }
}

func NewPagedElementResult2[T any](data []T, input *PagedElementDTO, totalElements uint64) PagedElementResult[T] {
  divider := input.Element
  if input.Element == 0 {
    divider = uint64(len(data))
  }
  currentPage := input.Page
  if input.Page == 0 {
    currentPage = 1
  }
  var totalPages uint64
  if totalElements == 0 {
    totalPages = 0
  } else {
    totalPages = uint64(math.Ceil(float64(totalElements) / float64(divider)))
  }
  return NewPagedElementResult(data, currentPage, totalElements, totalPages)
}

type PagedElementResult[T any] struct {
  Data []T

  Element       uint64
  Page          uint64
  TotalElements uint64
  TotalPages    uint64
}
