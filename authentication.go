package main

import (
	"os"

	"github.com/golang-jwt/jwt"
)

func generateToken() (string, error) {
	// Generate a new token for the client
	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			"userId":     "lerner",
			"authorized": true,
		})

	// Read the secret from the environment
	secretKey := os.Getenv(SECRET_KEY)
	secretKeyByte := []byte(secretKey)
	// Sign the token with a valid secret
	tokenString, err := token.SignedString(secretKeyByte)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
