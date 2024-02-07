package service

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return UserService{userRepository}
}

func (s UserService) GetCurrentUser(c *gin.Context) {
	var user model.User

	if err := s.userRepository.GetCurrentUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, user)
}
