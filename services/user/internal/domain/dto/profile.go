package dto

import (
	"nexa/shared/wrapper"
)

type ProfileResponse struct {
	UserId    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Bio       string `json:"bio"`
	PhotoURL  string `json:"photo_url"`
}

//type ProfileCreateInput struct {
//	UserId    string `json:"user_id"`
//	FirstName string `json:"first_name"`
//	LastName  string `json:"last_name"`
//	Bio       string `json:"bio"`
//}

type ProfileUpdateInput struct {
	UserId    string                   `json:"user_id" validate:"uuid4"`
	FirstName wrapper.Nullable[string] `json:"first_name" validate:""`
	LastName  wrapper.Nullable[string] `json:"last_name" validate:""`
	Bio       wrapper.Nullable[string] `json:"bio" validate:""`
}

// TODO: Implement Upload file which is received from api gateway
type ProfilePictureUpdateInput struct {
	UserId   string `json:"user_id" validate:"uuid4"`
	Filename string `json:"filename" validate:"required"`
	Bytes    []byte `validate:"required,image"`
}
