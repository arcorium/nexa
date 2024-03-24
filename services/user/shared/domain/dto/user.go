package dto

import (
	"nexa/shared/types"
	"time"
)

type UserResponse struct {
	Id        types.Id `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	FirstName string   `json:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty"`
	Bio       string   `json:"bio,omitempty"`
}

type UserCreateInput struct {
	Username  string `json:"username" validate:"required,gte=6"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,gte=6"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Bio       string `json:"bio,omitempty"`
}

type UserUpdateInput struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserUpdatePasswordInput struct {
	Id           string `json:"id"`
	LastPassword string `json:"last_password"`
	NewPassword  string `json:"new_password"`
}

type UserBannedInput struct {
	Id       string        `json:"id"`
	Duration time.Duration `json:"duration"`
}

type UserResetPasswordInput struct {
	Id          string `json:"id"`
	NewPassword string `json:"new_password"`
}
