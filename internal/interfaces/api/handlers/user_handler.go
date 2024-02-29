package handlers

import (
	"funbanking/internal/domain/users"
	"funbanking/internal/infrastructure/auth"
	"funbanking/internal/infrastructure/mailing"
	"funbanking/package/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService users.UserService
}

func NewUserHandler() UserHandler {
	userRepository := users.NewUserRepository()
	return UserHandler{
		userService: users.NewUserService(
			userRepository,
			auth.NewUserAuth(userRepository),
			mailing.NewWelcomeMailer(),
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

	h.userService.AddVisitor(&users.Visitor{
		UserID:     &user.ID,
		CustomerID: nil,
		IPAddress:  c.ClientIP(),
	})

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
	itemsPerPage, _ := strconv.Atoi(c.Query("itemsPerPage"))
	pageNumber, _ := strconv.Atoi(c.Query("pageNumber"))

	params := map[string]string{
		"ID":       c.Query("id"),
		"Username": c.Query("username"),
		"Email":    c.Query("email"),
	}

	users, err := h.userService.FindAll(itemsPerPage, pageNumber, params)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h UserHandler) Update(c *gin.Context) {
	var user users.User

	currentUserID := c.MustGet("user_id").(string)
	userID := c.Param("id")

	if currentUserID != userID && !utils.IsAdmin(currentUserID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to do that"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := h.userService.Update(userID, &user); err != nil {
		if strings.Contains(err.Error(), "users_username_key") {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "That username is already in use"})
			return
		}

		if strings.Contains(err.Error(), "users_email_key") {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "That email is already in use"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h UserHandler) Create(c *gin.Context) {
	var request users.NewUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	user, err := h.userService.Create(&request)

	if err != nil {
		if strings.Contains(err.Error(), "users_username_key") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "That username is already in use"})
			return
		}

		if strings.Contains(err.Error(), "users_email_key") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "That email is already in use"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h UserHandler) Login(c *gin.Context) {
	var request users.LoginRequest

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
		User  users.User `json:"user"`
	}{Token: token, User: user}

	c.JSON(http.StatusOK, response)

}
