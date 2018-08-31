package main

import (
	"encoding/json"
	"errors"
	"github.com/tobscore/aity/model"
	"log"
	"net/http"
	"time"
)

type authResponse struct {
	alreadyRegistered bool           `json:"alreadyRegistered"`
	userinfo          model.UserInfo `json:"user"`
}

type GoogleUser struct {
	Id           string `json:"id"`
	Email        string `json:"email"`
	VerifiedMail bool   `json:"verified_email"`
	Name         string `json:"name"`
	GivenName    string `json:"given_name"`
	FamilyName   string `json:"family_name"`
	Link         string `json:"link"`
	PictureUrl   string `json:"picture"`
	Gender       string `json:"gender"`
	Locale       string `json:"locale"`
}

func (u *GoogleUser) toUser() *model.User {
	return &model.User{
		Email:          u.Email,
		Username:       u.Name,
		RegisteredDate: time.Now(),
	}
}

func (u *GoogleUser) toUserInfo(username string) *model.UserInfo {
	var uInfo model.UserInfo
	uInfo.Email = u.Email
	uInfo.Username = username
	return &uInfo
}

func lookupGoogleUser(token string) (*GoogleUser, error, int) {
	// Receive the user's information from google's servers in order to check if the user already exists in the database
	// or create it otherwise
	client := &http.Client{}
	var uInfo GoogleUser
	req, err := http.NewRequest("GET", "https://www.googleapis.com/userinfo/v2/me", nil)
	if err != nil {
		return &uInfo, err, http.StatusBadRequest
	}
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Cannot contact google server for authorization.")
		log.Println(err.Error())
		return &uInfo, err, http.StatusBadRequest
	} else if resp.StatusCode != 200 {
		return &uInfo, errors.New("invalid token"), http.StatusUnauthorized
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&uInfo)
	if err != nil {
		log.Println(err.Error())
		return &uInfo, err, http.StatusBadRequest
	}
	return &uInfo, nil, http.StatusOK
}
