package api

import (
	"funbanking/internal/interfaces/api/handlers"

	"github.com/gin-gonic/gin"
)

func Run() {
	r := runner{router: gin.Default()}
	r.setup()
	r.router.Run()
}

type runner struct {
	router *gin.Engine
}

func (r runner) setup() {
	r.setupHealthRoutes()
	r.setupUserRoutes()
}

func (r runner) setupHealthRoutes() {
	handler := handlers.NewHealthHandler()
	r.router.Group("/health").GET("/", handler.GetHealthCheck)
}

func (r runner) setupUserRoutes() {
	handler := handlers.NewUserHandler()
	r.router.GET("/current-user", handler.GetCurrentUser)
}
