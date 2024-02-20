package banking

import (
	"funbanking/internal/infrastructure/pagination"
)

type AccountService interface {
	FindByID(accountID string) (Account, error)
	FindTransactions(accountID string, statuses []string, itemsPerPage int, pageNumber int) (pagination.PaginatedResponse[Transaction], error)
	Update(accountID string, account *Account) error
}

type accountService struct {
	accountRepository AccountRepository
}

func NewAccountService(accountRepository AccountRepository) AccountService {
	return accountService{accountRepository}
}

func (s accountService) FindByID(accountID string) (Account, error) {
	var account Account
	err := s.accountRepository.FindByID(accountID, &account)
	return account, err
}

// TODO: THIS IS GOING TO BE PAGINATED
func (s accountService) FindTransactions(accountID string, statuses []string, itemsPerPage int, pageNumber int) (pagination.PaginatedResponse[Transaction], error) {
	return s.accountRepository.FindTransactions(accountID, statuses, itemsPerPage, pageNumber)
}

func (s accountService) Update(accountID string, account *Account) error {
	return s.accountRepository.Update(accountID, account)
}
