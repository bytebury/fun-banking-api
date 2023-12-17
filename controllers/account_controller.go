package controllers

import (
	"golfer/models"
	"golfer/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	service services.AccountService
}

func NewAccountController(account services.AccountService) *AccountController {
	return &AccountController{account}
}

func (controller AccountController) FindByID(c *gin.Context) {
	accountID := c.Param("id")
	var account models.Account
	err := controller.service.FindByID(accountID, &account)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, account)
}
