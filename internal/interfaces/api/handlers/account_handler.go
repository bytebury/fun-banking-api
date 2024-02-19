package handlers

import (
	"funbanking/internal/domain/banking"
	"funbanking/internal/domain/users"
	"funbanking/internal/infrastructure/auth"
	"funbanking/package/constants"
	"funbanking/package/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountService  banking.AccountService
	bankService     banking.BankService
	employeeService banking.EmployeeService
	userService     users.UserService
}

func NewAccountHandler() AccountHandler {
	userRepository := users.NewUserRepository()
	return AccountHandler{
		accountService: banking.NewAccountService(
			banking.NewAccountRepository(),
		),
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

	transactions, err := h.accountService.FindTransactions(accountID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find that account"})
		return
	}

	c.JSON(http.StatusOK, utils.Listify(transactions))
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusAccepted, account)
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
