package controllers

import (
	"golfer/models"
	"golfer/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SessionsController struct {
	userService services.UserService
}

func NewSessionsController(userService services.UserService) *SessionsController {
	return &SessionsController{
		userService,
	}
}

func (controller *SessionsController) Login(c *gin.Context) {
	var request models.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	token, user, err := controller.userService.Login(request)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred, please try again."})
		return
	}

	if token == "" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}
