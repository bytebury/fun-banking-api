package middleware

import (
	"golfer/config"
	"golfer/models"
	"golfer/repositories"
	"golfer/services"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &services.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return config.JwtKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*services.UserClaims); ok && token.Valid {
			c.Set("user_id", claims.UserID)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
		}
	}
}

/**
 * Auditing will check to see if a user exists, but will not fail it if it is not present.
 * this is used primarily for tracking if possible. Useful for routes that might want to see
 * if a user is doing an action; but technically the user might not be required.
 */
func Audit() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.Next()
		}

		token, err := jwt.ParseWithClaims(tokenString, &services.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return config.JwtKey, nil
		})

		if err != nil {
			c.Next()
		}

		if claims, ok := token.Claims.(*services.UserClaims); ok && token.Valid {
			c.Set("user_id", claims.UserID)
			c.Next()
		} else {
			c.Next()
		}
	}
}

func Admin(repository repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &services.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return config.JwtKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*services.UserClaims); ok && token.Valid {
			c.Set("user_id", claims.UserID)

			var user models.User
			if err := repository.FindByID(claims.UserID, &user); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
				c.Abort()
				return
			}

			if user.Role != config.AdminRole {
				c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to do that"})
				c.Abort()
				return
			}

			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
		}
	}
}
