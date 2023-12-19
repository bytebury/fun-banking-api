package controllers

import (
	"golfer/models"
	"golfer/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
	if err := controller.service.FindByID(accountID, &account); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
		return
	}

	if !controller.isBankStaff(account, c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that resource"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (controller AccountController) FindMoneyTransfers(c *gin.Context) {
	accountID := c.Param("id")

	var account models.Account
	if err := controller.service.FindByID(accountID, &account); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
		return
	}

	if !controller.isBankStaff(account, c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that resource"})
		return
	}

	var moneyTransfers []models.MoneyTransfer
	if err := controller.moneyTransferService.FindByAccount(accountID, &moneyTransfers, c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
		return
	}

	if moneyTransfers == nil {
		moneyTransfers = make([]models.MoneyTransfer, 0)
	}

	c.JSON(http.StatusOK, moneyTransfers)
}

func (controller AccountController) isBankStaff(account models.Account, c *gin.Context) bool {
	userID := c.MustGet("user_id").(string)
	return strconv.Itoa(int(account.Customer.Bank.UserID)) == userID
}
