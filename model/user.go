package model

import "time"

type User struct {
	Email          string
	Username       string
	RegisteredDate time.Time
}

type UserInfo struct {
	Email    string
	Username string
}
