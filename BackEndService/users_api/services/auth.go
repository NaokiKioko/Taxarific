package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	bcrypt "golang.org/x/crypto/bcrypt"
	"github.com/joho/godotenv"
)

type claims struct {
	Id   string `json:"username"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(id string, role string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &claims{
		Id:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	err := godotenv.Load(".env")
	if err != nil {
		return "", err
	}
	secret := os.Getenv("JWT_SECRET_KEY")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWTToken(jwtToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		jwtToken,
		&claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		},
	)
	if err != nil {
		return
	}
	claims, valid := token.Claims.(*claims)
	if !valid {
		err = errors.New("invalid token")
		return
	}
	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		err = errors.New("token expired")
		return
	}
	return
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
func CheckPassword(providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(providedPassword))
}
