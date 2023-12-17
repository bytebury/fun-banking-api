package controllers

import (
	"golfer/models"
	"golfer/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MoneyTransferController struct {
	service        services.MoneyTransferService
	accountService services.AccountService
}

func NewMoneyTransferController(
	service services.MoneyTransferService,
	accountService services.AccountService,
) *MoneyTransferController {
	return &MoneyTransferController{service, accountService}
}

func (controller MoneyTransferController) Create(c *gin.Context) {
	var moneyTransfer models.MoneyTransfer

	if err := c.ShouldBindJSON(&moneyTransfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	if moneyTransfer.Amount == 0 || moneyTransfer.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	var account models.Account
	if err := controller.accountService.FindByID(strconv.Itoa(int(moneyTransfer.AccountID)), &account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to find that account"})
		return
	}

	moneyTransfer.CurrentBalance = account.Balance

	err := controller.service.Create(&moneyTransfer)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong creating your transfer"})
		return
	}

	c.JSON(http.StatusCreated, moneyTransfer)
}

// TODO: CANNOT APPROVE TRANSACTIONS UNLESS IT IS YOUR BANK
func (controller MoneyTransferController) Approve(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	moneyTransferID := c.Param("id")

	transfer, err := controller.service.Approve(moneyTransferID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to approve the transfer"})
		return
	}

	c.JSON(http.StatusOK, transfer)
}

// TODO: CANNOT DECLINE TRANSACTIONS UNLESS IT IS YOUR BANK
func (controller MoneyTransferController) Decline(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	moneyTransferID := c.Param("id")

	transfer, err := controller.service.Decline(moneyTransferID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to decline the transfer"})
		return
	}

	c.JSON(http.StatusOK, transfer)
}
