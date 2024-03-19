package dto

import "nexa/shared/types"

type UserResponse struct {
	Id       types.Id `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Role     string   `json:"role"`
}
