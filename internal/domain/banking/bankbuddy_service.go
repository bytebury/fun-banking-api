package banking

import (
	"errors"
	"fmt"
	"funbanking/internal/infrastructure/persistence"
	"math"
	"strconv"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

type BankBuddyService interface {
	Transfer(transfer *BankBuddyTransfer) error
	FindRecipients(bankID string) ([]Customer, error)
}

type bankBuddyService struct {
	transactionService TransactionService
	accountService     AccountService
}

func NewBankBuddyService(transactionService TransactionService, accountService AccountService) BankBuddyService {
	return bankBuddyService{
		transactionService,
		accountService,
	}
}

func (s bankBuddyService) Transfer(transfer *BankBuddyTransfer) error {
	transfer.Amount = math.Abs(transfer.Amount)

	return persistence.DB.Transaction(func(tx *gorm.DB) error {
		fromAccount, noAccFound1 := s.accountService.FindByID(strconv.Itoa(int(transfer.FromAccountID)))
		toAccount, noAccFound2 := s.accountService.FindByID(strconv.Itoa(int(transfer.ToAccountID)))

		if noAccFound1 != nil {
			return noAccFound1
		}

		if noAccFound2 != nil {
			return noAccFound2
		}

		if fromAccount.Balance < transfer.Amount {
			return errors.New("insufficient funds")
		}

		transaction := Transaction{
			AccountID: transfer.FromAccountID,
			Amount:    transfer.Amount * -1,
			Description: fmt.Sprintf(
				"%s. Sent to %s %s",
				transfer.Description,
				cases.Title(language.AmericanEnglish).String(toAccount.Customer.FirstName),
				cases.Title(language.AmericanEnglish).String(toAccount.Customer.LastName),
			),
			BankBuddySenderID: &fromAccount.CustomerID,
			Type:              "bank_buddy",
		}

		if err := s.transactionService.Create("", &transaction); err != nil {
			return err
		}

		if _, err := s.transactionService.Approve("", strconv.Itoa(int(transaction.ID))); err != nil {
			return err
		}

		transaction = Transaction{
			AccountID:         transfer.ToAccountID,
			Amount:            transfer.Amount,
			Description:       transfer.Description,
			Type:              "bank_buddy",
			BankBuddySenderID: &fromAccount.CustomerID,
		}

		if err := s.transactionService.Create("", &transaction); err != nil {
			return err
		}

		if _, err := s.transactionService.Approve("", strconv.Itoa(int(transaction.ID))); err != nil {
			return err
		}

		return nil
	})
}

func (s bankBuddyService) FindRecipients(bankID string) ([]Customer, error) {
	var customers []Customer

	if err := persistence.DB.Preload("Accounts").Find(&customers, "bank_id = ?", bankID).Error; err != nil {
		return make([]Customer, 0), err
	}

	return censorPrivateFields(customers), nil
}

func censorPrivateFields(customers []Customer) []Customer {
	for i := range customers {
		customers[i].PIN = ""
		for j := range customers[i].Accounts {
			customers[i].Accounts[j].Balance = 0
		}
	}
	return customers
}
