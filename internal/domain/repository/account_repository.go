package repository

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/infrastructure/persistence"

	"gorm.io/gorm"
)

type AccountRepository interface {
	FindByID(accountID string, account *model.Account) error
	FindTransactions(accountID string, transactions *[]model.Transaction) error
	Update(account *model.Account) error
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository() AccountRepository {
	return accountRepository{db: persistence.DB}
}

func (r accountRepository) FindByID(accountID string, account *model.Account) error {
	return r.db.First(&account, "id = ?", accountID).Error
}

func (r accountRepository) FindTransactions(accountID string, transactions *[]model.Transaction) error {
	return r.db.Find(&transactions, "account_id = ?", accountID).Error
}

func (r accountRepository) Update(account *model.Account) error {
	var foundAccount model.Account

	if err := r.db.First(&account).Error; err != nil {
		return err
	}

	if account.Name == "" {
		account.Name = foundAccount.Name
	}

	return r.db.Model(&account).Select("Name").Updates(account).Error
}
