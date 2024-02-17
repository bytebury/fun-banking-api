package handlers

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"funbanking/internal/domain/service"
	"funbanking/internal/infrastructure/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler() UserHandler {
	return UserHandler{
		userService: service.NewUserService(
			repository.NewUserRepository(),
		),
	}
}

func (h UserHandler) GetCurrentUser(c *gin.Context) {
	// TODO: need to wait for middleware to be done to do this one
	h.userService.FindByID("0")
}

func (h UserHandler) FindByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.FindByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h UserHandler) FindBanks(c *gin.Context) {
	id := c.Param("id")

	banks, err := h.userService.FindBanks(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find user"})
		return
	}

	c.JSON(http.StatusOK, banks)
}

func (h UserHandler) Update(c *gin.Context) {
	var user model.User
	id := c.Param("id")

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := h.userService.Update(id, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h UserHandler) Create(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := h.userService.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h UserHandler) Login(c *gin.Context) {
	var request auth.UserLoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	token, user, err := h.userService.Login(request.UsernameOrEmail, request.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unable to log you in, invalid credentials"})
		return
	}

	response := struct {
		Token string     `json:"token"`
		User  model.User `json:"user"`
	}{Token: token, User: user}

	c.JSON(http.StatusOK, response)

}
