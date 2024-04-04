package banking

import (
	"errors"
	"funbanking/internal/infrastructure/pagination"
	"strconv"
)

type AccountService interface {
	FindByID(accountID string) (Account, error)
	FindTransactions(accountID string, statuses []string, itemsPerPage int, pageNumber int, params map[string]string) (pagination.PaginatedResponse[Transaction], error)
	MonthlyTransactionInsights(accountID string) ([]AccountMonthlySummary, error)
	Update(accountID string, account *Account) error
	AddToBalance(accountID string, amount float64) (Account, error)
	Transfer(customerID string, transferRequest TransferRequest) error
	Create(userID string, account *Account) error
}

type accountService struct {
	accountRepository  AccountRepository
	customerRepository CustomerRepository
}

func NewAccountService(accountRepository AccountRepository, customerRepository CustomerRepository) AccountService {
	return accountService{accountRepository, customerRepository}
}

func (s accountService) FindByID(accountID string) (Account, error) {
	var account Account
	err := s.accountRepository.FindByID(accountID, &account)
	return account, err
}

func (s accountService) FindTransactions(accountID string, statuses []string, itemsPerPage int, pageNumber int, params map[string]string) (pagination.PaginatedResponse[Transaction], error) {
	return s.accountRepository.FindTransactions(accountID, statuses, itemsPerPage, pageNumber, params)
}

func (s accountService) MonthlyTransactionInsights(accountID string) ([]AccountMonthlySummary, error) {
	return s.accountRepository.MonthlyTransactionInsights(accountID)
}

func (s accountService) Update(accountID string, account *Account) error {
	return s.accountRepository.Update(accountID, account)
}

func (s accountService) AddToBalance(accountID string, amount float64) (Account, error) {
	return s.accountRepository.AddToBalance(accountID, amount)
}

func (s accountService) Create(userID string, account *Account) error {
	var customer Customer
	customerID := strconv.Itoa(int(account.CustomerID))

	if err := s.customerRepository.FindByID(customerID, &customer); err != nil {
		return err
	}

	if !s.userHasAccessToCustomer(userID, customer) {
		return errors.New("not allowed")
	}

	if len(customer.Accounts) > 1 {
		return errors.New("maximum number of accounts reached")
	}

	return s.accountRepository.Create(account)
}

func (s accountService) Transfer(customerID string, transferRequest TransferRequest) error {
	// Check to make sure that the customer owns both accounts
	return s.accountRepository.Transfer(transferRequest)
}

func (s accountService) userHasAccessToCustomer(userID string, customer Customer) bool {
	// TODO Need to account for employees and admins, too
	return userID == strconv.Itoa(int(customer.Bank.UserID))
}
