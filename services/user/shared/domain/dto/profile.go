package dto

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
	UserId    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Bio       string `json:"bio"`
}

// TODO: Implement Upload file which is received from api gateway
type ProfilePictureUpdateInput struct {
	UserId   string `json:"user_id"`
	Filename string `json:"filename"`
}
