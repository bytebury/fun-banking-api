package middleware

import (
	"funbanking/internal/infrastructure/auth"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Audit() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		tokenStrings := strings.Split(tokenString, " ")
		tokenString = tokenStrings[len(tokenStrings)-1]

		if tokenString == "" {
			c.Next()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &auth.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.Next()
			return
		}

		if claims, ok := token.Claims.(*auth.UserClaims); ok && token.Valid {
			c.Set("user_id", claims.UserID)
			c.Next()
			return
		}

		c.Next()
	}
}
