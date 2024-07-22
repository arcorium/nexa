package model

import (
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  "github.com/arcorium/nexa/shared/variadic"
  "nexa/services/authentication/internal/domain/entity"
  "time"
)

type CredentialMapOption = repo.DataAccessModelMapOption[*entity.Credential, *Credential]

func FromCredentialModel(domain *entity.Credential, opts ...CredentialMapOption) Credential {

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

func (c *Credential) ToDomain() entity.Credential {
  return entity.Credential{
    Id:            types.DropError(types.IdFromString(c.Id)),
    UserId:        types.DropError(types.IdFromString(c.UserId)),
    AccessTokenId: types.DropError(types.IdFromString(c.AccessTokenId)),
    Device: entity.Device{
      Name: c.Device,
    },
    RefreshToken: c.RefreshToken,
    ExpiresAt:    c.ExpiresAt,
  }
}
