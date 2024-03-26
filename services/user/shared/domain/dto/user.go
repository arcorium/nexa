package dto

import (
	"nexa/shared/types"
	"nexa/shared/wrapper"
	"time"
)

type UserResponse struct {
	Id         types.Id `json:"id"`
	Username   string   `json:"username"`
	Email      string   `json:"email"`
	IsVerified bool     `json:"is_verified"`
	FirstName  string   `json:"first_name,omitempty"`
	LastName   string   `json:"last_name,omitempty"`
	Bio        string   `json:"bio,omitempty"`
}

type UserCreateInput struct {
	Username  string                   `json:"username" validate:"required,gte=6"`
	Email     string                   `json:"email" validate:"required,email"`
	Password  string                   `json:"password" validate:"required,gte=6"`
	FirstName string                   `json:"first_name,omitempty" validate:"required"`
	LastName  wrapper.Nullable[string] `json:"last_name,omitempty" validate:""`
	Bio       wrapper.Nullable[string] `json:"bio,omitempty" validate:""`
}

type UserUpdateInput struct {
	Id       string                   `json:"id"`
	Username wrapper.Nullable[string] `json:"username" validate:""`
	Email    wrapper.Nullable[string] `json:"email" validate:""`
}

type UserUpdatePasswordInput struct {
	Id           string `json:"id"`
	LastPassword string `json:"last_password" validate:"required"`
	NewPassword  string `json:"new_password" validate:"required"`
}

type UserBannedInput struct {
	Id       string        `json:"id"`
	Duration time.Duration `json:"duration"`
}

type UserResetPasswordInput struct {
	Id          string `json:"id"`
	NewPassword string `json:"new_password"`
}
