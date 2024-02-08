package service

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"funbanking/package/utils"
)

type AccountService interface {
	FindByID(id string) (model.Account, error)
	FindTransactions(id string) ([]model.Transaction, error)
	Update(id string, account *model.Account) error
}

type accountService struct {
	accountRepository repository.AccountRepository
}

func NewAccountService(accountRepository repository.AccountRepository) AccountService {
	return accountService{accountRepository}
}

func (s accountService) FindByID(id string) (model.Account, error) {
	var account model.Account
	err := s.accountRepository.FindByID(id, &account)
	return account, err
}

// TODO: THIS IS GOING TO BE PAGINATED
func (s accountService) FindTransactions(id string) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := s.accountRepository.FindTransactions(id, &transactions)
	return utils.Listify(transactions), err
}

func (s accountService) Update(id string, account *model.Account) error {
	return s.accountRepository.Update(id, account)
}
