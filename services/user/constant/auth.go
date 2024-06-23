package constant

import "time"

const TOKEN_METADATA_KEY = "bearer"

// Permissions

// Resource
const USER_RESOURCE = "users"

// Actions
const (
  BANNED_ACTION = "banned"
)

const (
  EMAIL_VERIFICAITON_TOKEN_TTL time.Duration = time.Hour * 24
  FORGOT_PASSWORD_TOKEN_TTL                  = time.Hour * 24
)
