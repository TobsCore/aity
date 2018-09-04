package model

type AuthResponse struct {
	AlreadyRegistered  bool    `json:"already_registered"`
	AuthToken          string  `json:"auth_token"`
	ExpirationDuration float64 `json:"expires_in"`
	UserInfo           User    `json:"user"`
}
