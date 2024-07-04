package external

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  userv1 "nexa/proto/gen/go/user/v1"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/external"
  "nexa/services/authentication/util"
  "nexa/shared/types"
  spanUtil "nexa/shared/util/span"
)

func NewUserClient(conn grpc.ClientConnInterface) external.IUserClient {
  return &userClient{
    client: userv1.NewUserServiceClient(conn),
    trace:  util.GetTracer(),
  }
}

type userClient struct {
  external.IUserClient
  client userv1.UserServiceClient

  trace trace.Tracer
}

func (u *userClient) Validate(ctx context.Context, email types.Email, password types.Password) (dto.UserResponseDTO, error) {
  ctx, span := u.trace.Start(ctx, "UserClient.Validate")
  defer span.End()

  // Call
  dtos := userv1.ValidateUserRequest{
    Email:    email.String(),
    Password: password.String(),
  }

  response, err := u.client.Validate(ctx, &dtos)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.UserResponseDTO{}, err
  }

  // Map to response DTO
  user := response.User
  id, err := types.IdFromString(user.Id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return dto.UserResponseDTO{}, err
  }

  responseDTO := dto.UserResponseDTO{
    UserId:   id,
    Username: user.Username,
  }

  return responseDTO, nil
}

func (u *userClient) Create(ctx context.Context, request *dto.RegisterDTO) (types.Id, error) {
  ctx, span := u.trace.Start(ctx, "UserClient.Create")
  defer span.End()

  dtos := userv1.CreateUserRequest{
    Username:  request.Username,
    Email:     request.Email.String(),
    Password:  request.Password.String(),
    FirstName: request.FirstName,
    LastName:  request.LastName.Value(),
    Bio:       request.Bio.Value(),
  }

  res, err := u.client.Create(ctx, &dtos)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), err
  }

  userId, err := types.IdFromString(res.Id)
  if err != nil {
    spanUtil.RecordError(err, span)
    return types.NullId(), err
  }

  return userId, nil
}
