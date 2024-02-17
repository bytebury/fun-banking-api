package middleware

import (
	"funbanking/internal/infrastructure/auth"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		tokenStrings := strings.Split(tokenString, " ")
		tokenString = tokenStrings[len(tokenStrings)-1]

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &auth.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to do this action"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*auth.UserClaims); ok && token.Valid {
			c.Set("user_id", claims.UserID)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to do this action"})
			c.Abort()
		}
	}
}
