package services

import (
	// "errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	bcrypt "golang.org/x/crypto/bcrypt"
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
	secret := "5!|LjlIV2~WO5xw%rW>?DG8d^4q&96St>3f80b!5UE_D"
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// func ValidateJWTToken(jwtToken string) (err error) {
// 	token, err := jwt.ParseWithClaims(
// 		jwtToken,
// 		&claims{},
// 		func(token *jwt.Token) (interface{}, error) {
// 			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
// 		},
// 	)
// 	if err != nil {
// 		return
// 	}
// 	claims, valid := token.Claims.(*claims)
// 	if !valid {
// 		err = errors.New("invalid token")
// 		return
// 	}
// 	if claims.ExpiresAt < time.Now().Unix() {
// 		err = errors.New("token expired")
// 		return
// 	}
// 	return
// }

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
