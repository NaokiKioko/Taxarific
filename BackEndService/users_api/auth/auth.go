package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	bcrypt "golang.org/x/crypto/bcrypt"
)

type claims struct {
	UserId string
	Role   string
	jwt.RegisteredClaims
}

func GenerateJWTToken(id string, role string) (string, error) {
	claims := &claims{
		UserId:           id,
		Role:             role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
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

func ValidateJWTToken(jwtToken string) (*claims, error) {
	token, err := jwt.ParseWithClaims(
		jwtToken,
		&claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, valid := token.Claims.(*claims)
	if !valid {
		err = errors.New("invalid token")
		return nil, err
	}
	if !claims.ExpiresAt.Time.IsZero() && claims.ExpiresAt.Time.Before(time.Now()) {
		err = errors.New("token expired")
		return nil, err
	}
	return claims, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
func CheckPassword(providedPassword string, storedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(storedPassword))
}
