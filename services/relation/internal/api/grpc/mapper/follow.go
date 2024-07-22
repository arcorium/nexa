package mapper

import (
  relationv1 "github.com/arcorium/nexa/proto/gen/go/relation/v1"
  "nexa/services/relation/internal/domain/dto"
  "nexa/services/relation/internal/domain/entity"
)

func ToProtoFollowStatus(status entity.FollowStatus) relationv1.FollowStatus {
  switch status {
  case entity.FollowStatusNone:
    return relationv1.FollowStatus_NONE
  case entity.FollowStatusFollower:
    return relationv1.FollowStatus_FOLLOWER
  case entity.FollowStatusMutual:
    return relationv1.FollowStatus_MUTUAL
  default:
    return relationv1.FollowStatus_UNKNOWN
  }
}

func ToProtoFollowCount(count *dto.FollowCountResponseDTO) *relationv1.FollowCount {
  return &relationv1.FollowCount{
    TotalFollower: count.TotalFollowers,
    TotalFollowee: count.TotalFollowings,
  }
}