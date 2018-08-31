package model

type AuthResponse struct {
	AlreadyRegistered bool       `json:"already_registered"`
	UserInfo          User `json:"user"`
}