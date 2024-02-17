package auth

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"strconv"
	"strings"
)

type CustomerLoginRequest struct {
	BankID string `json:"bank_id"`
	PIN    string `json:"pin"`
}

type CustomerAuth interface {
	Login(request CustomerLoginRequest) (string, model.Customer, error)
}

type customerAuth struct {
	customerRepository repository.CustomerRepository
	jwt                JWTService
}

func NewCustomerAuth(customerRepository repository.CustomerRepository) CustomerAuth {
	return customerAuth{
		customerRepository: customerRepository,
		jwt:                NewJWTService(),
	}
}

func (auth customerAuth) Login(request CustomerLoginRequest) (string, model.Customer, error) {
	request.PIN = strings.TrimSpace(request.PIN)

	var customer model.Customer
	if err := auth.customerRepository.FindByBankAndPIN(request.BankID, request.PIN, &customer); err != nil {
		return "", customer, err
	}

	token, err := auth.jwt.GenerateCustomerToken(strconv.Itoa(int(customer.ID)))

	if err != nil {
		return "", customer, err
	}

	return token, customer, err
}
