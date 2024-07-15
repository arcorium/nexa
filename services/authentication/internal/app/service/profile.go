package service

import (
  "context"
  "github.com/arcorium/nexa/shared/status"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "go.opentelemetry.io/otel/trace"
  dto2 "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/entity"
  "nexa/services/authentication/internal/domain/external"
  "nexa/services/authentication/internal/domain/mapper"
  "nexa/services/authentication/internal/domain/repository"
  "nexa/services/authentication/internal/domain/service"
  "nexa/services/authentication/util"
)

func NewProfile(repo repository.IProfile, storageExt external.IFileStorageClient) service.IProfile {
  return &profileService{
    profileRepo: repo,
    storageExt:  storageExt,
    tracer:      util.GetTracer(),
  }
}

type profileService struct {
  profileRepo repository.IProfile
  storageExt  external.IFileStorageClient

  tracer trace.Tracer
}

func (p profileService) Find(ctx context.Context, userIds ...types.Id) ([]dto2.ProfileResponseDTO, status.Object) {
  ctx, span := p.tracer.Start(ctx, "ProfileService.Find")
  defer span.End()

  profiles, err := p.profileRepo.FindByIds(ctx, userIds...)

  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, status.FromRepository(err, status.NullCode)
  }
  return sharedUtil.CastSliceP(profiles, mapper.ToProfileResponse), status.Success()
}

func (p profileService) Update(ctx context.Context, updateDto *dto2.ProfileUpdateDTO) status.Object {
  ctx, span := p.tracer.Start(ctx, "ProfileService.Update")
  defer span.End()

  profile := updateDto.ToDomain()

  err := p.profileRepo.Patch(ctx, &profile)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }
  return status.Updated()
}

func (p profileService) UpdateAvatar(ctx context.Context, updateDto *dto2.ProfileAvatarUpdateDTO) status.Object {
  ctx, span := p.tracer.Start(ctx, "ProfileService.UpdateAvatar")
  defer span.End()

  // Check if user already has photo
  profiles, err := p.profileRepo.FindByIds(ctx, updateDto.Id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.FromRepository(err, status.NullCode)
  }

  // Upload new avatar
  fileId, filePath, err := p.storageExt.UploadProfileImage(ctx, &dto2.UploadImageDTO{
    Filename: updateDto.Filename,
    Data:     updateDto.Bytes,
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return status.ErrExternal(err)
  }

  // Update profiles data
  profile := entity.PatchedProfile{
    Id:       updateDto.Id,
    PhotoId:  types.SomeNullable(fileId),
    PhotoURL: types.SomeNullable(filePath),
  }

  err = p.profileRepo.Patch(ctx, &profile)
  if err != nil {
    // Delete new avatar when error happens
    spanUtil.RecordError(err, span)
    extErr := p.storageExt.DeleteProfileImage(ctx, fileId)
    if extErr != nil {
      spanUtil.RecordError(extErr, span)
      return status.ErrExternal(extErr)
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
