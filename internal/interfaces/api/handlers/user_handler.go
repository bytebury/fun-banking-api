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

func (h UserHandler) FindByID(c *gin.Context) {
	h.user.FindByID(c)
}

func (h UserHandler) FindByUsernameOrEmail(c *gin.Context) {
	h.user.FindByUsernameOrEmail(c)
}

func (h UserHandler) Update(c *gin.Context) {
	h.user.Update(c)
}

func (h UserHandler) Create(c *gin.Context) {
	h.user.Create(c)
}
