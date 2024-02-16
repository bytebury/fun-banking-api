package auth

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTService struct{}

func NewJWTService() JWTService {
	return JWTService{}
}

func (j JWTService) GenerateUserToken(userID string) (string, error) {
	claims := struct {
		UserID string `json:"user_id"`
		jwt.StandardClaims
	}{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * 365 * 100 * time.Hour).Unix(), // One-hundred years from today
			IssuedAt:  time.Now().Unix(),
		},
	}
	return generateToken(claims)
}

func generateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(os.Getenv("JWT_SECRET"))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
