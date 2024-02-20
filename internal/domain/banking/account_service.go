package banking

import (
	"funbanking/package/utils"
)

type AccountService interface {
	FindByID(accountID string) (Account, error)
	FindTransactions(accountID string) ([]Transaction, error)
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
func (s accountService) FindTransactions(accountID string) ([]Transaction, error) {
	var transactions []Transaction
	err := s.accountRepository.FindTransactions(accountID, &transactions)
	return utils.Listify(transactions), err
}

func (s accountService) Update(accountID string, account *Account) error {
	return s.accountRepository.Update(accountID, account)
}
