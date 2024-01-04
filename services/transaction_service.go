package services

import (
	"errors"
	"golfer/models"
	"golfer/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionService struct {
	repository      repositories.TransactionRepository
	accountService  AccountService
	userService     UserService
	employeeService EmployeeService
}

func NewTransactionService(
	repository repositories.TransactionRepository,
	accountService AccountService,
	userService UserService,
	employeeService EmployeeService,
) *TransactionService {
	return &TransactionService{
		repository,
		accountService,
		userService,
		employeeService,
	}
}

func (ts TransactionService) Create(request *models.Transaction, userID string) error {
	if err := ts.repository.Create(request); err != nil {
		return err
	}

	accountID := strconv.Itoa(int(request.AccountID))
	requestID := strconv.Itoa(int(request.ID))

	if !ts.isBankStaff(accountID, userID) {
		return nil
	}

	_, err := ts.Approve(requestID, userID)
	return err
}

func (ts TransactionService) FindByID(transferID string, transfer *models.Transaction) error {
	return ts.repository.FindByID(transferID, transfer)
}

func (ts TransactionService) FindByAccount(accountID string, transfers *[]models.Transaction, count *int64, c *gin.Context) error {
	return ts.repository.FindByAccount(accountID, transfers, count, c)
}

func (ts TransactionService) Approve(transferID, userID string) (models.Transaction, error) {
	var transfer models.Transaction

	if err := ts.repository.FindByID(transferID, &transfer); err != nil {
		return transfer, err
	}

	if transfer.Status != "pending" {
		return transfer, errors.New("this transfer was already processed")
	}

	currentUserID, _ := stringToUintPtr(userID)

	transfer.Status = "approved"
	transfer.UserID = currentUserID
	transfer.CurrentBalance = transfer.Account.Balance + float64(transfer.Amount)

	if err := ts.repository.Update(&transfer); err != nil {
		return transfer, err
	}

	ts.accountService.UpdateBalance(
		strconv.Itoa((int(transfer.AccountID))),
		transfer.CurrentBalance)

	return transfer, nil
}

func (ts TransactionService) Decline(transferID, userID string) (models.Transaction, error) {
	var transfer models.Transaction

	if err := ts.repository.FindByID(transferID, &transfer); err != nil {
		return transfer, err
	}

	if transfer.Status != "pending" {
		return transfer, errors.New("this transfer was already processed")
	}

	currentUserID, _ := stringToUintPtr(userID)

	transfer.Status = "declined"
	transfer.UserID = currentUserID

	if err := ts.repository.Update(&transfer); err != nil {
		return transfer, err
	}

	return transfer, nil
}

func (ts TransactionService) Notifications(userID string, transfers *[]models.Transaction) error {
	return ts.repository.FindByUserID(userID, transfers)
}

func (ts TransactionService) isBankStaff(accountId, userId string) bool {
	var user models.User

	if err := ts.userService.FindByID(userId, &user); err != nil {
		return false
	}

	var account models.Account
	if err := ts.accountService.FindByID(accountId, &account); err != nil {
		return false
	}

	if account.Customer.Bank.UserID == user.ID {
		return true
	}

	var employees []models.Employee
	if err := ts.employeeService.FindByBank(strconv.Itoa(int(account.Customer.BankID)), &employees); err != nil {
		return false
	}

	for _, employee := range employees {
		if employee.UserID == user.ID {
			return true
		}
	}

	return false
}

func stringToUintPtr(s string) (*uint, error) {
	// Convert the string to a uint
	num, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return nil, err
	}

	// Convert the uint to a *uint
	uintPtr := uint(num)
	return &uintPtr, nil
}
