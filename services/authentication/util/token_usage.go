package util

import (
  tokenv1 "github.com/arcorium/nexa/proto/gen/go/token/v1"
  "nexa/services/authentication/internal/domain/dto"
)

func TokenPurposeToUsage(purpose dto.TokenUsage) tokenv1.TokenUsage {
  if purpose == dto.TokenUsageEmailVerification {
    return tokenv1.TokenUsage_EMAIL_VERIFICATION
  } else if purpose == dto.TokenUsageResetPassword {
    return tokenv1.TokenUsage_RESET_PASSWORD
  }
  return tokenv1.TokenUsage(purpose)
}

func TokenUsageToPurpose(usage tokenv1.TokenUsage) dto.TokenUsage {
  if usage == tokenv1.TokenUsage_EMAIL_VERIFICATION {
    return dto.TokenUsageEmailVerification
  } else if usage == tokenv1.TokenUsage_RESET_PASSWORD {
    return dto.TokenUsageResetPassword
  }
  return dto.TokenUsage(usage)
}
