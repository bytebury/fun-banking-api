package service

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	GetCurrentUser(c *gin.Context)
	FindByID(c *gin.Context)
	FindByUsernameOrEmail(c *gin.Context)
	Update(c *gin.Context)
	Create(c *gin.Context)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return userService{userRepository}
}

func (s userService) GetCurrentUser(c *gin.Context) {
	var user model.User

	if err := s.userRepository.GetCurrentUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s userService) FindByID(c *gin.Context) {
	var user model.User
	userID := c.Param("id")

	if err := s.userRepository.FindByID(userID, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s userService) FindByUsernameOrEmail(c *gin.Context) {
	var user model.User
	var request struct {
		username string
		email    string
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if strings.TrimSpace(request.username) == "" {
		request.username = request.email
	}

	if err := s.userRepository.FindByUsernameOrEmail(request.username, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s userService) Update(c *gin.Context) {
	var user model.User
	userID := c.Param("id")

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := s.userRepository.Update(userID, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusAccepted, user)
}

func (s userService) Create(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := s.userRepository.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, user)
}
