package repositories

import (
	"golfer/database"
	"golfer/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransferRepository struct {
	db *gorm.DB
}

func NewTransferRepository() *TransferRepository {
	return &TransferRepository{
		db: database.DB,
	}
}

func (repository TransferRepository) Create(transfer *models.Transfer) error {
	return repository.db.Create(&transfer).Error
}

func (repository TransferRepository) FindByID(transferID string, transfer *models.Transfer) error {
	return repository.db.Preload("Account").Preload("Account.Customer").Preload("Account.Customer.Bank").First(&transfer, transferID).Error
}

func (repository TransferRepository) FindByAccount(accountID string, transfers *[]models.Transfer, count *int64, c *gin.Context) error {
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

	query := repository.db.Preload("User")
	query = query.Where("account_id = ?", accountID)

	if len(statuses) > 0 {
		query = query.Where("status IN ?", statuses)
	}

	if err := query.Model(&transfers).Count(count).Error; err != nil {
		return err
	}

	return query.Order("updated_at DESC").Limit(limit).Offset(offset).Find(&transfers).Error
}

func (repository TransferRepository) FindByUserID(userID string, transfers *[]models.Transfer) error {
	unionQuery := "(SELECT bank_id FROM employees WHERE user_id = ? UNION SELECT id FROM banks WHERE user_id = ?)"

	// Main query
	return repository.db.Model(&models.Transfer{}).
		Joins("JOIN accounts ON transfers.account_id = accounts.id").
		Joins("JOIN customers ON accounts.customer_id = customers.id").
		Joins("JOIN banks ON customers.bank_id = banks.id").
		Where("banks.id IN (?)", gorm.Expr(unionQuery, userID, userID)).
		Where("transfers.status = ?", "pending").
		Preload("Account.Customer").
		Find(&transfers).Error
}

func (repository TransferRepository) Update(transfer *models.Transfer) error {
	return repository.db.Save(&transfer).Error
}

func (repository TransferRepository) Delete(transferID string) error {
	return repository.db.Delete(&models.Transfer{}, "id = ?", transferID).Error
}
