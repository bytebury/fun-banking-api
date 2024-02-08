package service

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"funbanking/package/utils"
)

type BankService interface {
	FindByID(id string) (model.Bank, error)
	FindByUsernameAndSlug(username, slug string) (model.Bank, error)
	FindAllCustomers(id string) ([]model.Customer, error)
	Create(bank *model.Bank) error
	Update(id string, bank *model.Bank) error
	Delete(id string) error
}

type bankService struct {
	bankRepository repository.BankRepository
}

func NewBankService(bankRepository repository.BankRepository) BankService {
	return bankService{bankRepository}
}

func (s bankService) FindByID(id string) (model.Bank, error) {
	var bank model.Bank
	err := s.bankRepository.FindByID(id, &bank)
	return bank, err
}

func (s bankService) FindByUsernameAndSlug(username, slug string) (model.Bank, error) {
	var bank model.Bank
	err := s.bankRepository.FindByUsernameAndSlug(username, slug, &bank)
	return bank, err
}

func (s bankService) FindAllCustomers(id string) ([]model.Customer, error) {
	var customers []model.Customer
	err := s.bankRepository.FindAllCustomers(id, &customers)
	return utils.Listify(customers), err
}

func (s bankService) Create(bank *model.Bank) error {
	return s.bankRepository.Create(bank)
}

func (s bankService) Update(id string, bank *model.Bank) error {
	return s.bankRepository.Update(id, bank)
}

func (s bankService) Delete(id string) error {
	return s.bankRepository.Delete(id)
}
