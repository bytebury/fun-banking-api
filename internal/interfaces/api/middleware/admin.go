package middleware

import (
	"funbanking/internal/infrastructure/auth"
	"funbanking/package/utils"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Admin() gin.HandlerFunc {
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

			if utils.IsAdmin(claims.UserID) {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to do this action"})
		c.Abort()
	}
}
