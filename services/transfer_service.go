package services

import (
	"errors"
	"golfer/models"
	"golfer/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransferService struct {
	repository      repositories.TransferRepository
	accountService  AccountService
	userService     UserService
	employeeService EmployeeService
}

func NewTransferService(
	repository repositories.TransferRepository,
	accountService AccountService,
	userService UserService,
	employeeService EmployeeService,
) *TransferService {
	return &TransferService{
		repository,
		accountService,
		userService,
		employeeService,
	}
}

func (service TransferService) Create(request *models.Transfer, userID string) error {
	if err := service.repository.Create(request); err != nil {
		return err
	}

	accountID := strconv.Itoa(int(request.AccountID))
	requestID := strconv.Itoa(int(request.ID))

	if !service.isBankStaff(accountID, userID) {
		return nil
	}

	_, err := service.Approve(requestID, userID)
	return err
}

func (service TransferService) FindByID(transferID string, transfer *models.Transfer) error {
	return service.repository.FindByID(transferID, transfer)
}

func (service TransferService) FindByAccount(accountID string, transfers *[]models.Transfer, count *int64, c *gin.Context) error {
	return service.repository.FindByAccount(accountID, transfers, count, c)
}

func (service TransferService) Approve(transferID, userID string) (models.Transfer, error) {
	var transfer models.Transfer

	if err := service.repository.FindByID(transferID, &transfer); err != nil {
		return transfer, err
	}

	if transfer.Status != "pending" {
		return transfer, errors.New("this transfer was already processed")
	}

	currentUserID, _ := stringToUintPtr(userID)

	transfer.Status = "approved"
	transfer.UserID = currentUserID
	transfer.CurrentBalance = transfer.Account.Balance + float64(transfer.Amount)

	if err := service.repository.Update(&transfer); err != nil {
		return transfer, err
	}

	service.accountService.UpdateBalance(
		strconv.Itoa((int(transfer.AccountID))),
		transfer.CurrentBalance)

	return transfer, nil
}

func (service TransferService) Decline(transferID, userID string) (models.Transfer, error) {
	var transfer models.Transfer

	if err := service.repository.FindByID(transferID, &transfer); err != nil {
		return transfer, err
	}

	if transfer.Status != "pending" {
		return transfer, errors.New("this transfer was already processed")
	}

	currentUserID, _ := stringToUintPtr(userID)

	transfer.Status = "declined"
	transfer.UserID = currentUserID

	if err := service.repository.Update(&transfer); err != nil {
		return transfer, err
	}

	return transfer, nil
}

func (service TransferService) Notifications(userID string, transfers *[]models.Transfer) error {
	return service.repository.FindByUserID(userID, transfers)
}

func (service TransferService) isBankStaff(accountId, userId string) bool {
	var user models.User

	if err := service.userService.FindByID(userId, &user); err != nil {
		return false
	}

	var account models.Account
	if err := service.accountService.FindByID(accountId, &account); err != nil {
		return false
	}

	if account.Customer.Bank.UserID == user.ID {
		return true
	}

	var employees []models.Employee
	if err := service.employeeService.FindByBank(strconv.Itoa(int(account.Customer.BankID)), &employees); err != nil {
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
