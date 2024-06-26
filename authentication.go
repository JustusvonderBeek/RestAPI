package main

import (
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
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
	// log.Printf("Env key: %s", secretKey)
	secretKeyByte := []byte(secretKey)
	// Sign the token with a valid secret
	tokenString, err := token.SignedString(secretKeyByte)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func IPWhiteList(whitelist map[string]bool) gin.HandlerFunc {
	f := func(c *gin.Context) {
		// If the IP isn't in the whitelist, forbid the request.
		ip := c.ClientIP()
		// log.Printf("IP: %s", ip)

		re := regexp.MustCompile(`([\d]+).([\d]+).[\d]+.[\d]+`)
		ipRange := re.ReplaceAllString(ip, "$1.$2.0.0")
		// log.Printf("Ip Range: %s", ipRange)
		if !whitelist[ip] && !whitelist[ipRange] {
			log.Printf("Unauthorized access from %s", ip)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "fuck you in the ass blyad"})
			return
		}
		c.Next()
	}
	return f
}
