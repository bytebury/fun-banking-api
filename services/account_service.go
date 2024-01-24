package services

import (
	"golfer/models"
	"golfer/repositories"
)

type AccountService struct {
	repository repositories.AccountRepository
}

func NewAccountService(repository repositories.AccountRepository) *AccountService {
	return &AccountService{
		repository,
	}
}

func (service AccountService) Create(request *models.Account) error {
	return service.repository.Create(request)
}

func (service AccountService) FindByID(accountID string, account *models.Account) error {
	return service.repository.FindByID(accountID, account)
}

func (service AccountService) FindByCustomer(customerID string, accounts *[]models.Account) error {
	return service.repository.FindByCustomer(customerID, accounts)
}

func (service AccountService) Update(accountID string, request *models.Account) (models.Account, error) {
	var account models.Account

	if err := service.repository.FindByID(accountID, &account); err != nil {
		return account, err
	}

	if request.Name != "" {
		account.Name = request.Name
	}

	if err := service.repository.Update(&account); err != nil {
		return account, err
	}

	return account, nil
}

func (service AccountService) UpdateBalance(accountID string, balance float64) (models.Account, error) {
	var account models.Account

	if err := service.repository.FindByID(accountID, &account); err != nil {
		return account, err
	}

	account.Balance = balance

	if err := service.repository.Update(&account); err != nil {
		return account, err
	}

	return account, nil
}

func (service AccountService) Delete(accountID string) error {
	return service.repository.Delete(accountID)
}

func (service AccountService) GetMonthlyData(accountID string) ([]models.AccountMonthlySummary, error) {
	return service.repository.GetMonthlyData(accountID)
}
