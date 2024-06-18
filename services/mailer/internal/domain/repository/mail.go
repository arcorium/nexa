package repository

import (
  "context"
  domain "nexa/services/mailer/internal/domain/entity"
  "nexa/shared/types"
  "nexa/shared/util/repo"
)

type IMail interface {
  FindAll(ctx context.Context, query repo.QueryParameter) (repo.PaginatedResult[domain.Mail], error)
  FindByIds(ctx context.Context, ids ...types.Id) ([]domain.Mail, error)
  FindByTag(ctx context.Context, tag types.Id) ([]domain.Mail, error)
  Create(ctx context.Context, mail *domain.Mail) error
  AppendTag(ctx context.Context, mailId types.Id, tagIds []types.Id) error
  RemoveTag(ctx context.Context, mailId types.Id, tagIds []types.Id) error
  Patch(ctx context.Context, mail *domain.Mail) error
  Remove(ctx context.Context, id types.Id) error
}
