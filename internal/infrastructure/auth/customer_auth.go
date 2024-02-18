package auth

import (
	"funbanking/internal/domain/banking"
	"strconv"
	"strings"
)

type CustomerAuth interface {
	Login(request banking.CustomerLoginRequest) (string, banking.Customer, error)
}

type customerAuth struct {
	customerRepository banking.CustomerRepository
	jwt                JWTService
}

func NewCustomerAuth(customerRepository banking.CustomerRepository) CustomerAuth {
	return customerAuth{
		customerRepository: customerRepository,
		jwt:                NewJWTService(),
	}
}

func (auth customerAuth) Login(request banking.CustomerLoginRequest) (string, banking.Customer, error) {
	request.PIN = strings.TrimSpace(request.PIN)

	var customer banking.Customer
	if err := auth.customerRepository.FindByBankAndPIN(request.BankID, request.PIN, &customer); err != nil {
		return "", customer, err
	}

	token, err := auth.jwt.GenerateCustomerToken(strconv.Itoa(int(customer.ID)))

	if err != nil {
		return "", customer, err
	}

	return token, customer, err
}
