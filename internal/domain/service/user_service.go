package service

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserService struct {
	repo repository.User
}

func NewUserService(userRepository repository.User) UserService {
	return UserService{repo: userRepository}
}

func (s UserService) GetCurrentUser(c *gin.Context) {
	var user model.User
	err := s.repo.GetCurrentUser(&user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, user)
}
