package model

import (
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util/repo"
  "github.com/arcorium/nexa/shared/variadic"
  "github.com/uptrace/bun"
  "nexa/services/relation/internal/domain/entity"
  "time"
)

type FollowMapOption = repo.DataAccessModelMapOption[*entity.Follow, *Follow]

func FromFollowDomain(ent *entity.Follow, opts ...FollowMapOption) Follow {
  follow := Follow{
    FollowerId: ent.FollowerId.String(),
    FolloweeId: ent.FolloweeId.String(),
    CreatedAt:  ent.CreatedAt,
  }

  variadic.New(opts...).
    DoAll(repo.MapOptionFunc(ent, &follow))

  return follow
}

type Follow struct {
  bun.BaseModel `bun:"table:follows"`

  FollowerId string `bun:",type:uuid,nullzero,pk"`
  FolloweeId string `bun:",type:uuid,nullzero,pk"`

  CreatedAt time.Time `bun:",nullzero,notnull"`
}

func (f *Follow) ToDomain() (entity.Follow, error) {
  followerId, err := types.IdFromString(f.FollowerId)
  if err != nil {
    return entity.Follow{}, err
  }

  followeeId, err := types.IdFromString(f.FolloweeId)
  if err != nil {
    return entity.Follow{}, err
  }

  return entity.Follow{
    FollowerId: followerId,
    FolloweeId: followeeId,
    CreatedAt:  f.CreatedAt,
  }, nil
}
