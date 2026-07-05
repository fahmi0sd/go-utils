package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Generate(username, secret string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(ttl).Unix(),
		"iat":      time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}
