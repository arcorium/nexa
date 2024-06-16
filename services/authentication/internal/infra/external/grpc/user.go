package grpc

import (
  "context"
  "google.golang.org/grpc"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/external"
  userDto "nexa/services/user/shared/domain/dto"
  "nexa/services/user/shared/proto"
  "nexa/shared/types"
)

func NewUserClient(conn grpc.ClientConnInterface) external.IUserClient {
  return &userClient{client: proto.NewUserServiceClient(conn)}
}

type userClient struct {
  client proto.UserServiceClient
}

func (u userClient) ValidateUser(ctx context.Context, email types.Email, password string) (userDto.UserResponseDTO, error) {
  input := proto.ValidateUserInput{
    Email:    email.Underlying(),
    Password: password,
  }
  output, err := u.client.Validate(ctx, &input)
  return userDto.UserResponseDTO{
    Id:         types.IdFromString(output.User.Id),
    Username:   output.User.Username,
    Email:      output.User.Email,
    IsVerified: output.User.IsVerified,
  }, err
}

func (u userClient) RegisterUser(ctx context.Context, request *dto.RegisterDTO) error {
  dtos := proto.CreateUserInput{
    Username:  request.Username,
    Email:     request.Email,
    Password:  request.Password,
    FirstName: request.FirstName,
    LastName:  request.LastName.Value(),
    Bio:       request.Bio.Value(),
  }
  _, err := u.client.Create(ctx, &dtos)
  return err
}
