package service

import (
	"errors"
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
)

type TransactionService interface {
	FindByID(id string) (model.Transaction, error)
	Approve(id string, transaction *model.Transaction) error
	Decline(id string, transaction *model.Transaction) error
	Create(transaction *model.Transaction) error
}

type transactionService struct {
	transactionRepository repository.TransactionRepository
}

func NewTransactionService(transactionRepository repository.TransactionRepository) TransactionService {
	return transactionService{transactionRepository}
}

func (s transactionService) FindByID(id string) (model.Transaction, error) {
	var transaction model.Transaction
	err := s.transactionRepository.FindByID(id, &transaction)
	return transaction, err
}

func (s transactionService) Approve(id string, transaction *model.Transaction) error {
	if transaction.Status != model.TransactionPending {
		return errors.New("transaction has already been processed")
	}

	transaction.Status = model.TransactionApproved

	return s.transactionRepository.Update(id, transaction)
}

func (s transactionService) Decline(id string, transaction *model.Transaction) error {
	if transaction.Status != model.TransactionPending {
		return errors.New("transaction has already been processed")
	}

	transaction.Status = model.TransactionDeclined

	return s.transactionRepository.Update(id, transaction)
}

func (s transactionService) Create(transaction *model.Transaction) error {
	return s.transactionRepository.Create(transaction)
}
