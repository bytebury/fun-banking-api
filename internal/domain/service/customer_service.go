package service

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"funbanking/package/utils"
)

type CustomerService interface {
	FindByID(id string) (model.Customer, error)
	FindAccounts(id string) ([]model.Account, error)
	Create(customer *model.Customer) error
	Update(id string, customer *model.Customer) error
	Delete(id string) error
}

type customerService struct {
	customerRepository repository.CustomerRepository
}

func NewCustomerService(customerRepository repository.CustomerRepository) CustomerService {
	return customerService{customerRepository}
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

func (s customerService) Create(customer *model.Customer) error {
	return s.customerRepository.Create(customer)
}

func (s customerService) Update(id string, customer *model.Customer) error {
	return s.customerRepository.Update(id, customer)
}

func (s customerService) Delete(id string) error {
	return s.customerRepository.Delete(id)
}
