package model

import (
  domain "nexa/services/authentication/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util"
  "nexa/shared/util/repo"
  "nexa/shared/variadic"
  "time"
)

type CredentialMapOption = repo.DataAccessModelMapOption[*domain.Credential, *Credential]

func FromCredentialModel(domain *domain.Credential, opts ...CredentialMapOption) Credential {

  cred := Credential{
    Id:            util.ReturnOnEqual(domain.Id.Underlying().String(), types.NullIdStr, ""),
    UserId:        util.ReturnOnEqual(domain.UserId.Underlying().String(), types.NullIdStr, ""),
    AccessTokenId: util.ReturnOnEqual(domain.AccessTokenId.Underlying().String(), types.NullIdStr, ""),
    Device:        domain.Device.Name,
    RefreshToken:  domain.RefreshToken,
    ExpiresAt:     domain.ExpiresAt,
  }

  variadic.New(opts...).
    DoAll(repo.MapOptionFunc(domain, &cred))

  return cred
}

type Credential struct {
  Id            string    `redis:"-"`
  UserId        string    `redis:"user_id"`
  AccessTokenId string    `redis:"access_token_id"`
  Device        string    `redis:"device"`
  RefreshToken  string    `redis:"refresh_token"`
  ExpiresAt     time.Time `redis:"-"`
}

// OmitZero used for data patching
func (c *Credential) OmitZero() map[string]any {
  result := make(map[string]any)

  if c.AccessTokenId != "" {
    result["access_token_id"] = c.AccessTokenId
  }

  if c.Device != "" {
    result["device"] = c.Device
  }

  if c.RefreshToken != "" {
    result["refresh_token"] = c.RefreshToken
  }

  return result
}

func (c *Credential) ToDomain() domain.Credential {
  return domain.Credential{
    Id:            types.DropError(types.IdFromString(c.Id)),
    UserId:        types.DropError(types.IdFromString(c.UserId)),
    AccessTokenId: types.DropError(types.IdFromString(c.AccessTokenId)),
    Device: domain.Device{
      Name: c.Device,
    },
    RefreshToken: c.RefreshToken,
    ExpiresAt:    c.ExpiresAt,
  }
}
