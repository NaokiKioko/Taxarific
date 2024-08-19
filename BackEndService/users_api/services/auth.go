package services

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type claims struct {
	Id   string `json:"username"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(id string, role string) (string, error) {
	experationTime := time.Now().Add(time.Hour)
	claims := &claims{
		Id:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(experationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET_KEY")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
