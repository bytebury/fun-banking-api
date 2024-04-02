package banking

import (
	"funbanking/internal/infrastructure/persistence"
	"math"
	"strconv"

	"gorm.io/gorm"
)

type BankBuddyService interface {
	Transfer(transfer *BankBuddyTransfer) error
	FindRecipients(bankID string) ([]Customer, error)
}

type bankBuddyService struct {
	transactionService TransactionService
}

func NewBankBuddyService(transactionService TransactionService) BankBuddyService {
	return bankBuddyService{
		transactionService,
	}
}

func (s bankBuddyService) Transfer(transfer *BankBuddyTransfer) error {
	transfer.Amount = math.Abs(transfer.Amount)

	return persistence.DB.Transaction(func(tx *gorm.DB) error {
		transaction := Transaction{
			AccountID:   transfer.FromAccountID,
			Amount:      transfer.Amount * -1,
			Description: transfer.Description,
			Type:        "bankbuddy",
		}

		if err := s.transactionService.Create("", &transaction); err != nil {
			return err
		}

		if _, err := s.transactionService.Approve("", strconv.Itoa(int(transaction.ID))); err != nil {
			return err
		}

		transaction = Transaction{
			AccountID:   transfer.ToAccountID,
			Amount:      transfer.Amount,
			Description: transfer.Description,
			Type:        "bankbuddy",
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
