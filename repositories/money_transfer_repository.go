package repositories

import (
	"golfer/database"
	"golfer/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MoneyTransferRepository struct {
	db *gorm.DB
}

func NewMoneyTransferRepository() *MoneyTransferRepository {
	return &MoneyTransferRepository{
		db: database.DB,
	}
}

func (repository MoneyTransferRepository) Create(moneyTransfer *models.MoneyTransfer) error {
	return repository.db.Create(&moneyTransfer).Error
}

func (repository MoneyTransferRepository) FindByID(moneyTransferID string, moneyTransfer *models.MoneyTransfer) error {
	return repository.db.Preload("Account").Preload("Account.Customer").Preload("Account.Customer.Bank").First(&moneyTransfer, moneyTransferID).Error
}

func (repository MoneyTransferRepository) FindByAccount(accountID string, moneyTransfers *[]models.MoneyTransfer, count *int64, c *gin.Context) error {
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

	if err := query.Model(&moneyTransfers).Count(count).Error; err != nil {
		return err
	}

	return query.Order("updated_at DESC").Limit(limit).Offset(offset).Find(&moneyTransfers).Error
}

func (repository MoneyTransferRepository) FindByUserID(userID string, transfers *[]models.MoneyTransfer) error {
	return repository.db.Model(&models.MoneyTransfer{}).
		Joins("JOIN accounts ON accounts.id = money_transfers.account_id").
		Joins("JOIN customers ON customers.id = accounts.customer_id").
		Joins("JOIN banks ON banks.id = customers.bank_id").
		Joins("JOIN users ON users.id = banks.user_id").
		Where("money_transfers.status = ?", "pending").
		Where("users.id = ?", userID).
		Find(&transfers).Error
}

func (repository MoneyTransferRepository) Update(moneyTransfer *models.MoneyTransfer) error {
	return repository.db.Save(&moneyTransfer).Error
}

func (repository MoneyTransferRepository) Delete(moneyTransferID string) error {
	return repository.db.Delete(&models.MoneyTransfer{}, "id = ?", moneyTransferID).Error
}
