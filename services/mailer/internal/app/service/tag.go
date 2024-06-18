package service

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/mailer/internal/domain/dto"
  domain "nexa/services/mailer/internal/domain/entity"
  "nexa/services/mailer/internal/domain/mapper"
  "nexa/services/mailer/internal/domain/repository"
  "nexa/services/mailer/internal/domain/service"
  "nexa/services/mailer/util"
  sharedDto "nexa/shared/dto"
  spanUtil "nexa/shared/span"
  "nexa/shared/status"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
)

func NewTag(repo repository.ITag) service.ITag {
  return &tagService{
    tagRepo: repo,
    tracer:  util.GetTracer(),
  }
}

type tagService struct {
  tagRepo repository.ITag
  tracer  trace.Tracer
}

func (t *tagService) Find(ctx context.Context, elementDTO *sharedDto.PagedElementDTO) (*sharedDto.PagedElementResult[dto.TagResponseDTO], status.Object) {
  ctx, span := t.tracer.Start(ctx, "TagService.Find")
  defer span.End()

  result, err := t.tagRepo.FindAll(ctx, elementDTO.ToQueryParam())
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  tags := sharedUtil.CastSlice(result.Data, func(tag *domain.Tag) dto.TagResponseDTO {
    return mapper.ToResponseDTO(tag)
  })

  res := sharedDto.NewPagedElementResult2(tags, elementDTO, result.Total)
  return &res, status.Success()
}

func (t *tagService) FindByIds(ctx context.Context, ids ...types.Id) ([]dto.TagResponseDTO, status.Object) {
  ctx, span := t.tracer.Start(ctx, "TagService.FindByIds")
  defer span.End()

  result, err := t.tagRepo.FindByIds(ctx, ids...)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }

  tags := sharedUtil.CastSlice(result, func(tag *domain.Tag) dto.TagResponseDTO {
    return mapper.ToResponseDTO(tag)
  })

  return tags, status.Success()
}

func (t *tagService) FindByName(ctx context.Context, name string) (dto.TagResponseDTO, status.Object) {
  ctx, span := t.tracer.Start(ctx, "TagService.FindByName")
  defer span.End()

  result, err := t.tagRepo.FindByName(ctx, name)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.TagResponseDTO{}, status.FromRepository(err, status.NullCode)
  }

  return mapper.ToResponseDTO(result), status.Success()
}

func (t *tagService) Create(ctx context.Context, dto *dto.CreateTagDTO) (types.Id, status.Object) {
  ctx, span := t.tracer.Start(ctx, "TagService.Create")
  defer span.End()

  tag := mapper.MapCreateTagDTO(dto)

  err := t.tagRepo.Create(ctx, &tag)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), status.FromRepository(err, status.NullCode)
  }

  return tag.Id, status.Created()
}

func (t *tagService) Update(ctx context.Context, tagDto *dto.UpdateTagDTO) status.Object {
  ctx, span := t.tracer.Start(ctx, "TagService.Update")
  defer span.End()

  tag := mapper.MapUpdateTagDTO(tagDto)
  err := t.tagRepo.Patch(ctx, &tag)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Updated()
}

func (t *tagService) Remove(ctx context.Context, id types.Id) status.Object {
  ctx, span := t.tracer.Start(ctx, "TagService.Remove")
  defer span.End()

  err := t.tagRepo.Remove(ctx, id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  return status.Deleted()
}
