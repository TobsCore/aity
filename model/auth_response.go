package model

type AuthResponse struct {
	AlreadyRegistered bool   `json:"already_registered"`
	AuthToken         string `json:"auth_token"`
	UserInfo          User   `json:"user"`
}
