package controllers

import (
	"golfer/models"
	"golfer/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: SHOULD ONLY SEE ACCOUNTS THAT ARE A PART OF BANKS THAT YOU OWN!

type AccountController struct {
	service              services.AccountService
	moneyTransferService services.MoneyTransferService
}

func NewAccountController(
	account services.AccountService,
	moneyTransferService services.MoneyTransferService,
) *AccountController {
	return &AccountController{account, moneyTransferService}
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

func (controller AccountController) FindMoneyTransfers(c *gin.Context) {
	accountID := c.Param("id")
	var moneyTransfers []models.MoneyTransfer
	err := controller.moneyTransferService.FindByAccount(accountID, &moneyTransfers)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
		return
	}

	if moneyTransfers == nil {
		moneyTransfers = make([]models.MoneyTransfer, 0)
	}

	c.JSON(http.StatusOK, moneyTransfers)
}
