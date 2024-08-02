package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("your-secret-key")

func JwtCreateToken(username string, role string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"iss": "mchosting",
		"aud": role,
		"exp": time.Now().Add(time.Minute * 30).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
