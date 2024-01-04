package repositories

import (
	"golfer/database"
	"golfer/models"
	"time"

	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{
		db: database.DB,
	}
}

func (repository AccountRepository) Create(account *models.Account) error {
	return repository.db.Create(&account).Error
}

func (repository AccountRepository) FindByID(accountID string, account *models.Account) error {
	return repository.db.Preload("Customer").Preload("Customer.Bank").First(&account, accountID).Error
}

func (repository AccountRepository) FindByCustomer(customerID string, accounts *[]models.Account) error {
	return repository.db.Find(&accounts, "customer_id = ?", customerID).Error
}

func (repository AccountRepository) Update(account *models.Account) error {
	return repository.db.Save(&account).Error
}

func (repository AccountRepository) Delete(accountID string) error {
	return repository.db.Delete(&models.Account{}, "id = ?", accountID).Error
}

func (repository AccountRepository) GetTransactionHistoricalData(accountID string, daysAgo int) ([]models.DailyTransferSummary, error) {
	var dailySummaries []models.DailyTransferSummary
	xDaysAgo := time.Now().AddDate(0, 0, -daysAgo)

	subQuery := repository.db.Model(&models.Transaction{}).
		Select("MAX(updated_at) as max_updated_at").
		Where("updated_at >= ? AND account_id = ?", xDaysAgo, accountID).
		Group("DATE(updated_at)")

	result := repository.db.Model(&models.Transaction{}).
		Select("DATE(updated_at) as date, current_balance as total_balance").
		Joins("JOIN (?) as sub on sub.max_updated_at = transactions.updated_at", subQuery).
		Where("transactions.account_id = ? AND status = ?", accountID, "approved").
		Order("date").
		Scan(&dailySummaries)

	return dailySummaries, result.Error
}
