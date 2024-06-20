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
  Create(ctx context.Context, mails ...domain.Mail) error
  AppendTags(ctx context.Context, mailId types.Id, tagIds []types.Id) error
  RemoveTags(ctx context.Context, mailId types.Id, tagIds []types.Id) error
  Patch(ctx context.Context, mail *domain.Mail) error
  Remove(ctx context.Context, id types.Id) error
  AppendMultipleTags(ctx context.Context, mailTags ...types.Pair[types.Id, []types.Id]) error
}
