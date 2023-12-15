package controllers

import (
	"golfer/config"
	"golfer/services"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type PasswordController struct {
	passwordService services.PasswordService
}

func NewPasswordController(passwordService services.PasswordService) *PasswordController {
	return &PasswordController{passwordService}
}

func (controller PasswordController) SendForgotPasswordEmail(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	if err := controller.passwordService.SendForgotPasswordEmail(request.Email); err != nil {
		if isRecordNotFoundError(err) {
			// We don't want malicious users to know who has an account and who doesn't by default
			c.JSON(http.StatusOK, gin.H{"message": "Successfully sent the e-mail if that account exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "We were unable to send the e-mail"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully sent the e-mail if that account exists"})
}

func (controller PasswordController) ResetPassword(c *gin.Context) {
	var resetPasswordRequest struct {
		Password     string `json:"password"`
		Confirmation string `json:"confirmation"`
		Token        string `json:"token"`
	}

	if err := c.ShouldBindJSON(&resetPasswordRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	if resetPasswordRequest.Password != resetPasswordRequest.Confirmation {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Your passwords don't match"})
		return
	}

	token, err := jwt.ParseWithClaims(resetPasswordRequest.Token, &services.ForgotPasswordClaims{}, func(t *jwt.Token) (interface{}, error) {
		return config.JwtKey, nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unable to parse token"})
		return
	}

	if claims, ok := token.Claims.(*services.ForgotPasswordClaims); ok && token.Valid {
		if err := controller.passwordService.UpdatePassword(claims.Email, resetPasswordRequest.Password, resetPasswordRequest.Confirmation); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully updated your password"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Token is bad or expired"})
	}
}
