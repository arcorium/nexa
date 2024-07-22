package mapper

import (
  "github.com/arcorium/nexa/proto/gen/go/common"
  relationv1 "github.com/arcorium/nexa/proto/gen/go/relation/v1"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  "nexa/services/relation/internal/domain/dto"
)

//func ToProtoBlock() relationv1.GetBlockedResponse{
//
//}

func ToProtoPagedElementOutput[T any](result *sharedDto.PagedElementResult[T]) *common.PagedElementOutput {
  return &common.PagedElementOutput{
    Element:       result.Element,
    Page:          result.Page,
    TotalElements: result.TotalElements,
    TotalPages:    result.TotalPages,
  }
}

func ToProtoBlockCount(dto *dto.BlockCountResponseDTO) *relationv1.BlockCount {
  return &relationv1.BlockCount{
    TotalBlocked: dto.Total,
  }
}
