package mapper

import (
  authZv1 "nexa/proto/gen/go/authorization/v1"
  "nexa/services/authorization/internal/domain/dto"
)

func ToIsAuthorizationDTO(request *authZv1.CheckUserRequest) dto.IsAuthorizationDTO {
  return dto.IsAuthorizationDTO{
    UserId:             request.UserId,
    ExpectedPermission: request.ExpectedPermissions,
  }
}
