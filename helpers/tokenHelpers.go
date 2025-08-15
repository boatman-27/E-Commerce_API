package helpers

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateAccessToken(userid, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"UserId": userid,
		"Email":  email,
		"Role":   role,
		"exp":    time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
}

func GenerateRefreshToken(userid, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"UserId": userid,
		"Email":  email,
		"Role":   role,
		"exp":    time.Now().Add(7 * 24 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
}
