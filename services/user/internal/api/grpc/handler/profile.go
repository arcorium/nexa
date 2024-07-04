package handler

import (
  "context"
  middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  "io"
  "nexa/proto/gen/go/user/v1"
  "nexa/services/user/internal/api/grpc/mapper"
  "nexa/services/user/internal/domain/dto"
  "nexa/services/user/internal/domain/service"
  "nexa/services/user/util"
  sharedErr "nexa/shared/errors"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  spanUtil "nexa/shared/util/span"
)

func NewProfileHandler(profile service.IProfile) ProfileHandler {
  return ProfileHandler{
    profileService: profile,
    tracer:         util.GetTracer(),
  }
}

type ProfileHandler struct {
  userv1.UnimplementedProfileServiceServer

  profileService service.IProfile
  tracer         trace.Tracer
}

func (p *ProfileHandler) Register(server *grpc.Server) {
  userv1.RegisterProfileServiceServer(server, p)
}

func (p *ProfileHandler) Find(ctx context.Context, request *userv1.FindProfileRequest) (*userv1.FindProfileResponse, error) {
  ctx, span := p.tracer.Start(ctx, "ProfileHandler.Find")
  defer span.End()

  // Input validation
  ids, ierr := sharedUtil.CastSliceErrs(request.UserIds, types.IdFromString)
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    err := sharedErr.GrpcFieldIndexedErrors("user_ids", ierr)
    return nil, err
  }

  profiles, stats := p.profileService.Find(ctx, ids...)
  response := &userv1.FindProfileResponse{
    Profiles: sharedUtil.CastSliceP(profiles, mapper.ToProtoProfile),
  }

  return response, stats.ToGRPCErrorWithSpan(span)
}

func (p *ProfileHandler) Update(ctx context.Context, request *userv1.UpdateProfileRequest) (*emptypb.Empty, error) {
  ctx, span := p.tracer.Start(ctx, "ProfileHandler.Update")
  defer span.End()

  dtoInput, err := mapper.ToProfileUpdateDTO(request)
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

func (p *ProfileHandler) UpdateAvatar(server userv1.ProfileService_UpdateAvatarServer) error {
  ss := server.(grpc.ServerStream)
  wrappedStream, ok := ss.(*middleware.WrappedServerStream)
  if !ok {
    return server.SendAndClose(nil)
  }

  ctx, span := p.tracer.Start(wrappedStream.WrappedContext, "ProfileHandler.Update")
  defer span.End()

  var userId string
  var filename string
  var bytes []byte
  for {
    recv, err := server.Recv()
    if err != nil {
      if err == io.EOF {
        break
      }
      spanUtil.RecordError(err, span)
      return err
    }
    userId = recv.UserId
    filename = recv.Filename
    bytes = append(bytes, recv.Chunk...)
  }
  id, err := types.IdFromString(userId)
  if err != nil {
    spanUtil.RecordError(err, span)
    err = sharedErr.GrpcFieldErrors2(sharedErr.NewFieldError("user_id", err))
    return err
  }

  // Mapping and Validation
  dtoInput := dto.ProfileAvatarUpdateDTO{
    UserId:   id,
    Filename: filename,
    Bytes:    bytes,
  }

  err = sharedUtil.ValidateStructCtx(ctx, &dtoInput)
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }

  stat := p.profileService.UpdateAvatar(ctx, &dtoInput)
  if stat.IsError() {
    spanUtil.RecordError(stat.Error, span)
    return server.SendAndClose(nil)
  }

  return server.SendAndClose(&emptypb.Empty{})
}
