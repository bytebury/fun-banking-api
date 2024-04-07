package auth

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type CustomerClaims struct {
	CustomerID string `json:"customer_id"`
	jwt.StandardClaims
}

type ForgotPasswordClaims struct {
	Recipient string `json:"recipient"`
	jwt.StandardClaims
}

type VerificationClaims struct {
	Recipient string `json:"recipient"`
	jwt.StandardClaims
}

type JWTService struct{}

func NewJWTService() JWTService {
	return JWTService{}
}

func (j JWTService) GenerateUserToken(userID string) (string, error) {
	claims := UserClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * 365 * 100 * time.Hour).Unix(), // One-hundred years from today
			IssuedAt:  time.Now().Unix(),
		},
	}
	return generateToken(claims)
}

func (j JWTService) GenerateCustomerToken(customerID string) (string, error) {
	claims := CustomerClaims{
		CustomerID: customerID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	return generateToken(claims)
}

func (j JWTService) GenerateForgotPasswordToken(recipient string) (string, error) {
	claims := ForgotPasswordClaims{
		Recipient: recipient,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	return generateToken(claims)
}

func (j JWTService) GenerateVerificationToken(recipient string) (string, error) {
	claims := VerificationClaims{
		Recipient: recipient,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	return generateToken(claims)
}

func generateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
