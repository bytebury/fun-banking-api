package services

import (
	"golfer/models"
	"golfer/repositories"
	"strconv"
	"strings"
)

type CustomerService struct {
	repository repositories.CustomerRepository
	jwtService JwtService
}

func NewCustomerService(repository repositories.CustomerRepository, jwtService JwtService) *CustomerService {
	return &CustomerService{
		repository,
		jwtService,
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

func (service CustomerService) Login(request models.CustomerSignInRequest) (string, models.Customer, error) {
	var customer models.Customer

	if err := service.repository.FindByBankAndPIN(request.BankID, request.PIN, &customer); err != nil {
		return "", customer, err
	}

	token, err := service.jwtService.GenerateCustomerToken(strconv.Itoa(int(customer.ID)))

	if err != nil {
		return "", customer, err
	}

	return token, customer, err
}
