package repository

import (
  "context"
  "nexa/services/authentication/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util/repo"
)

type ICredential interface {
  Create(ctx context.Context, credential *entity.Credential) error
  Patch(ctx context.Context, credential *entity.Credential) error
  Delete(ctx context.Context, credIds ...types.Id) error
  // DeleteByUserId delete user credentials based on the user id
  DeleteByUserId(ctx context.Context, userId types.Id, credIds ...types.Id) error
  Find(ctx context.Context, refreshTokenId types.Id) (*entity.Credential, error)
  FindByUserId(ctx context.Context, userId types.Id) ([]entity.Credential, error)
  Get(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Credential], error)
}
