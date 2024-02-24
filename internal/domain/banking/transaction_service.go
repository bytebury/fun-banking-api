package banking

import (
	"errors"
	"funbanking/internal/infrastructure/persistence"
	"funbanking/package/utils"
	"strconv"

	"gorm.io/gorm"
)

type TransactionService interface {
	FindByID(transactionID string) (Transaction, error)
	FindAllPendingTransactions(userID string) ([]Transaction, error)
	Approve(userID string, transactionID string) (Transaction, error)
	Decline(userID string, transactionID string) (Transaction, error)
	Create(userID string, transaction *Transaction) error
}

type transactionService struct {
	transactionRepository TransactionRepository
	bankService           BankService
	accountService        AccountService
}

func NewTransactionService(transactionRepository TransactionRepository) TransactionService {
	return transactionService{
		transactionRepository: transactionRepository,
		bankService: NewBankService(
			NewBankRepository(),
		),
		accountService: NewAccountService(
			NewAccountRepository(),
		),
	}
}

func (s transactionService) FindByID(transactionID string) (Transaction, error) {
	var transaction Transaction
	err := s.transactionRepository.FindByID(transactionID, &transaction)
	return transaction, err
}

func (s transactionService) FindAllPendingTransactions(userID string) ([]Transaction, error) {
	var transactions []Transaction
	err := s.transactionRepository.FindAllPendingTransactions(userID, &transactions)
	return utils.Listify(transactions), err
}

func (s transactionService) Approve(userID string, transactionID string) (Transaction, error) {
	var transaction Transaction

	err := persistence.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.transactionRepository.FindByID(transactionID, &transaction); err != nil {
			return err
		}

		if transaction.Status != TransactionPending {
			return errors.New("transaction has already been processed")
		}

		// Reflect the account's balance
		transaction.Status = TransactionApproved
		userIDPtr, _ := utils.StringToUIntPointer(userID)
		transaction.UserID = userIDPtr
		transaction.CurrentBalance = transaction.Account.Balance + float64(transaction.Amount)

		if err := s.transactionRepository.Update(transactionID, &transaction); err != nil {
			return err
		}

		if _, err := s.accountService.AddToBalance(strconv.Itoa(int(transaction.AccountID)), transaction.Amount); err != nil {
			return err
		}

		return nil
	})

	return transaction, err
}

func (s transactionService) Decline(userID string, transactionID string) (Transaction, error) {
	var transaction Transaction

	if err := s.transactionRepository.FindByID(transactionID, &transaction); err != nil {
		return transaction, err
	}

	if transaction.Status != TransactionPending {
		return transaction, errors.New("this transaction has already been processed")
	}

	transaction.Status = TransactionDeclined
	userIDPtr, _ := utils.StringToUIntPointer(userID)
	transaction.UserID = userIDPtr

	err := s.transactionRepository.Update(transactionID, &transaction)

	return transaction, err
}

func (s transactionService) Create(userID string, transaction *Transaction) error {
	err := s.transactionRepository.Create(transaction)

	if err != nil {
		return err
	}

	account, err := s.accountService.FindByID(strconv.Itoa(int(transaction.AccountID)))

	if err != nil {
		return err
	}

	if s.bankService.IsEmployee(strconv.Itoa(int(account.Customer.BankID)), userID) {
		t, err := s.Approve(userID, strconv.Itoa(int(transaction.ID)))
		*transaction = t
		return err
	}

	return err
}
