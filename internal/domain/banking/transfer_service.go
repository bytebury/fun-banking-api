package banking

import (
	"errors"
	"fmt"
	"funbanking/internal/infrastructure/persistence"
	"math"
	"strconv"

	"gorm.io/gorm"
)

type TransferService interface {
	Transfer(customerID string, transfer TransferRequest) error
}

type transferService struct {
	accountService     AccountService
	transactionService TransactionService
}

func NewTransferService(accountService AccountService, transactionService TransactionService) TransferService {
	return transferService{accountService, transactionService}
}

func (s transferService) Transfer(customerID string, transfer TransferRequest) error {
	var toAccount Account
	var fromAccount Account

	if account, err := s.accountService.FindByID(strconv.Itoa(int(transfer.FromAccountID))); err != nil {
		return errors.New("account does not exist")
	} else {
		fromAccount = account
	}

	if account, err := s.accountService.FindByID(strconv.Itoa(int(transfer.ToAccountID))); err != nil {
		return errors.New("account does not exist")
	} else {
		toAccount = account
	}

	if toAccount.CustomerID != fromAccount.CustomerID || strconv.Itoa(int(toAccount.CustomerID)) != customerID {
		return errors.New("you can only transfer between accounts you own")
	}

	return s.transfer(fromAccount, toAccount, transfer)
}

func (s transferService) transfer(fromAccount, toAccount Account, transfer TransferRequest) error {
	transfer.Amount = math.Abs(transfer.Amount)

	return persistence.DB.Transaction(func(tx *gorm.DB) error {
		if fromAccount.Balance < transfer.Amount {
			return errors.New("insufficient funds")
		}

		description := fmt.Sprintf("Funds transfer from %s to %s", fromAccount.Name, toAccount.Name)

		withdrawTransaction := Transaction{
			AccountID:   fromAccount.ID,
			Amount:      transfer.Amount * -1,
			Description: description,
			Type:        "transfer",
		}

		if err := s.transactionService.Create("", &withdrawTransaction); err != nil {
			return err
		}

		if _, err := s.transactionService.Approve("", strconv.Itoa(int(withdrawTransaction.ID))); err != nil {
			return err
		}

		depositTransaction := Transaction{
			AccountID:   toAccount.ID,
			Amount:      transfer.Amount,
			Description: description,
			Type:        "transfer",
		}

		if err := s.transactionService.Create("", &depositTransaction); err != nil {
			return err
		}

		if _, err := s.transactionService.Approve("", strconv.Itoa(int(depositTransaction.ID))); err != nil {
			return err
		}

		return nil
	})
}
