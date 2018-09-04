package main

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
)

var secret = []byte("password")

const tokenExpiration = time.Second * time.Duration(20)

type AityClaims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

func TokenForUser(email string) (string, error) {
	expTime := time.Now().Add(tokenExpiration)
	claims := AityClaims{
		User: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
			Issuer:    "AITY Backend",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ValidateToken checks the header if a token is present. If it is present it checks it, whether it is a valid token. This function returns not only the valid state but also recommends an http status code.
func ValidateToken(header http.Header) (bool, int, *jwt.Token) {
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return false, http.StatusUnauthorized, &jwt.Token{}
	}
	splitted := strings.Split(authHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	if len(splitted) != 2 {
		return false, http.StatusUnauthorized, &jwt.Token{}
	}

	tokenString := splitted[1]

	token, err := jwt.ParseWithClaims(tokenString, &AityClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return false, http.StatusUnauthorized, &jwt.Token{}
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return false, http.StatusUnauthorized, &jwt.Token{}
		}
	} else if token.Valid {
		return true, http.StatusOK, token
	}
	return false, http.StatusForbidden, &jwt.Token{}
}
