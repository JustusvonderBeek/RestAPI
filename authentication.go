package main

import (
	"github.com/golang-jwt/jwt"
)

var sampleSecretKey = []byte("G0qBl4O*ÊLJ0$<©Rî?Gl@ëCR5¢2<3l7pzÃ]M<DõUY:2>0m±o5{CdÑ582&4d«aI'6")

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
