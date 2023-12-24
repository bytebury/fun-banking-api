package services

import (
	"golfer/models"
	"golfer/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MoneyTransferService struct {
	repository     repositories.MoneyTransferRepository
	accountService AccountService
	userService    UserService
}

func NewMoneyTransferService(
	repository repositories.MoneyTransferRepository,
	accountService AccountService,
	userService UserService,
) *MoneyTransferService {
	return &MoneyTransferService{
		repository,
		accountService,
		userService,
	}
}

func (service MoneyTransferService) Create(request *models.MoneyTransfer, userID string) error {
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

func (service MoneyTransferService) FindByID(moneyTransferID string, moneyTransfer *models.MoneyTransfer) error {
	return service.repository.FindByID(moneyTransferID, moneyTransfer)
}

func (service MoneyTransferService) FindByAccount(accountID string, moneyTransfers *[]models.MoneyTransfer, count *int64, c *gin.Context) error {
	return service.repository.FindByAccount(accountID, moneyTransfers, count, c)
}

func (service MoneyTransferService) Approve(moneyTransferID, userID string) (models.MoneyTransfer, error) {
	var moneyTransfer models.MoneyTransfer

	if err := service.repository.FindByID(moneyTransferID, &moneyTransfer); err != nil {
		return moneyTransfer, err
	}

	currentUserID, _ := stringToUintPtr(userID)

	moneyTransfer.Status = "approved"
	moneyTransfer.UserID = currentUserID
	moneyTransfer.CurrentBalance = moneyTransfer.Account.Balance + float64(moneyTransfer.Amount)

	if err := service.repository.Update(&moneyTransfer); err != nil {
		return moneyTransfer, err
	}

	service.accountService.UpdateBalance(
		strconv.Itoa((int(moneyTransfer.AccountID))),
		moneyTransfer.CurrentBalance)

	return moneyTransfer, nil
}

func (service MoneyTransferService) Decline(moneyTransferID, userID string) (models.MoneyTransfer, error) {
	var moneyTransfer models.MoneyTransfer

	if err := service.repository.FindByID(moneyTransferID, &moneyTransfer); err != nil {
		return moneyTransfer, err
	}

	currentUserID, _ := stringToUintPtr(userID)

	moneyTransfer.Status = "declined"
	moneyTransfer.UserID = currentUserID

	if err := service.repository.Update(&moneyTransfer); err != nil {
		return moneyTransfer, err
	}

	return moneyTransfer, nil
}

func (service MoneyTransferService) Notifications(userID string, transfers *[]models.MoneyTransfer) error {
	return service.repository.FindByUserID(userID, transfers)
}

func (service MoneyTransferService) isBankStaff(accountId, userId string) bool {
	var user models.User

	if err := service.userService.FindByID(userId, &user); err != nil {
		return false
	}

	var account models.Account
	if err := service.accountService.FindByID(accountId, &account); err != nil {
		return false
	}

	return account.Customer.Bank.UserID == user.ID
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
