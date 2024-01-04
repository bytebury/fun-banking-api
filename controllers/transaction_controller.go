package controllers

import (
	"golfer/config"
	"golfer/models"
	"golfer/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	service         services.TransactionService
	accountService  services.AccountService
	employeeService services.EmployeeService
}

func NewTransactionController(
	service services.TransactionService,
	accountService services.AccountService,
	employeeService services.EmployeeService,
) *TransactionController {
	return &TransactionController{service, accountService, employeeService}
}

func (controller TransactionController) Create(c *gin.Context) {
	var transfer models.Transaction
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

	// TODO MOVE THIS TO THE SERVICE!!
	if transfer.Amount > config.MAX_BANKING_TRANSFER_AMOUNT {
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

func (controller TransactionController) Approve(c *gin.Context) {
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

func (controller TransactionController) Decline(c *gin.Context) {
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

func (controller TransactionController) Notifications(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	var transfers []models.Transaction

	if err := controller.service.Notifications(userID, &transfers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, transfers)

}

func (controller TransactionController) isBankStaff(transferID string, c *gin.Context) bool {
	userID := c.MustGet("user_id").(string)

	var transfer models.Transaction
	if err := controller.service.FindByID(transferID, &transfer); err != nil {
		return false
	}

	if strconv.Itoa(int(transfer.Account.Customer.Bank.UserID)) == userID {
		return true
	}

	var employees []models.Employee
	if err := controller.employeeService.FindByBank(strconv.Itoa(int(transfer.Account.Customer.Bank.ID)), &employees); err != nil {
		return false
	}

	for _, employee := range employees {
		if strconv.Itoa(int(employee.UserID)) == userID {
			return true
		}
	}

	return false
}
