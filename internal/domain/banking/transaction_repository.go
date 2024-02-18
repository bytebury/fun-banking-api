package banking

import (
	"funbanking/internal/infrastructure/persistence"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindByID(transactionID string, transaction *Transaction) error
	Create(transaction *Transaction) error
	Update(transactionID string, transaction *Transaction) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository() TransactionRepository {
	return transactionRepository{db: persistence.DB}
}

func (r transactionRepository) FindByID(transactionID string, transaction *Transaction) error {
	return r.db.Preload("Account.Customer.Bank").First(&transaction, "id = ?", transactionID).Error
}

func (r transactionRepository) Create(transaction *Transaction) error {
	return r.db.Create(&transaction).Error
}

func (r transactionRepository) Update(transactionID string, transaction *Transaction) error {
	var foundTransaction Transaction
	if err := r.db.First(&foundTransaction, "id = ?", transactionID).Error; err != nil {
		return err
	}
	return r.db.Model(&foundTransaction).Select("Status", "User").Updates(&transaction).Error
}
