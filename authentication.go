package main

import (
	"github.com/golang-jwt/jwt"
)

var sampleSecretKey = []byte("0>i3(r61D11¤Tqu5$kz$4ãMb(1ð>rOëpoP7=o§æ[16#Mt?çoe0206;s4)KÁD3<<o")

func generateToken() (string, error) {
	// Generate a new token for the client
	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.MapClaims{
			"userId":     "lerner",
			"authorized": true,
		})

	// Sign the token with a valid secret
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
