package controllers

import (
	"golfer/models"
	"golfer/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	accountService     services.AccountService
	transactionService services.TransactionService
}

func NewAccountController(
	accountService services.AccountService,
	transactionService services.TransactionService,
) *AccountController {
	return &AccountController{
		accountService,
		transactionService,
	}
}

func (ac AccountController) FindByID(c *gin.Context) {
	accountID := c.Param("id")

	var account models.Account
	if err := ac.accountService.FindByID(accountID, &account); err != nil {
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

func (ac AccountController) FindTransactions(c *gin.Context) {
	accountID := c.Param("id")

	var account models.Account
	if err := ac.accountService.FindByID(accountID, &account); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
		return
	}

	// TODO: Only that specific customer or bank staff can look accounts up.
	// if !controller.isBankStaff(account, c) {
	// 	c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that resource"})
	// 	return
	// }

	var transfers []models.Transaction
	var count int64
	if err := ac.transactionService.FindByAccount(accountID, &transfers, &count, c); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
		return
	}

	if transfers == nil {
		transfers = make([]models.Transaction, 0)
	}

	pageNumber, err := strconv.Atoi(c.Query("page"))

	if err != nil {
		pageNumber = 1
	}

	paginatedResponse := models.PaginatedResponse[models.Transaction]{
		Items: transfers,
		PagingInfo: models.PagingInfo{
			TotalItems: uint(count),
			PageNumber: uint(pageNumber),
		},
	}

	c.JSON(http.StatusOK, paginatedResponse)
}

func (ac AccountController) GetTransactionHistoricalData(c *gin.Context) {
	accountID := c.Param("id")

	summary, err := ac.accountService.GetMonthlyData(accountID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something happened while retrieving data"})
		return
	}

	var labels []string
	var deposits []float64
	var withdrawals []float64

	for _, agg := range summary {
		labels = append(labels, agg.Month)
		deposits = append(deposits, agg.Deposits)
		withdrawals = append(withdrawals, agg.Withdrawals)
	}

	c.JSON(http.StatusOK, gin.H{
		"label":       labels,
		"deposits":    deposits,
		"withdrawals": withdrawals,
	})
}

// func (controller AccountController) isBankStaff(account models.Account, c *gin.Context) bool {
// 	userID := c.GetString("user_id")
// 	return strconv.Itoa(int(account.Customer.Bank.UserID)) == userID
// }
