package handlers

import (
	"funbanking/internal/domain/repository"
	"funbanking/internal/domain/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	user service.UserService
}

func NewUserHandler() UserHandler {
	return UserHandler{
		user: service.NewUserService(repository.NewUserRepository()),
	}
}

func (h UserHandler) GetCurrentUser(c *gin.Context) {
	h.user.GetCurrentUser(c)
}
