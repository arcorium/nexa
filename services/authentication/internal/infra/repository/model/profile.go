package model

import (
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  "github.com/arcorium/nexa/shared/variadic"
  "github.com/uptrace/bun"
  "nexa/services/authentication/internal/domain/entity"
  "time"
)

type ProfileMapOption = repo.DataAccessModelMapOption[*entity.Profile, *Profile]
type PatchedProfileMapOption = repo.DataAccessModelMapOption[*entity.PatchedProfile, *Profile]

func FromPatchedProfileDomain(ent *entity.PatchedProfile, opts ...PatchedProfileMapOption) Profile {
  profile := Profile{
    Id:        ent.Id.String(),
    FirstName: ent.FirstName,
    LastName:  ent.LastName.ValueOrNil(),
    PhotoId:   types.GetValueOrNilCasted(ent.PhotoId, sharedUtil.ToString[types.Id]),
    PhotoURL:  types.GetValueOrNilCasted(ent.PhotoURL, sharedUtil.ToString[types.FilePath]),
    Bio:       ent.Bio.ValueOrNil(),
  }

  variadic.New(opts...).DoAll(repo.MapOptionFunc(ent, &profile))

  return profile
}

func FromProfileDomain(ent *entity.Profile, opts ...ProfileMapOption) Profile {
  photoId := ent.PhotoId.String()
  photoUrl := ent.PhotoURL.String()

  profile := Profile{
    Id:        ent.Id.String(),
    UserId:    ent.UserId.String(),
    FirstName: ent.FirstName,
    LastName:  &ent.LastName,
    PhotoId:   &photoId,
    PhotoURL:  &photoUrl,
    Bio:       &ent.Bio,
  }

  variadic.New(opts...).DoAll(repo.MapOptionFunc(ent, &profile))

  return profile
}

type Profile struct {
  bun.BaseModel `bun:"table:profiles,alias:p"`

  Id        string  `bun:",type:uuid,pk"`
  UserId    string  `bun:",type:uuid,notnull,nullzero"` // Profile is unique per user
  FirstName string  `bun:",notnull,nullzero"`
  LastName  *string `bun:","`
  PhotoId   *string `bun:",type:uuid"` // Id of profile image on file storage
  PhotoURL  *string `bun:","`
  Bio       *string `bun:","`

  UpdatedAt time.Time `bun:",nullzero"`

  User *User `bun:"rel:belongs-to,join:user_id=id,on_delete:CASCADE"` // Only used as reference
}

func (p *Profile) ToDomain() (entity.Profile, error) {
  id, err := types.IdFromString(p.Id)
  if err != nil {
    return entity.Profile{}, err
  }

  userId, err := types.IdFromString(p.UserId)
  if err != nil {
    return entity.Profile{}, err
  }

  photoId, err := types.IdFromString(types.OnNil(p.PhotoId, ""))
  if err != nil {
    return entity.Profile{}, err
  }

  return entity.Profile{
    Id:        id,
    UserId:    userId,
    FirstName: p.FirstName,
    LastName:  types.OnNil(p.LastName, ""),
    PhotoId:   photoId,
    PhotoURL:  types.FilePathFromString(types.OnNil(p.PhotoURL, "")),
    Bio:       types.OnNil(p.Bio, ""),
  }, nil
}
