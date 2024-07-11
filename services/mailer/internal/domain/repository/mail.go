package repository

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  "github.com/arcorium/nexa/shared/util/repo"
  domain "nexa/services/mailer/internal/domain/entity"
)

type MailTags = types.Pair[types.Id, []types.Id]

type IMail interface {
  Get(ctx context.Context, query repo.QueryParameter) (repo.PaginatedResult[domain.Mail], error)
  FindByIds(ctx context.Context, ids ...types.Id) ([]domain.Mail, error)
  FindByTag(ctx context.Context, tag types.Id) ([]domain.Mail, error)
  Create(ctx context.Context, mails ...domain.Mail) error
  AppendTags(ctx context.Context, mailId types.Id, tagIds []types.Id) error
  RemoveTags(ctx context.Context, mailId types.Id, tagIds []types.Id) error
  // Patch only able to set several fields, that are status, updated_at and delivered_at. Other fields are supposed to read only
  Patch(ctx context.Context, mail *domain.Mail) error
  Remove(ctx context.Context, id types.Id) error
  // AppendMultipleTags insert multiple tags for each mail, the first data is the mail id and the second
  // is the ids of the tag
  AppendMultipleTags(ctx context.Context, mailTags ...MailTags) error
}
