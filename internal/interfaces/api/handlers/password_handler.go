package handlers

import (
	"funbanking/internal/domain/users"
	"funbanking/internal/infrastructure/mailing"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type PasswordHandler struct {
	passwordService users.PasswordService
}

func NewPasswordHandler() PasswordHandler {
	return PasswordHandler{
		passwordService: users.NewPasswordService(
			mailing.NewForgotPasswordMailer(),
		),
	}
}

func (h PasswordHandler) ForgotPassword(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := h.passwordService.ForgotPassword(request.Email); err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			c.JSON(http.StatusOK, gin.H{"message": "We sent you an e-mail if your account exists with directions on how to reset your password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong sending the e-mail"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "We sent you an e-mail if your account exists with directions on how to reset your password"})
}

func (h PasswordHandler) ResetPassword(c *gin.Context) {
	email := c.MustGet("reset_email").(string)
	password := c.GetString("reset_password")
	passwordConfirmation := c.GetString("reset_confirmation")

	if err := h.passwordService.ResetPassword(email, password, passwordConfirmation); err != nil {
		if strings.Contains(err.Error(), "passwords do not match") {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Passwords do not match"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong updating your password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "You have successfully reset your password"})
}
