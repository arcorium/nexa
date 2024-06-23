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
  Delete(ctx context.Context, accessTokenId ...types.Id) error
  DeleteByUserId(ctx context.Context, userId types.Id) error
  Find(ctx context.Context, accessTokenId types.Id) (entity.Credential, error)
  FindByUserId(ctx context.Context, userId types.Id) ([]entity.Credential, error)
  FindAll(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Credential], error)
}
