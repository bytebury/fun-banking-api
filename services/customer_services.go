package services

import (
	"golfer/models"
	"golfer/repositories"
	"strings"
)

type CustomerService struct {
	repository repositories.CustomerRepository
}

func NewCustomerService(repository repositories.CustomerRepository) *CustomerService {
	return &CustomerService{
		repository,
	}
}

func (service CustomerService) Create(request *models.Customer) error {
	request.FirstName = strings.ToLower(request.FirstName)
	request.LastName = strings.ToLower(request.LastName)
	return service.repository.Create(request)
}

func (service CustomerService) FindByID(customerID string, customer *models.Customer) error {
	return service.repository.FindByID(customerID, customer)
}

func (service CustomerService) Update(customerID string, request *models.Customer) (models.Customer, error) {
	var customer models.Customer

	if err := service.repository.FindByID(customerID, &customer); err != nil {
		return customer, err
	}

	if request.FirstName != "" {
		customer.FirstName = strings.ToLower(request.FirstName)
	}

	if request.LastName != "" {
		customer.LastName = strings.ToLower(request.LastName)
	}

	if request.PIN != "" {
		customer.PIN = request.PIN
	}

	if err := service.repository.Update(&customer); err != nil {
		return customer, err
	}

	return customer, nil
}

func (service CustomerService) Delete(customerID string) error {
	return service.repository.Delete(customerID)
}

func (service CustomerService) Login(bankID string, pin string, customer *models.Customer) error {
	return service.repository.FindByBankAndPIN(bankID, pin, customer)
}
