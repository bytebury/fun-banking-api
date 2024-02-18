package handlers

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"funbanking/internal/domain/service"
	"funbanking/internal/infrastructure/auth"
	"funbanking/package/utils"
	"net/http"
	"strconv"

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
	userID := c.MustGet("user_id").(string)

	user, err := h.userService.FindByID(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong trying to get user information"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h UserHandler) FindByUsername(c *gin.Context) {
	username := c.Param("username")
	userID := c.GetString("user_id")

	user, err := h.userService.FindByUsernameOrEmail(username)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find that user"})
		return
	}

	if userID == "" || !utils.IsAdmin(userID) || strconv.Itoa(int(user.ID)) != userID {
		user.Email = ""
	}

	c.JSON(http.StatusOK, user)
}

func (h UserHandler) Search(c *gin.Context) {
	// TODO: This will be used by admins to search for users
	// will be paginated
}

func (h UserHandler) FindBanks(c *gin.Context) {
	username := c.Param("username")

	user, err := h.userService.FindByUsernameOrEmail(username)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find user"})
		return
	}

	banks, err := h.userService.FindBanks(strconv.Itoa(int(user.ID)))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to get banks for that user"})
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
