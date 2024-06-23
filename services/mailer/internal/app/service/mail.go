package service

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "log"
  "nexa/services/mailer/internal/app/uow"
  "nexa/services/mailer/internal/domain/dto"
  domain "nexa/services/mailer/internal/domain/entity"
  "nexa/services/mailer/internal/domain/external"
  "nexa/services/mailer/internal/domain/mapper"
  "nexa/services/mailer/internal/domain/service"
  "nexa/services/mailer/util"
  sharedDto "nexa/shared/dto"
  spanUtil "nexa/shared/span"
  "nexa/shared/status"
  "nexa/shared/types"
  sharedUOW "nexa/shared/uow"
  sharedUtil "nexa/shared/util"
  "nexa/shared/wrapper"
  "sync/atomic"
)

func NewMail(work sharedUOW.IUnitOfWork[uow.MailStorage], mailExt external.IMail) service.IMail {
  return &mailService{
    mailUow: work,
    mailExt: mailExt,
    tracer:  util.GetTracer(),
  }
}

type mailService struct {
  mailExt external.IMail
  mailUow sharedUOW.IUnitOfWork[uow.MailStorage]
  tracer  trace.Tracer

  workCount atomic.Int64
}

func (m *mailService) Find(ctx context.Context, pagedDTO *sharedDto.PagedElementDTO) (*sharedDto.PagedElementResult[dto.MailResponseDTO], status.Object) {
  ctx, span := m.tracer.Start(ctx, "MailService.Find")
  defer span.End()

  repos := m.mailUow.Repositories()
  result, err := repos.Mail().FindAll(ctx, pagedDTO.ToQueryParam())
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  mails := sharedUtil.CastSliceP(result.Data, func(mail *domain.Mail) dto.MailResponseDTO {
    return mapper.ToMailResponseDTO(mail)
  })

  pagedResult := sharedDto.NewPagedElementResult2(mails, pagedDTO, result.Total)
  return &pagedResult, status.Success()
}

func (m *mailService) FindByIds(ctx context.Context, mailIds ...types.Id) ([]dto.MailResponseDTO, status.Object) {
  ctx, span := m.tracer.Start(ctx, "MailService.FindByIds")
  defer span.End()

  repos := m.mailUow.Repositories()
  mails, err := repos.Mail().FindByIds(ctx, mailIds...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  mailResponses := sharedUtil.CastSliceP(mails, func(mail *domain.Mail) dto.MailResponseDTO {
    return mapper.ToMailResponseDTO(mail)
  })

  return mailResponses, status.Success()
}

func (m *mailService) FindByTag(ctx context.Context, tagId types.Id) ([]dto.MailResponseDTO, status.Object) {
  ctx, span := m.tracer.Start(ctx, "MailService.FindByTag")
  defer span.End()

  repos := m.mailUow.Repositories()
  mails, err := repos.Mail().FindByTag(ctx, tagId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  mailResponses := sharedUtil.CastSliceP(mails, func(mail *domain.Mail) dto.MailResponseDTO {
    return mapper.ToMailResponseDTO(mail)
  })

  return mailResponses, status.Success()
}

func (m *mailService) Send(ctx context.Context, mailDTO *dto.SendMailDTO) ([]types.Id, status.Object) {
  ctx, span := m.tracer.Start(ctx, "MailService.Send")
  defer span.End()

  mails := mapper.MapSendMailDTO(mailDTO)

  // Save metadata
  err := m.mailUow.DoTx(ctx, func(ctx context.Context, storage uow.MailStorage) error {
    ctx, txSpan := m.tracer.Start(ctx, "UOW.Send")
    defer txSpan.End()

    err := storage.Mail().Create(ctx, mails...)
    if err != nil {
      spanUtil.RecordError(err, txSpan)
      return err
    }

    var tags []types.Pair[types.Id, []types.Id]

    for _, mail := range mails {
      result := types.NewPair(mail.Id, sharedUtil.CastSliceP(mail.Tags, func(from *domain.Tag) types.Id {
        return from.Id
      }))

      tags = append(tags, result)
    }

    err = storage.Mail().AppendMultipleTags(ctx, tags...)
    return err
  })

  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  // Send mails
  go func() {
    ctx := context.Background()
    defer m.workCount.Add(int64(len(mails) * -1))

    for _, mail := range mails {
      m.workCount.Add(1)

      err = m.mailExt.Send(ctx, &mail, mailDTO.Attachments)
      stat := domain.StatusDelivered
      if err != nil {
        stat = domain.StatusFailed
      }
      repos := m.mailUow.Repositories()

      newMail := domain.Mail{
        Id:     mail.Id,
        Status: stat,
      }

      err = repos.Mail().Patch(ctx, &newMail)
      if err != nil {
        log.Println("Error set mail status:", err)
      }
    }
  }()

  mailIds := sharedUtil.CastSliceP(mails, func(from *domain.Mail) types.Id {
    return from.Id
  })

  return mailIds, status.Created()
}

func (m *mailService) Update(ctx context.Context, mailDTO *dto.UpdateMailDTO) status.Object {
  ctx, span := m.tracer.Start(ctx, "MailService.Update")
  defer span.End()

  mailId := wrapper.DropError(types.IdFromString(mailDTO.Id))

  addedTagIds := sharedUtil.CastSlice(mailDTO.AddedTagIds, func(from string) types.Id {
    return wrapper.DropError(types.IdFromString(from))
  })

  removedTagIds := sharedUtil.CastSlice(mailDTO.RemovedTagIds, func(from string) types.Id {
    return wrapper.DropError(types.IdFromString(from))
  })

  err := m.mailUow.DoTx(ctx, func(ctx context.Context, storage uow.MailStorage) error {
    ctx, txSpan := m.tracer.Start(ctx, "UOW.Update")
    defer txSpan.End()

    // Append Tags
    if len(addedTagIds) > 0 {
      err := storage.Mail().AppendTags(ctx, mailId, addedTagIds)
      if err != nil {
        spanUtil.RecordError(err, span)
        return err
      }
    }

    // Remove Tags
    if len(removedTagIds) > 0 {
      err := storage.Mail().RemoveTags(ctx, mailId, removedTagIds)
      if err != nil {
        spanUtil.RecordError(err, span)
        return err
      }
    }

    return nil
  })

  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Updated()
}

func (m *mailService) Remove(ctx context.Context, mailId types.Id) status.Object {
  ctx, span := m.tracer.Start(ctx, "MailService.Remove")
  defer span.End()

  repos := m.mailUow.Repositories()
  err := repos.Mail().Remove(ctx, mailId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Deleted()
}

func (m *mailService) HasWork() bool {
  return m.workCount.Load() <= 0
}
