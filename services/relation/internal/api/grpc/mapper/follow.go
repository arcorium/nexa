package mapper

import (
  "github.com/arcorium/nexa/proto/gen/go/common"
  relationv1 "github.com/arcorium/nexa/proto/gen/go/relation/v1"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  "nexa/services/relation/internal/domain/dto"
  "nexa/services/relation/internal/domain/entity"
)

func ToPagedElementDTO(input *common.PagedElementInput) sharedDto.PagedElementDTO {
  if input == nil {
    return sharedDto.PagedElementDTO{
      Element: 0,
      Page:    0,
    }
  }
  return sharedDto.PagedElementDTO{
    Element: input.Element,
    Page:    input.Page,
  }
}

func ToProtoFollowStatus(status entity.FollowStatus) relationv1.Relation {
  switch status {
  case entity.FollowStatusNone:
    return relationv1.Relation_NONE
  case entity.FollowStatusFollower:
    return relationv1.Relation_FOLLOWER
  case entity.FollowStatusMutual:
    return relationv1.Relation_MUTUAL
  default:
    return relationv1.Relation_UNKNOWN
  }
}

func ToProtoFollowCount(count *dto.FollowCountResponseDTO) *relationv1.FollowCount {
  return &relationv1.FollowCount{
    TotalFollower: count.TotalFollowers,
    TotalFollowee: count.TotalFollowings,
  }
}
