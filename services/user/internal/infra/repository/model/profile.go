package model

import (
  "github.com/uptrace/bun"
  domain "nexa/services/user/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util/repo"
  "nexa/shared/variadic"
  "nexa/shared/wrapper"
  "time"
)

type ProfileMapOption = repo.DataAccessModelMapOption[*domain.Profile, *Profile]

func FromProfileDomain(domain *domain.Profile, opts ...ProfileMapOption) Profile {
  pfl := Profile{
    UserId:    domain.Id.Underlying().String(),
    FirstName: domain.FirstName,
    LastName:  domain.LastName,
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
  PhotoURL  string `bun:",nullzero"` // TODO: Change into uuid
  Bio       string `bun:",nullzero"`

  UpdatedAt time.Time `bun:",nullzero"`

  User *User `bun:"rel:belongs-to,join:user_id=id,on_delete:CASCADE"`
}

func (p *Profile) ToDomain() domain.Profile {
  return domain.Profile{
    Id:        wrapper.DropError(types.IdFromString(p.UserId)),
    FirstName: p.FirstName,
    LastName:  p.LastName,
    PhotoURL:  types.FilePathFromString(p.PhotoURL),
    Bio:       p.Bio,
  }
}
