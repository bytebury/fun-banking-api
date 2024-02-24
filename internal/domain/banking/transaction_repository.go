package banking

import (
	"funbanking/internal/infrastructure/persistence"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindByID(transactionID string, transaction *Transaction) error
	Create(transaction *Transaction) error
	Update(transactionID string, transaction *Transaction) error
	FindAllPendingTransactions(userID string, transactions *[]Transaction) error
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

func (r transactionRepository) FindAllPendingTransactions(userID string, transactions *[]Transaction) error {
	unionQuery := "(SELECT bank_id FROM employees WHERE user_id = ? UNION SELECT id FROM banks WHERE user_id = ?)"

	return r.db.Model(&Transaction{}).
		Joins("JOIN accounts ON transactions.account_id = accounts.id").
		Joins("JOIN customers ON accounts.customer_id = customers.id").
		Joins("JOIN banks ON customers.bank_id = banks.id").
		Where("banks.id IN (?)", gorm.Expr(unionQuery, userID, userID)).
		Where("transactions.status = ?", "pending").
		Preload("Account.Customer").
		Find(&transactions).Error
}

func (r transactionRepository) Create(transaction *Transaction) error {
	return r.db.Create(&transaction).Error
}

func (r transactionRepository) Update(transactionID string, transaction *Transaction) error {
	var foundTransaction Transaction
	if err := r.db.First(&foundTransaction, "id = ?", transactionID).Error; err != nil {
		return err
	}

	if transaction.CurrentBalance == 0 {
		transaction.CurrentBalance = foundTransaction.CurrentBalance
	}

	if transaction.UserID == nil {
		transaction.UserID = foundTransaction.UserID
	}

	if transaction.Status == TransactionPending {
		transaction.Status = foundTransaction.Status
	}

	return r.db.Model(&foundTransaction).Select("Status", "UserID", "CurrentBalance").Updates(&transaction).Error
}
