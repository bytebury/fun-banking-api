package services

import (
	"golfer/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtService struct{}

type UserClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type CustomerClaims struct {
	CustomerID string `json:"customer_id"`
	jwt.StandardClaims
}

type ForgotPasswordClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func (service JwtService) GenerateUserToken(userID string) (string, error) {
	claims := &UserClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * 365 * 100 * time.Hour).Unix(), // One-hundred years from today
			IssuedAt:  time.Now().Unix(),
		},
	}
	return service.generateToken(claims)
}

func (service JwtService) GenerateCustomerToken(customerID string) (string, error) {
	claims := &CustomerClaims{
		CustomerID: customerID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * 365 * 100 * time.Hour).Unix(), // One-hundred years from today
			IssuedAt:  time.Now().Unix(),
		},
	}
	return service.generateToken(claims)
}

func (service JwtService) GeneratePasswordResetToken(email string) (string, error) {
	claims := &ForgotPasswordClaims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	return service.generateToken(claims)
}

func (service JwtService) generateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
