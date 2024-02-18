package banking

import (
	"funbanking/package/utils"
)

type AccountService interface {
	FindByID(id string) (Account, error)
	FindTransactions(id string) ([]Transaction, error)
	Update(id string, account *Account) error
}

type accountService struct {
	accountRepository AccountRepository
}

func NewAccountService(accountRepository AccountRepository) AccountService {
	return accountService{accountRepository}
}

func (s accountService) FindByID(id string) (Account, error) {
	var account Account
	err := s.accountRepository.FindByID(id, &account)
	return account, err
}

// TODO: THIS IS GOING TO BE PAGINATED
func (s accountService) FindTransactions(id string) ([]Transaction, error) {
	var transactions []Transaction
	err := s.accountRepository.FindTransactions(id, &transactions)
	return utils.Listify(transactions), err
}

func (s accountService) Update(id string, account *Account) error {
	return s.accountRepository.Update(id, account)
}
