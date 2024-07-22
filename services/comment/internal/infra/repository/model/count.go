package model

import "nexa/services/comment/internal/domain/entity"

type Count struct {
  Id            string `bun:"ids,scanonly"`
  TotalComments uint64 `bun:",scanonly"`
}

func (c *Count) ToDomain() entity.Count {
  return entity.Count{
    TotalComments: c.TotalComments,
  }
}
