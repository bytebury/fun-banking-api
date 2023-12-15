package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // You can set this to a specific origin or origins
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	return cors.New(config)
}
