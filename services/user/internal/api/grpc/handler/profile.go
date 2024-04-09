package handler

import (
  "context"
  "google.golang.org/grpc"
  "google.golang.org/protobuf/types/known/emptypb"
  "io"
  proto "nexa/proto/generated/golang/user/v1"
  "nexa/services/user/internal/api/grpc/mapper"
  "nexa/services/user/internal/domain/dto"
  "nexa/services/user/internal/domain/service"
  sharedErr "nexa/shared/errors"
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
  ids, ierr := util.CastSliceErrs(request.UserIds, func(from *string) (types.Id, error) {
    return types.IdFromString(*from)
  })

  if ierr != nil {
    return nil, sharedErr.GrpcFieldIndexedErrors("user_id", ierr)
  }

  profiles, stats := p.profileService.Find(ctx, ids)
  response := &proto.FindProfileResponse{
    Profiles: util.CastSlice(profiles, mapper.ToProtoProfile),
  }

  return response, stats.ToGRPCError()
}

func (p *ProfileHandler) Update(ctx context.Context, request *proto.UpdateProfileRequest) (*emptypb.Empty, error) {
  dtoInput := mapper.ToDTOProfileUpdateInput(request)

  err := util.ValidateStruct(ctx, &dtoInput)
  if err != nil {
    return nil, err
  }

  stats := p.profileService.Update(ctx, &dtoInput)
  // PERF: return nil for both success and error response for emptypb.Empty
  if stats.IsError() {
    return nil, stats.ToGRPCError()
  }
  return &emptypb.Empty{}, nil
}

func (p *ProfileHandler) UpdateAvatar(server proto.ProfileService_UpdateAvatarServer) error {
  // NOTE: Doesn't works yet
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

  p.profileService.UpdateAvatar(server.Context(), &dtoInput)
  return server.SendAndClose(&emptypb.Empty{})

}
