package middleware

import (
	"funbanking/internal/infrastructure/auth"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func PasswordReset() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Token                string `json:"token"`
			Password             string `json:"password"`
			PasswordConfirmation string `json:"password_confirmation"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(request.Token, &auth.ForgotPasswordClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Token expired or invalid"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*auth.ForgotPasswordClaims); ok && token.Valid {
			c.Set("reset_email", claims.Recipient)
			c.Set("reset_password", request.Password)
			c.Set("reset_confirmation", request.PasswordConfirmation)
			c.Next()
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"message": "Token expired or invalid"})
		c.Abort()
	}
}
