package controllers

import (
	"golfer/models"
	"golfer/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	service        services.AccountService
	transferervice services.TransferService
}

func NewAccountController(
	account services.AccountService,
	transferervice services.TransferService,
) *AccountController {
	return &AccountController{account, transferervice}
}

func (controller AccountController) FindByID(c *gin.Context) {
	accountID := c.Param("id")

	var account models.Account
	if err := controller.service.FindByID(accountID, &account); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
		return
	}

	// TODO: Only that specific customer or bank staff can look accounts up.
	// if !controller.isBankStaff(account, c) {
	// 	c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that resource"})
	// 	return
	// }

	c.JSON(http.StatusOK, account)
}

func (controller AccountController) FindTransfers(c *gin.Context) {
	accountID := c.Param("id")

	var account models.Account
	if err := controller.service.FindByID(accountID, &account); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
		return
	}

	// TODO: Only that specific customer or bank staff can look accounts up.
	// if !controller.isBankStaff(account, c) {
	// 	c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that resource"})
	// 	return
	// }

	var transfers []models.Transfer
	var count int64
	if err := controller.transferervice.FindByAccount(accountID, &transfers, &count, c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
		return
	}

	if transfers == nil {
		transfers = make([]models.Transfer, 0)
	}

	pageNumber, err := strconv.Atoi(c.Query("page"))

	if err != nil {
		pageNumber = 1
	}

	paginatedResponse := models.PaginatedResponse[models.Transfer]{
		Items: transfers,
		PagingInfo: models.PagingInfo{
			TotalItems: uint(count),
			PageNumber: uint(pageNumber),
		},
	}

	c.JSON(http.StatusOK, paginatedResponse)
}

func (controller AccountController) GetTransferHistoricalData(c *gin.Context) {
	accountID := c.Param("id")
	daysAgo, err := strconv.Atoi(c.Query("days-ago"))

	if err != nil || daysAgo <= 0 {
		daysAgo = 30
	}

	summary, err := controller.service.GetTransferHistoricalData(accountID, daysAgo)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something happened while retrieving data"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// func (controller AccountController) isBankStaff(account models.Account, c *gin.Context) bool {
// 	userID := c.GetString("user_id")
// 	return strconv.Itoa(int(account.Customer.Bank.UserID)) == userID
// }
