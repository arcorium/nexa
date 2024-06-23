package dto

type SendEmailVerificationDTO struct {
  Recipient string `json:"recipient"`
  Token     string `json:"token"`
}

type SendForgotPasswordDTO struct {
  Recipient string `json:"recipient"`
  Token     string `json:"token"`
}
