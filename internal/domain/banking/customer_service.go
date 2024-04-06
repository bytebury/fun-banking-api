package banking

import (
	"errors"
	"funbanking/package/utils"
	"strconv"
)

type LoginRequest struct {
	BankID string `json:"bank_id"`
	PIN    string `json:"pin"`
}

type CustomerAuth interface {
	Login(request LoginRequest) (string, Customer, error)
}

type CustomerService interface {
	FindByID(id string) (Customer, error)
	FindAccounts(id string) ([]Account, error)
	Login(bankId string, pin string) (string, Customer, error)
	Create(customer *Customer) error
	Update(id string, customer *Customer) error
	Delete(id string) error
}

type customerService struct {
	authService        CustomerAuth
	customerRepository CustomerRepository
	bankService        BankService
}

func NewCustomerService(customerRepository CustomerRepository, authService CustomerAuth, bankService BankService) CustomerService {
	return customerService{
		customerRepository: customerRepository,
		authService:        authService,
		bankService:        bankService,
	}
}

func (s customerService) FindByID(id string) (Customer, error) {
	var customer Customer
	err := s.customerRepository.FindByID(id, &customer)
	return customer, err
}

func (s customerService) FindAccounts(id string) ([]Account, error) {
	var accounts []Account
	err := s.customerRepository.FindAccounts(id, &accounts)
	return utils.Listify(accounts), err
}

func (s customerService) Login(bankID string, pin string) (string, Customer, error) {
	request := LoginRequest{
		BankID: bankID,
		PIN:    pin,
	}
	return s.authService.Login(request)
}

func (s customerService) Create(customer *Customer) error {
	if customers, err := s.bankService.FindAllCustomers(strconv.Itoa(int(customer.BankID))); err != nil {
		return err
	} else if EnablePremium && len(customers) >= BankConfig.Limits.Free.Customers {
		return errors.New("limit reached")
	}

	return s.customerRepository.Create(customer)
}

func (s customerService) Update(customerID string, customer *Customer) error {
	return s.customerRepository.Update(customerID, customer)
}

func (s customerService) Delete(customerID string) error {
	return s.customerRepository.Delete(customerID)
}
