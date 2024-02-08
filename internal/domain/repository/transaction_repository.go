package repository

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/infrastructure/persistence"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindByID(transactionID string, transaction *model.Transaction) error
	Create(transaction *model.Transaction) error
	Update(transaction *model.Transaction) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository() TransactionRepository {
	return transactionRepository{db: persistence.DB}
}

func (r transactionRepository) FindByID(transactionID string, transaction *model.Transaction) error {
	return r.db.Preload("Account.Customer.Bank").First(&transaction, "id = ?", transactionID).Error
}

func (r transactionRepository) Create(transaction *model.Transaction) error {
	return r.db.Create(&transaction).Error
}

func (r transactionRepository) Update(transaction *model.Transaction) error {
	return r.db.Model(&transaction).Select("Status", "User").Updates(&transaction).Error
}
