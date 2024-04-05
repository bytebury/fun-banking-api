package handlers

import (
	"funbanking/internal/domain/banking"
	"funbanking/internal/domain/users"
	"funbanking/internal/infrastructure/auth"
	"funbanking/internal/infrastructure/mailing"
	"funbanking/package/constants"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountService  banking.AccountService
	bankService     banking.BankService
	employeeService banking.EmployeeService
	userService     users.UserService
	transferService banking.TransferService
}

func NewAccountHandler() AccountHandler {
	userRepository := users.NewUserRepository()

	accountService := banking.NewAccountService(
		banking.NewAccountRepository(),
		banking.NewCustomerRepository(),
	)

	transactionService := banking.NewTransactionService(
		banking.NewTransactionRepository(),
	)

	return AccountHandler{
		accountService: accountService,
		bankService: banking.NewBankService(
			banking.NewBankRepository(),
		),
		employeeService: banking.NewEmployeeService(
			banking.NewEmployeeRepository(),
		),
		userService: users.NewUserService(
			userRepository,
			auth.NewUserAuth(
				userRepository,
			),
			mailing.NewWelcomeMailer(),
		),
		transferService: banking.NewTransferService(
			accountService,
			transactionService,
		),
	}
}

func (h AccountHandler) FindByID(c *gin.Context) {
	accountID := c.Param("id")
	userID := c.GetString("user_id")
	customerID := c.GetString("customer_id")

	if !h.isEmployee(accountID, userID) && !h.isOwner(accountID, customerID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have access to view that account"})
		return
	}

	account, err := h.accountService.FindByID(accountID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find that account"})
		return
	}

	if customerID != "" {
		h.userService.AddVisitor(&users.Visitor{
			UserID:     nil,
			CustomerID: &account.CustomerID,
			IPAddress:  c.ClientIP(),
		})
	}

	c.JSON(http.StatusOK, account)
}

func (h AccountHandler) FindTransactions(c *gin.Context) {
	accountID := c.Param("id")
	userID := c.GetString("user_id")
	customerID := c.GetString("customer_id")

	if !h.isEmployee(accountID, userID) && !h.isOwner(accountID, customerID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have access to view these transactions"})
		return
	}

	statuses := c.QueryArray("status")
	itemsPerPage, _ := strconv.Atoi(c.Query("itemsPerPage"))
	pageNumber, _ := strconv.Atoi(c.Query("pageNumber"))

	params := map[string]string{
		"StartDate": c.Query("startDate"),
		"EndDate":   c.Query("endDate"),
		"Direction": c.Query("direction"),
	}

	transactions, err := h.accountService.FindTransactions(accountID, statuses, itemsPerPage, pageNumber, params)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find that account"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h AccountHandler) MonthlyTransactionInsights(c *gin.Context) {
	accountID := c.Param("id")

	summary, err := h.accountService.MonthlyTransactionInsights(accountID)

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

func (h AccountHandler) Update(c *gin.Context) {
	var account banking.Account

	accountID := c.Param("id")
	userID := c.GetString("user_id")

	if !h.isEmployee(accountID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have access to update this account"})
		return
	}

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := h.accountService.Update(accountID, &account); err != nil {
		if strings.Contains(err.Error(), "name is too long") {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusAccepted, account)
}

func (h AccountHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	var account banking.Account

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := h.accountService.Create(userID, &account); err != nil {
		if strings.Contains(err.Error(), "required") {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if strings.Contains(err.Error(), "balances") {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if strings.Contains(err.Error(), "not allowed") {
			c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to do that"})
			return
		}

		if strings.Contains(err.Error(), "maximum") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Maximum number of accounts already reached"})
			return
		}
	}

	c.JSON(http.StatusOK, account)
}

func (h AccountHandler) TransferBetweenAccounts(c *gin.Context) {
	userID := c.GetString("user_id")
	customerID := c.GetString("customer_id")

	var transferRequest banking.TransferRequest

	if err := c.ShouldBindJSON(&transferRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	fromAccountId := strconv.Itoa(int(transferRequest.FromAccountID))

	if h.isEmployee(fromAccountId, userID) {
		if account, err := h.accountService.FindByID(fromAccountId); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "One or both accounts do not exist"})
			return
		} else {
			customerID = strconv.Itoa(int(account.CustomerID))
		}
	}

	if customerID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not find you or your customer"})
		return
	}

	if err := h.transferService.Transfer(customerID, transferRequest); err != nil {
		if strings.Contains(err.Error(), "account does not exist") {
			c.JSON(http.StatusNotFound, gin.H{"message": "One or both accounts do not exist"})
			return
		}
		if strings.Contains(err.Error(), "between accounts you own") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "You can only transfer between accounts you own"})
			return
		}
		if strings.Contains(err.Error(), "insufficient funds") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "You do not have enough funds to transfer that much"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed transfering funds between accounts"})
		return
	}

	c.JSON(http.StatusAccepted, nil)
}

func (h AccountHandler) isOwner(accountID string, customerID string) bool {
	if customerID == "" {
		return false
	}

	account, err := h.accountService.FindByID(accountID)

	if err != nil {
		return false
	}

	return strconv.Itoa(int(account.CustomerID)) == customerID
}

func (h AccountHandler) isEmployee(accountID string, userID string) bool {
	if userID == "" {
		return false
	}

	account, err := h.accountService.FindByID(accountID)

	if err != nil {
		return false
	}

	bank, err := h.bankService.FindByID(strconv.Itoa(int(account.Customer.BankID)))

	if err != nil {
		return false
	}

	if strconv.Itoa(int(bank.UserID)) == userID {
		return true
	}

	user, err := h.userService.FindByID(userID)

	if err != nil {
		return false
	}

	if user.Role >= constants.AdminRole {
		return true
	}

	employees, err := h.employeeService.FindAllByBankID(strconv.Itoa(int(bank.ID)))

	if err != nil {
		return false
	}

	for _, employee := range employees {
		if strconv.Itoa(int(employee.UserID)) == userID {
			return true
		}
	}

	return false
}
