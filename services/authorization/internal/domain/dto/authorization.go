package dto

type IsAuthorizationDTO struct {
  UserId             string   `validate:"required,uuid4"`
  ExpectedPermission []string `validate:"required"`
}
