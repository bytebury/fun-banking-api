package middleware

import (
	"funbanking/internal/infrastructure/auth"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Customer() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &auth.CustomerClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			Auth()(c)
			return
		}

		if claims, ok := token.Claims.(*auth.CustomerClaims); ok && token.Valid {
			if claims.CustomerID != "" {
				c.Set("customer_id", claims.CustomerID)
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to do this action"})
		c.Abort()
	}
}
