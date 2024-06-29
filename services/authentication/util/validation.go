package util

import (
  "fmt"
  "github.com/go-playground/validator/v10"
  "nexa/services/authentication/internal/domain/entity"
)

func RegisterValidationTags(validate *validator.Validate) {
  validate.RegisterAlias("usage", fmt.Sprintf("oneof=%s %s", entity.UsageVerifStr, entity.UsageResetStr))
  validate.RegisterAlias("usage_enum", fmt.Sprintf("lt=%d", entity.TokenUsageUnknown.Underlying()))
}
