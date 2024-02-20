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
	FindAllByUserID(userID string) ([]Bank, error)
	Create(userID string, bank *Bank) error
	Update(id string, bank *Bank) error
	Delete(id string) error
	IsOwner(bankID, userID string) bool
	IsEmployee(bankID, userID string) bool
}

type bankService struct {
	bankRepository     BankRepository
	employeeRepository EmployeeRepository
}

func NewBankService(bankRepository BankRepository) BankService {
	return bankService{
		bankRepository:     bankRepository,
		employeeRepository: NewEmployeeRepository(),
	}
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

func (s bankService) FindAllByUserID(userID string) ([]Bank, error) {
	var banks []Bank
	var employeeOf []Employee
	var employeeAtBanks []Bank = make([]Bank, 0)

	if err := s.bankRepository.FindAllByUserID(userID, &banks); err != nil {
		return banks, err
	}

	if err := s.employeeRepository.FindAllByUserID(userID, &employeeOf); err != nil {
		return banks, err
	}

	for _, employee := range employeeOf {
		employeeAtBanks = append(employeeAtBanks, employee.Bank)
	}

	if banks == nil {
		banks = make([]Bank, 0)
	}

	return append(banks, employeeAtBanks...), nil
}

func (s bankService) Create(userID string, bank *Bank) error {
	userIDAsUInt, err := strconv.Atoi(userID)

	if err != nil {
		return err
	}

	bank.UserID = uint(userIDAsUInt)
	return s.bankRepository.Create(bank)
}

func (s bankService) Update(id string, bank *Bank) error {
	return s.bankRepository.Update(id, bank)
}

func (s bankService) Delete(id string) error {
	return s.bankRepository.Delete(id)
}

func (s bankService) IsOwner(bankID string, userID string) bool {
	var bank Bank

	if err := s.bankRepository.FindByID(bankID, &bank); err != nil {
		return false
	}

	return strconv.Itoa(int(bank.UserID)) == userID || utils.IsAdmin(userID)
}

func (s bankService) IsEmployee(bankID string, userID string) bool {
	var bank Bank
	var employees []Employee

	if err := s.bankRepository.FindByID(bankID, &bank); err != nil {
		return false
	}

	if s.IsOwner(bankID, userID) {
		return true
	}

	if err := s.employeeRepository.FindAllByBankID(bankID, &employees); err != nil {
		return false
	}

	for _, employee := range employees {
		if strconv.Itoa(int(employee.UserID)) == userID {
			return true
		}
	}

	return false
}
