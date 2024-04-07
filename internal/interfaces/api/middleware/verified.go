package middleware

import (
	"funbanking/internal/domain/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Verified() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)
		userRepository := users.NewUserRepository()

		var user users.User
		if err := userRepository.FindByID(userID, &user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Could not find user"})
			c.Abort()
			return
		}

		if !user.Verified {
			c.JSON(http.StatusForbidden, gin.H{"message": "You need to verify your account first"})
			c.Abort()
			return
		}

		c.Next()
	}
}
