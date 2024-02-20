package banking

import (
	"funbanking/internal/infrastructure/persistence"

	"gorm.io/gorm"
)

type AccountRepository interface {
	FindByID(accountID string, account *Account) error
	FindTransactions(accountID string, transactions *[]Transaction) error
	Update(accountID string, account *Account) error
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository() AccountRepository {
	return accountRepository{db: persistence.DB}
}

func (r accountRepository) FindByID(accountID string, account *Account) error {
	return r.db.Preload("Customer").First(&account, "id = ?", accountID).Error
}

func (r accountRepository) FindTransactions(accountID string, transactions *[]Transaction) error {
	return r.db.Find(&transactions, "account_id = ?", accountID).Error
}

func (r accountRepository) Update(accountID string, account *Account) error {
	var foundAccount Account

	if err := r.db.First(&account, "id = ?", accountID).Error; err != nil {
		return err
	}

	if account.Name == "" {
		account.Name = foundAccount.Name
	}

	return r.db.Model(&foundAccount).Select("Name").Updates(account).Error
}
