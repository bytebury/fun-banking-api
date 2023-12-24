package controllers

import (
	"golfer/models"
	"golfer/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransferController struct {
	service        services.TransferService
	accountService services.AccountService
}

func NewTransferController(
	service services.TransferService,
	accountService services.AccountService,
) *TransferController {
	return &TransferController{service, accountService}
}

func (controller TransferController) Create(c *gin.Context) {
	var transfer models.Transfer
	userID := c.GetString("user_id")

	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	if transfer.Amount == 0 || transfer.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	var account models.Account
	if err := controller.accountService.FindByID(strconv.Itoa(int(transfer.AccountID)), &account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to find that account"})
		return
	}

	if transfer.Amount > 100_000_000 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "You can't transfer that much money at once"})
		return
	}

	transfer.CurrentBalance = account.Balance

	err := controller.service.Create(&transfer, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong creating your transfer"})
		return
	}

	c.JSON(http.StatusCreated, transfer)
}

func (controller TransferController) Approve(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	transferID := c.Param("id")

	if !controller.isBankStaff(transferID, c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that"})
		return
	}

	transfer, err := controller.service.Approve(transferID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to approve the transfer"})
		return
	}

	c.JSON(http.StatusOK, transfer)
}

func (controller TransferController) Decline(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	transferID := c.Param("id")

	if !controller.isBankStaff(transferID, c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that"})
		return
	}

	transfer, err := controller.service.Decline(transferID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to decline the transfer"})
		return
	}

	c.JSON(http.StatusOK, transfer)
}

func (controller TransferController) Notifications(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	var transfers []models.Transfer

	if err := controller.service.Notifications(userID, &transfers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, transfers)

}

func (controller TransferController) isBankStaff(transferID string, c *gin.Context) bool {
	userID := c.MustGet("user_id").(string)

	var transfer models.Transfer
	if err := controller.service.FindByID(transferID, &transfer); err != nil {
		return false
	}

	return strconv.Itoa(int(transfer.Account.Customer.Bank.UserID)) == userID
}