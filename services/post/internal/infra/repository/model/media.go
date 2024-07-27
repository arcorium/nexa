package model

import (
  "github.com/arcorium/nexa/shared/types"
  "github.com/uptrace/bun"
  "nexa/services/post/internal/domain/entity"
)

type Media struct {
  bun.BaseModel `bun:"table:medias"`
  Id            uint64 `bun:",autoincrement,pk"`
  VersionId     string `bun:",type:uuid,notnull,nullzero"`
  FileId        string `bun:",type:uuid,notnull,nullzero"`

  PostVersion *PostVersion `bun:"rel:belongs-to,join:version_id=id,on_delete:CASCADE"`
}

func (m *Media) ToDomain() (entity.Media, error) {
  fileId, err := types.IdFromString(m.FileId)
  if err != nil {
    return entity.Media{}, err
  }

  return entity.Media{
    Id: fileId,
    //Url: "", // Set on service
  }, nil
}
