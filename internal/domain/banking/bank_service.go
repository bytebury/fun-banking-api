package banking

import (
	"funbanking/package/utils"
	"strconv"
	"strings"
)

type BankService interface {
	FindByID(id string) (Bank, error)
	FindByUsernameAndSlug(username, slug string) (Bank, error)
	FindAllCustomers(id string) ([]Customer, error)
	Create(bank *Bank) error
	Update(id string, bank *Bank) error
	Delete(id string) error
	IsBankOwner(bankID, userID string) bool
}

type bankService struct {
	bankRepository BankRepository
}

func NewBankService(bankRepository BankRepository) BankService {
	return bankService{bankRepository}
}

func (s bankService) FindByID(id string) (Bank, error) {
	var bank Bank
	err := s.bankRepository.FindByID(id, &bank)
	return bank, err
}

func (s bankService) FindByUsernameAndSlug(username, slug string) (Bank, error) {
	var bank Bank

	// Normalize inputs
	username = strings.TrimSpace(strings.ToLower(username))
	slug = strings.TrimSpace(strings.ToLower(slug))

	err := s.bankRepository.FindByUsernameAndSlug(username, slug, &bank)
	return bank, err
}

func (s bankService) FindAllCustomers(id string) ([]Customer, error) {
	var customers []Customer
	err := s.bankRepository.FindAllCustomers(id, &customers)
	return utils.Listify(customers), err
}

func (s bankService) Create(bank *Bank) error {
	return s.bankRepository.Create(bank)
}

func (s bankService) Update(id string, bank *Bank) error {
	return s.bankRepository.Update(id, bank)
}

func (s bankService) Delete(id string) error {
	return s.bankRepository.Delete(id)
}

func (s bankService) IsBankOwner(bankID string, userID string) bool {
	var bank Bank

	if err := s.bankRepository.FindByID(bankID, &bank); err != nil {
		return false
	}

	return strconv.Itoa(int(bank.UserID)) == userID || utils.IsAdmin(userID)
}
