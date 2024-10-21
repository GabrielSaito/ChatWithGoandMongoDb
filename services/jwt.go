package services

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("secret")

func GenerateToken(username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return ""
	}
	return tokenString
}
