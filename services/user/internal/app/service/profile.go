package service

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  userUow "nexa/services/user/internal/app/uow"
  "nexa/services/user/internal/domain/dto"
  "nexa/services/user/internal/domain/mapper"
  "nexa/services/user/internal/domain/service"
  util2 "nexa/services/user/util"
  spanUtil "nexa/shared/span"
  "nexa/shared/status"
  "nexa/shared/types"
  "nexa/shared/uow"
  "nexa/shared/util"
)

func NewProfile(uow uow.IUnitOfWork[userUow.UserStorage]) service.IProfile {
  return &profileService{
    unit:   uow,
    tracer: util2.GetTracer(),
  }
}

type profileService struct {
  unit uow.IUnitOfWork[userUow.UserStorage]

  tracer trace.Tracer
}

func (p profileService) Find(ctx context.Context, userIds []types.Id) ([]dto.ProfileResponseDTO, status.Object) {
  ctx, span := p.tracer.Start(ctx, "ProfileService.Find")
  defer span.End()

  repo := p.unit.Repositories()
  profiles, err := repo.Profile().FindByIds(ctx, userIds...)

  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }
  return util.CastSliceP(profiles, mapper.ToProfileResponse), status.Success()
}

func (p profileService) Update(ctx context.Context, input *dto.ProfileUpdateDTO) status.Object {
  ctx, span := p.tracer.Start(ctx, "ProfileService.Update")
  defer span.End()

  profile := mapper.MapProfileUpdateDTO(input)
  repo := p.unit.Repositories()
  err := repo.Profile().Patch(ctx, &profile)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Updated()
}

func (p profileService) UpdateAvatar(ctx context.Context, input *dto.ProfilePictureUpdateDTO) status.Object {
  ctx, span := p.tracer.Start(ctx, "ProfileService.UpdateAvatar")
  defer span.End()
  // TODO: Communicate with FileStorage Service

  url := types.FilePath("")
  profile := mapper.MapProfilePictureUpdateDTO(input)
  profile.PhotoURL = url

  repo := p.unit.Repositories()
  err := repo.Profile().Patch(ctx, &profile)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Updated()
}
