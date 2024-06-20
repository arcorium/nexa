package handler

import (
  "context"
  "errors"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  "io"
  proto "nexa/proto/generated/golang/user/v1"
  "nexa/services/user/internal/api/grpc/mapper"
  "nexa/services/user/internal/domain/dto"
  "nexa/services/user/internal/domain/service"
  sharedErr "nexa/shared/errors"
  spanUtil "nexa/shared/span"
  "nexa/shared/types"
  "nexa/shared/util"
)

func NewProfileHandler(profile service.IProfile) ProfileHandler {
  return ProfileHandler{
    profileService: profile,
  }
}

type ProfileHandler struct {
  proto.UnimplementedProfileServiceServer

  profileService service.IProfile
}

func (p *ProfileHandler) Register(server *grpc.Server) {
  proto.RegisterProfileServiceServer(server, p)
}

func (p *ProfileHandler) Find(ctx context.Context, request *proto.FindProfileRequest) (*proto.FindProfileResponse, error) {
  span := trace.SpanFromContext(ctx)

  ids, ierr := util.CastSliceErrsP(request.UserIds, func(from *string) (types.Id, error) {
    return types.IdFromString(*from)
  })

  if ierr != nil {
    errs := util.CastSlice(ierr, func(from sharedErr.IndexedError) error {
      return from.Err
    })
    err := errors.Join(errs...)
    spanUtil.RecordError(err, span)
    return nil, sharedErr.GrpcFieldIndexedErrors("user_id", ierr)
  }

  profiles, stats := p.profileService.Find(ctx, ids)
  response := &proto.FindProfileResponse{
    Profiles: util.CastSliceP(profiles, mapper.ToProtoProfile),
  }

  return response, stats.ToGRPCErrorWithSpan(span)
}

func (p *ProfileHandler) Update(ctx context.Context, request *proto.UpdateProfileRequest) (*emptypb.Empty, error) {
  ctx, span := p.tracer.Start(ctx, "ProfileHandler.Update")
  defer span.End()

  dtoInput := mapper.ToProfileUpdateDTO(request)

  err := util.ValidateStruct(ctx, &dtoInput)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  stats := p.profileService.Update(ctx, &dtoInput)
  if stats.IsError() {
    spanUtil.RecordError(stats.Error, span)
    return nil, stats.ToGRPCError()
  }
  return &emptypb.Empty{}, nil
}

func (p *ProfileHandler) UpdateAvatar(server proto.ProfileService_UpdateAvatarServer) error {
  ctx, span := p.tracer.Start(server.Context(), "ProfileHandler.UpdateAvatar")
  defer span.End()

  dtoInput := dto.ProfilePictureUpdateDTO{}
  for {
    recv, err := server.Recv()
    if err != nil {
      if err == io.EOF {
        break
      }
      return err
    }
    dtoInput.UserId = recv.UserId
    dtoInput.Filename = recv.Filename
    dtoInput.Bytes = append(dtoInput.Bytes, recv.Chunk...)
  }

  stat := p.profileService.UpdateAvatar(ctx, &dtoInput)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return server.SendAndClose(nil)
  }

  return server.SendAndClose(&emptypb.Empty{})
}
