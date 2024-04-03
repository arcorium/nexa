package handler

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"nexa/services/user/internal/api/grpc/mapper"
	"nexa/services/user/internal/domain/dto"
	"nexa/services/user/internal/domain/service"
	"nexa/services/user/shared/proto"
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

func (p *ProfileHandler) Find(request *proto.FindProfileRequest, server proto.ProfileService_FindServer) error {
	ids, err := util.CastSliceErr(request.UserIds, func(from *string) (types.Id, error) {
		id := types.IdFromString(*from)
		return id, id.Validate()
	})
	if err != nil {
		return err
	}

	profiles, stats := p.profileService.Find(server.Context(), ids)
	if stats.IsError() {
		return stats.Error
	}

	for _, profile := range profiles {
		response := mapper.ToProtoProfile(&profile)
		if err := server.Send(&response); err != nil {
			return err
		}
	}
	return nil
}

func (p *ProfileHandler) Update(ctx context.Context, request *proto.UpdateProfileRequest) (*emptypb.Empty, error) {
	dtoInput := mapper.ToDTOProfileUpdateInput(request)

	if err := util.GetValidator().Struct(&dtoInput); err != nil {
		return nil, err
	}

	stats := p.profileService.Update(ctx, &dtoInput)
	if stats.IsError() {
		return nil, stats.Error
	}
	return &emptypb.Empty{}, nil
}

func (p *ProfileHandler) UpdateAvatar(server proto.ProfileService_UpdateAvatarServer) error {
	dtoInput := dto.ProfilePictureUpdateInput{}
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
