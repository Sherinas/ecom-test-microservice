package jwtauth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(email, secret, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(secret))
}
