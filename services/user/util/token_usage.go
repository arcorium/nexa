package util

import (
  authNv1 "nexa/proto/gen/go/authentication/v1"
  "nexa/services/user/internal/domain/dto"
)

func TokenPurposeToUsage(purpose dto.TokenPurpose) authNv1.TokenUsage {
  if purpose == dto.EmailVerificationToken {
    return authNv1.TokenUsage_VERIFICATION
  } else if purpose == dto.ForgotPasswordToken {
    return authNv1.TokenUsage_RESET_PASSWORD
  }
  return authNv1.TokenUsage(purpose)
}

func TokenUsageToPurpose(usage authNv1.TokenUsage) dto.TokenPurpose {
  if usage == authNv1.TokenUsage_VERIFICATION {
    return dto.EmailVerificationToken
  } else if usage == authNv1.TokenUsage_RESET_PASSWORD {
    return dto.ForgotPasswordToken
  }
  return dto.TokenPurpose(usage)
}
