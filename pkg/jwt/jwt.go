package jwt

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	SecretKey = []byte("secret")
)

// Generates a jwt token and assign a username to its claims and return it
func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	/* Crete a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	// Set token claims
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal("Error in generating key")
		return "", err
	}
	return tokenString, nil
}

// Parses a jwt token and returns the username in its claims
func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	} else {
		return "", err
	}
}
