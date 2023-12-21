package services

import (
	"golfer/models"
	"golfer/repositories"
	"strings"
)

type BankService struct {
	bankRepository repositories.BankRepository
}

func NewBankService(bankRepository repositories.BankRepository) *BankService {
	return &BankService{
		bankRepository,
	}
}

func (service BankService) Create(request *models.Bank) error {
	request.Slug = strings.ToLower(request.Slug)
	return service.bankRepository.Create(request)
}

func (service BankService) FindByID(bankID string, bank *models.Bank) error {
	return service.bankRepository.FindByID(bankID, bank)
}

func (service BankService) FindBanksByUserID(userID string, banks *[]models.Bank) error {
	return service.bankRepository.FindBanksByUserID(userID, banks)
}

func (service BankService) FindCustomers(bankID string, customers *[]models.Customer) error {
	return service.bankRepository.FindCustomers(bankID, customers)
}

func (service BankService) Update(bankID string, request *models.Bank) (models.Bank, error) {
	var bank models.Bank

	if err := service.bankRepository.FindByID(bankID, &bank); err != nil {
		return bank, err
	}

	if request.Slug != "" {
		bank.Slug = strings.ToLower(request.Slug)
	}

	if request.Name != "" {
		bank.Name = request.Name
	}

	if request.Description != "" {
		bank.Description = request.Description
	}

	if err := service.bankRepository.Update(&bank); err != nil {
		return bank, err
	}

	return bank, nil
}

func (service BankService) FindByUsernameAndSlug(username, slug string, bank *models.Bank) error {
	return service.bankRepository.FindByUsernameAndSlug(username, slug, bank)
}

func (service BankService) Delete(bankID string) error {
	return service.bankRepository.Delete(bankID)
}
