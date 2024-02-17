package service

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"funbanking/internal/infrastructure/auth"
	"funbanking/package/utils"
)

type CustomerService interface {
	FindByID(id string) (model.Customer, error)
	FindAccounts(id string) ([]model.Account, error)
	Login(bankId string, pin string) (string, model.Customer, error)
	Create(customer *model.Customer) error
	Update(id string, customer *model.Customer) error
	Delete(id string) error
}

type customerService struct {
	authService        auth.CustomerAuth
	customerRepository repository.CustomerRepository
}

func NewCustomerService(customerRepository repository.CustomerRepository) CustomerService {
	return customerService{
		customerRepository: customerRepository,
		authService:        auth.NewCustomerAuth(customerRepository),
	}
}

func (s customerService) FindByID(id string) (model.Customer, error) {
	var customer model.Customer
	err := s.customerRepository.FindByID(id, &customer)
	return customer, err
}

func (s customerService) FindAccounts(id string) ([]model.Account, error) {
	var accounts []model.Account
	err := s.customerRepository.FindAccounts(id, &accounts)
	return utils.Listify(accounts), err
}

func (s customerService) Login(bankID string, pin string) (string, model.Customer, error) {
	request := auth.CustomerLoginRequest{
		BankID: bankID,
		PIN:    pin,
	}
	return s.authService.Login(request)
}

func (s customerService) Create(customer *model.Customer) error {
	return s.customerRepository.Create(customer)
}

func (s customerService) Update(id string, customer *model.Customer) error {
	return s.customerRepository.Update(id, customer)
}

func (s customerService) Delete(id string) error {
	return s.customerRepository.Delete(id)
}
