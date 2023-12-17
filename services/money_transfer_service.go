package services

import (
	"golfer/models"
	"golfer/repositories"
	"strconv"
)

type MoneyTransferService struct {
	repository     repositories.MoneyTransferRepository
	accountService AccountService
}

func NewMoneyTransferService(
	repository repositories.MoneyTransferRepository,
	accountService AccountService,
) *MoneyTransferService {
	return &MoneyTransferService{
		repository,
		accountService,
	}
}

func (service MoneyTransferService) Create(request *models.MoneyTransfer) error {
	return service.repository.Create(request)
}

func (service MoneyTransferService) FindByID(moneyTransferID string, moneyTransfer *models.MoneyTransfer) error {
	return service.repository.FindByID(moneyTransferID, moneyTransfer)
}

func (service MoneyTransferService) FindByAccount(accountID string, moneyTransfers *[]models.MoneyTransfer) error {
	return service.repository.FindByAccount(accountID, moneyTransfers)
}

func (service MoneyTransferService) Approve(moneyTransferID, userID string) (models.MoneyTransfer, error) {
	var moneyTransfer models.MoneyTransfer

	if err := service.repository.FindByID(moneyTransferID, &moneyTransfer); err != nil {
		return moneyTransfer, err
	}

	currentUserID, _ := stringToUintPtr(userID)

	moneyTransfer.Status = "approved"
	moneyTransfer.UserID = currentUserID

	if err := service.repository.Update(&moneyTransfer); err != nil {
		return moneyTransfer, err
	}

	service.accountService.UpdateBalance(
		strconv.Itoa((int(moneyTransfer.AccountID))),
		moneyTransfer.Account.Balance+float64(moneyTransfer.Amount))

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
