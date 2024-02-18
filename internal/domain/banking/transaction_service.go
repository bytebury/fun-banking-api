package banking

import (
	"errors"
)

type TransactionService interface {
	FindByID(id string) (Transaction, error)
	Approve(id string, transaction *Transaction) error
	Decline(id string, transaction *Transaction) error
	Create(transaction *Transaction) error
}

type transactionService struct {
	transactionRepository TransactionRepository
}

func NewTransactionService(transactionRepository TransactionRepository) TransactionService {
	return transactionService{transactionRepository}
}

func (s transactionService) FindByID(id string) (Transaction, error) {
	var transaction Transaction
	err := s.transactionRepository.FindByID(id, &transaction)
	return transaction, err
}

func (s transactionService) Approve(id string, transaction *Transaction) error {
	if transaction.Status != TransactionPending {
		return errors.New("transaction has already been processed")
	}

	transaction.Status = TransactionApproved

	return s.transactionRepository.Update(id, transaction)
}

func (s transactionService) Decline(id string, transaction *Transaction) error {
	if transaction.Status != TransactionPending {
		return errors.New("transaction has already been processed")
	}

	transaction.Status = TransactionDeclined

	return s.transactionRepository.Update(id, transaction)
}

func (s transactionService) Create(transaction *Transaction) error {
	return s.transactionRepository.Create(transaction)
}
