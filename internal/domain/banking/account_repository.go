package banking

import (
	"funbanking/internal/infrastructure/pagination"
	"funbanking/internal/infrastructure/persistence"

	"gorm.io/gorm"
)

type AccountRepository interface {
	FindByID(accountID string, account *Account) error
	FindTransactions(accountID string, statuses []string, itemsPerPage int, pageNumber int) (pagination.PaginatedResponse[Transaction], error)
	Update(accountID string, account *Account) error
	AddToBalance(accountID string, amount float64) (Account, error)
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

func (r accountRepository) FindTransactions(accountID string, statuses []string, itemsPerPage int, pageNumber int) (pagination.PaginatedResponse[Transaction], error) {
	query := r.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "first_name", "last_name", "username")
	})
	query = query.Where("account_id = ?", accountID)

	if len(statuses) > 0 {
		query = query.Where("status IN ?", statuses)
	}

	query = query.Order("updated_at DESC")

	return pagination.Find[Transaction](query, pageNumber, itemsPerPage)
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

func (r accountRepository) AddToBalance(accountID string, amount float64) (Account, error) {
	var foundAccount Account

	if err := r.db.First(&foundAccount, "id = ?", accountID).Error; err != nil {
		return foundAccount, err
	}

	foundAccount.Balance += amount

	err := r.db.Model(&foundAccount).Select("Balance").Updates(&foundAccount).Error
	return foundAccount, err
}
