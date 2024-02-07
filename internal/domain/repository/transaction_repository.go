package repository

import (
	"funbanking/internal/infrastructure/persistence"

	"gorm.io/gorm"
)

type TransactionRepository interface {

}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository() TransactionRepository {
	return transactionRepository{ db: persistence.DB }
}
