package controllers

import (
	"golfer/config"
	"golfer/models"
	"golfer/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	accountService     services.AccountService
	transactionService services.TransactionService
	bankService        services.BankService
	userService        services.UserService
	employeeService    services.EmployeeService
}

func NewAccountController(
	accountService services.AccountService,
	transactionService services.TransactionService,
	bankService services.BankService,
	userService services.UserService,
	employeeService services.EmployeeService,
) *AccountController {
	return &AccountController{
		accountService,
		transactionService,
		bankService,
		userService,
		employeeService,
	}
}

func (ac AccountController) FindByID(c *gin.Context) {
	accountID := c.Param("id")

	var account models.Account
	if err := ac.accountService.FindByID(accountID, &account); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
		return
	}

	if !ac.hasAccess(account, c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that resource"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (ac AccountController) FindTransactions(c *gin.Context) {
	accountID := c.Param("id")

	var account models.Account
	if err := ac.accountService.FindByID(accountID, &account); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
		return
	}

	if !ac.hasAccess(account, c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that resource"})
		return
	}

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
		"labels":      labels,
		"deposits":    deposits,
		"withdrawals": withdrawals,
	})
}

func (ac AccountController) hasAccess(account models.Account, c *gin.Context) bool {
	return ac.isAccountOwner(account, c) || ac.isBankStaff(account, c)
}

func (ac AccountController) isAccountOwner(account models.Account, c *gin.Context) bool {
	customerID, exists := c.Get("customer_id")

	if !exists {
		return false
	}

	return strconv.Itoa(int(account.CustomerID)) == customerID
}

func (controller AccountController) isBankStaff(account models.Account, c *gin.Context) bool {
	userIDString, exists := c.Get("user_id")

	if !exists {
		return false
	}

	userID, _ := strconv.Atoi(userIDString.(string))

	var bank models.Bank
	if err := controller.bankService.FindByID(strconv.Itoa(int(account.Customer.BankID)), &bank); err != nil {
		return false
	}

	// You own the bank, so you're good to go
	if bank.UserID == uint(userID) {
		return true
	}

	var user models.User
	if err := controller.userService.FindByID(strconv.Itoa(userID), &user); err != nil {
		return false
	}

	// Admins have access to everything
	if user.Role >= config.AdminRole {
		return true
	}

	var employees []models.Employee
	if err := controller.employeeService.FindByBank(strconv.Itoa(int(bank.ID)), &employees); err != nil {
		return false
	}

	for _, employee := range employees {
		if int(employee.UserID) == userID {
			return true
		}
	}

	return false
}
