package mapper

import (
  "nexa/services/relation/internal/domain/dto"
  "nexa/services/relation/internal/domain/entity"
)

func ToFollowerResponseDTO(follow *entity.Follow) dto.FollowResponseDTO {
  return dto.FollowResponseDTO{
    UserId:    follow.FollowerId,
    CreatedAt: follow.CreatedAt,
  }
}

func ToFolloweeResponseDTO(follow *entity.Follow) dto.FollowResponseDTO {
  return dto.FollowResponseDTO{
    UserId:    follow.FolloweeId,
    CreatedAt: follow.CreatedAt,
  }
}

func ToFollowStatus(vals bool) entity.FollowStatus {
  if vals {
    return entity.FollowStatusFollower
  }
  return entity.FollowStatusNone
}

func ToFollowCountResponseDTO(followCount *entity.FollowCount) dto.FollowCountResponseDTO {
  return dto.FollowCountResponseDTO{
    UserId:          followCount.UserId,
    TotalFollowers:  followCount.TotalFollowers,
    TotalFollowings: followCount.TotalFollowings,
  }
}
