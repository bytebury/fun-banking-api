package middleware

import (
	"funbanking/internal/infrastructure/auth"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func VerifyUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Token string `json:"token"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(request.Token, &auth.VerificationClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Token expired or invalid"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*auth.VerificationClaims); ok && token.Valid {
			c.Set("verify_email", claims.Recipient)
			c.Next()
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"message": "Token expired or invalid"})
		c.Abort()
	}
}
