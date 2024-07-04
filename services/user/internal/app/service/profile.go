package service

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/user/internal/domain/dto"
  "nexa/services/user/internal/domain/entity"
  "nexa/services/user/internal/domain/external"
  "nexa/services/user/internal/domain/mapper"
  "nexa/services/user/internal/domain/repository"
  "nexa/services/user/internal/domain/service"
  util2 "nexa/services/user/util"
  "nexa/shared/status"
  "nexa/shared/types"
  "nexa/shared/util"
  spanUtil "nexa/shared/util/span"
)

func NewProfile(repo repository.IProfile, storageExt external.IFileStorageClient) service.IProfile {
  return &profileService{
    profileRepo: repo,
    storageExt:  storageExt,
    tracer:      util2.GetTracer(),
  }
}

type profileService struct {
  profileRepo repository.IProfile
  storageExt  external.IFileStorageClient

  tracer trace.Tracer
}

func (p profileService) Find(ctx context.Context, userIds ...types.Id) ([]dto.ProfileResponseDTO, status.Object) {
  ctx, span := p.tracer.Start(ctx, "ProfileService.Find")
  defer span.End()

  profiles, err := p.profileRepo.FindByIds(ctx, userIds...)

  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }
  return util.CastSliceP(profiles, mapper.ToProfileResponse), status.Success()
}

func (p profileService) Update(ctx context.Context, input *dto.ProfileUpdateDTO) status.Object {
  ctx, span := p.tracer.Start(ctx, "ProfileService.Update")
  defer span.End()

  profile := input.ToDomain()

  err := p.profileRepo.Patch(ctx, &profile)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Updated()
}

func (p profileService) UpdateAvatar(ctx context.Context, input *dto.ProfileAvatarUpdateDTO) status.Object {
  ctx, span := p.tracer.Start(ctx, "ProfileService.UpdateAvatar")
  defer span.End()

  // Check if user already has photo
  profiles, err := p.profileRepo.FindByIds(ctx, input.UserId)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  // Upload new avatar
  fileId, filePath, err := p.storageExt.UploadProfileImage(ctx, &dto.UploadImageDTO{
    Filename: input.Filename,
    Data:     input.Bytes,
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrExternal(err)
  }

  // Update profiles data
  profile := entity.Profile{
    Id:       input.UserId,
    PhotoId:  fileId,
    PhotoURL: filePath,
  }

  err = p.profileRepo.Patch(ctx, &profile)
  if err != nil {
    // Delete new avatar when error happens
    spanUtil.RecordError(err, span)
    err = p.storageExt.DeleteProfileImage(ctx, fileId)
    if err != nil {
      spanUtil.RecordError(err, span)
      return status.ErrExternal(err)
    }
    return status.FromRepository(err, status.NullCode)
  }

  // Delete last avatar
  if profiles[0].HasAvatar() {
    err = p.storageExt.DeleteProfileImage(ctx, profiles[0].PhotoId)
    if err != nil {
      spanUtil.RecordError(err, span)
      return status.ErrExternal(err)
    }
  }

  return status.Updated()
}
