package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"time"
)

var secret = []byte("password")

const tokenExpiration = time.Hour * time.Duration(24)

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

// validateToken checks the header if a token is present. If it is present it checks it, whether it is a valid token. This function returns not only the valid state but also recommends an http status code.
func validateToken(header http.Header) (bool, int, *jwt.Token) {
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

func ValidateRequest(r *http.Request) (bool, string, int) {
	user := mux.Vars(r)["user"]
	var validToken, suggestedStatus, token = validateToken(r.Header)
	if !validToken {
		return false, "Cannot validate token", suggestedStatus
	}

	// Check, whether the user checks out. In order to do that, receive the user of the token and then check it agains the user in the url. To ensure everything works as expected, usernames are cast to lower case.
	claims, ok := token.Claims.(*AityClaims)
	if !ok || !token.Valid {
		return false, "Cannot parse claims", http.StatusBadRequest
	}
	claimsUser := claims.User
	if strings.ToLower(claimsUser) != strings.ToLower(user) {
		log.Printf("%s - %s", claimsUser, strings.ToLower(user))
		return false, "Not authorized", http.StatusUnauthorized
	}

	// If it is valid, return true and the other values don't matter
	return true, "", 0
}
