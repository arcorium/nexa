package model

import (
  "github.com/uptrace/bun"
  "nexa/services/authentication/shared/domain/entity"
  "nexa/services/authentication/shared/domain/valueobject"
  "nexa/shared/types"
  "nexa/shared/util/repo"
  "nexa/shared/variadic"
  "time"
)

type CredentialMapOption = repo.DataAccessModelMapOption[*entity.Credential, *Credential]

func FromCredentialModel(domain *entity.Credential, opts ...CredentialMapOption) Credential {
  cred := Credential{
    UserId:        domain.UserId.Underlying().String(),
    AccessTokenId: domain.AccessTokenId.Underlying().String(),
    Device:        domain.Device,
    Token:         domain.RefreshToken,
  }

  variadic.New(opts...).
    DoAll(repo.MapOptionFunc(domain, &cred))

  return cred
}

type Credential struct {
  bun.BaseModel `bun:"table:credentials"`

  UserId        string             `bun:",nullzero,type:uuid,pk"`
  AccessTokenId string             `bun:"access_id,nullzero,notnull,type:uuid,pk"`
  Device        valueobject.Device `bun:",nullzero,embed:device_"`
  Token         string             `bun:",nullzero,notnull"` // Refresh token

  UpdatedAt time.Time `bun:",nullzero"`
  CreatedAt time.Time `bun:",nullzero,notnull"`
}

func (c *Credential) ToDomain() entity.Credential {
  return entity.Credential{
    UserId:        types.IdFromString(c.UserId),
    AccessTokenId: types.IdFromString(c.AccessTokenId),
    Device:        c.Device,
    RefreshToken:  c.Token,
  }
}
