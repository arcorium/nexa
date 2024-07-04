package model

import (
  "github.com/uptrace/bun"
  domain "nexa/services/user/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util/repo"
  "nexa/shared/variadic"
  "time"
)

type ProfileMapOption = repo.DataAccessModelMapOption[*domain.Profile, *Profile]

func FromProfileDomain(domain *domain.Profile, opts ...ProfileMapOption) Profile {
  pfl := Profile{
    UserId:    domain.Id.String(),
    FirstName: domain.FirstName,
    LastName:  domain.LastName,
    PhotoId:   domain.PhotoId.String(),
    PhotoURL:  domain.PhotoURL.FileName(),
    Bio:       domain.Bio,
  }

  variadic.New(opts...).DoAll(repo.MapOptionFunc(domain, &pfl))

  return pfl
}

type Profile struct {
  bun.BaseModel `bun:"table:profiles"`

  UserId    string `bun:",type:uuid,pk"` // Profile is unique per user
  FirstName string `bun:",notnull"`
  LastName  string `bun:",nullzero"`
  PhotoId   string `bun:",type:uuid,notnull"` // Id of profile image on file storage
  PhotoURL  string `bun:",nullzero"`
  Bio       string `bun:",nullzero"`

  UpdatedAt time.Time `bun:",nullzero"`

  //User *User `bun:"rel:belongs-to,join:user_id=id,on_delete:CASCADE"`
}

func (p *Profile) ToDomain() (domain.Profile, error) {
  userId, err := types.IdFromString(p.UserId)
  if err != nil {
    return domain.Profile{}, err
  }

  photoId, err := types.IdFromString(p.PhotoId)
  if err != nil {
    return domain.Profile{}, err
  }

  return domain.Profile{
    Id:        userId,
    FirstName: p.FirstName,
    LastName:  p.LastName,
    PhotoId:   photoId,
    PhotoURL:  types.FilePathFromString(p.PhotoURL),
    Bio:       p.Bio,
  }, nil
}
