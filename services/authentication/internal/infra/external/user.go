package external

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  userv1 "nexa/proto/gen/go/user/v1"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/external"
  "nexa/services/authentication/util"
  spanUtil "nexa/shared/span"
  "nexa/shared/types"
)

func NewUserClient(conn grpc.ClientConnInterface) external.IUserClient {
  return &userClient{
    client: userv1.NewUserServiceClient(conn),
    trace:  util.GetTracer(),
  }
}

type userClient struct {
  client userv1.UserServiceClient

  trace trace.Tracer
}

func (u *userClient) Validate(ctx context.Context, email types.Email, password string) (dto.UserValidateResponseDTO, error) {
  ctx, span := u.trace.Start(ctx, "UserClient.Validate")
  defer span.End()

  // Call
  dtos := userv1.ValidateUserRequest{
    Email:    email.String(),
    Password: password,
  }
  response, err := u.client.Validate(ctx, &dtos)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.UserValidateResponseDTO{}, err
  }

  // Map to response DTO
  user := response.User
  id, err := types.IdFromString(user.Id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.UserValidateResponseDTO{}, err
  }

  responseDTO := dto.UserValidateResponseDTO{
    UserId:   id,
    Username: user.Username,
  }

  return responseDTO, nil
}

func (u *userClient) Create(ctx context.Context, request *dto.RegisterDTO) error {
  ctx, span := u.trace.Start(ctx, "UserClient.Create")
  defer span.End()

  dtos := userv1.CreateUserRequest{
    Username:  request.Username,
    Email:     request.Email,
    Password:  request.Password,
    FirstName: request.FirstName,
    LastName:  request.LastName.Value(),
    Bio:       request.Bio.Value(),
  }

  _, err := u.client.Create(ctx, &dtos)
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }
  return nil
}
