package auth

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type JWTResponse struct {
	Token string `json:"token"`
}

type JWTCustomClaims struct {
	UID   uint   `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(uid uint, email string) (string, error) {
	claims := &JWTCustomClaims{
		uid,
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
