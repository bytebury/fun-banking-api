package controllers

import (
	"golfer/models"
	"golfer/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	user services.UserService
}

func NewUserController(user services.UserService) *UserController {
	return &UserController{
		user: user,
	}
}

func (controller UserController) FindCurrentUser(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	var user models.User
	err := controller.user.FindByID(userID, &user)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (handler UserController) FindByID(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := handler.user.FindByID(userID, &user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (handler UserController) FindByUsername(c *gin.Context) {
	username := c.Param("username")

	var user models.User
	if err := handler.user.FindByUsername(username, &user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (handler UserController) Update(c *gin.Context) {
	var request models.UserRequest
	userID := c.Param("id")

	if !handler.canModify(c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	user, err := handler.user.Update(userID, &request)

	if isDuplicateError(err) {
		if strings.Contains(err.Error(), "username") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "That username already exists"})
			return
		}

		if strings.Contains(err.Error(), "email") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "That email is already associated to another account"})
			return
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong updating the user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (handler UserController) Create(c *gin.Context) {
	var request models.UserRequest
	var user models.User

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	err := handler.user.Create(&request, &user)

	if isPasswordsDoNotMatchError(err) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Those passwords do not match"})
		return
	}

	if isDuplicateError(err) {
		if strings.Contains(err.Error(), "username") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "That username already exists"})
			return
		}

		if strings.Contains(err.Error(), "email") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "That email is already associated to another account"})
			return
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong creating your account"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (handler UserController) Delete(c *gin.Context) {
	userID := c.Param("id")

	if !handler.canModify(c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that"})
		return
	}

	if err := handler.user.Delete(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong deleting the user"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

/**
 * Determines if the current user can modify the given user in context.
 * This would be useful when setting up administrative access. For right now,
 * The logic just makes sure they are who the ID is. For example, you can only
 * update or delete your own account.
 */
func (handler UserController) canModify(c *gin.Context) bool {
	userID := c.Param("id")
	currentUserID := c.MustGet("user_id").(string)
	return userID == currentUserID
}
