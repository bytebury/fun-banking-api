package repositories

import (
	"golfer/database"
	"golfer/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{
		db: database.DB,
	}
}

func (tr TransactionRepository) Create(transfer *models.Transaction) error {
	return tr.db.Create(&transfer).Error
}

func (tr TransactionRepository) FindByID(transferID string, transfer *models.Transaction) error {
	return tr.db.Preload("Account").Preload("Account.Customer").Preload("Account.Customer.Bank").First(&transfer, transferID).Error
}

func (tr TransactionRepository) FindByAccount(accountID string, transfers *[]models.Transaction, count *int64, c *gin.Context) error {
	statuses := c.QueryArray("status")
	limit, limitErr := strconv.Atoi(c.Query("limit"))
	page, pageErr := strconv.Atoi(c.Query("page"))

	if limitErr != nil {
		limit = 5
	}

	if pageErr != nil {
		page = 1
	}

	offset := (page - 1) * limit

	query := tr.db.Preload("User")
	query = query.Where("account_id = ?", accountID)

	if len(statuses) > 0 {
		query = query.Where("status IN ?", statuses)
	}

	if err := query.Model(&transfers).Count(count).Error; err != nil {
		return err
	}

	return query.Order("updated_at DESC").Limit(limit).Offset(offset).Find(&transfers).Error
}

func (tr TransactionRepository) FindByUserID(userID string, transfers *[]models.Transaction) error {
	unionQuery := "(SELECT bank_id FROM employees WHERE user_id = ? UNION SELECT id FROM banks WHERE user_id = ?)"

	return tr.db.Model(&models.Transaction{}).
		Joins("JOIN accounts ON transfers.account_id = accounts.id").
		Joins("JOIN customers ON accounts.customer_id = customers.id").
		Joins("JOIN banks ON customers.bank_id = banks.id").
		Where("banks.id IN (?)", gorm.Expr(unionQuery, userID, userID)).
		Where("transfers.status = ?", "pending").
		Preload("Account.Customer").
		Find(&transfers).Error
}

func (tr TransactionRepository) Update(transfer *models.Transaction) error {
	return tr.db.Save(&transfer).Error
}

func (tr TransactionRepository) Delete(transferID string) error {
	return tr.db.Delete(&models.Transaction{}, "id = ?", transferID).Error
}
