package service

import (
  "context"
  "nexa/services/mailer/internal/domain/dto"
  sharedDto "nexa/shared/dto"
  "nexa/shared/status"
  "nexa/shared/types"
)

type IMail interface {
  Find(ctx context.Context, pagedDTO *sharedDto.PagedElementDTO) (*sharedDto.PagedElementResult[dto.MailResponseDTO], status.Object)
  FindByIds(ctx context.Context, mailIds ...types.Id) ([]dto.MailResponseDTO, status.Object)
  FindByTag(ctx context.Context, tagId types.Id) ([]dto.MailResponseDTO, status.Object)
  Send(ctx context.Context, mailDTO *dto.SendMailDTO) (types.Id, status.Object)
  Update(ctx context.Context, mailDTO *dto.UpdateMailDTO) status.Object
  Remove(ctx context.Context, mailId types.Id) status.Object
}
