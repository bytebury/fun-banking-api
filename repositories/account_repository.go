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

func (repository AccountRepository) GetMonthlyData(accountID string) ([]models.AccountMonthlySummary, error) {
	// Get the current month and year
	currentMonth := time.Now().Month()
	currentYear := time.Now().Year()

	// Calculate the start of the current month
	startOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.UTC)

	var monthlyAggregations []models.AccountMonthlySummary
	err := repository.db.Model(&models.Transaction{}).
		Select("TO_CHAR(updated_at, 'Month') as month, SUM(CASE WHEN amount >= 0 THEN amount ELSE 0 END) as deposits, SUM(CASE WHEN amount < 0 THEN amount ELSE 0 END) as withdrawals").
		Where("updated_at >= ? AND account_id = ?", startOfMonth, accountID).
		Group("month").
		Order("month").
		Find(&monthlyAggregations).Error

	return monthlyAggregations, err

}
