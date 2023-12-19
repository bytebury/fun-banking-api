package repositories

import (
	"golfer/database"
	"golfer/models"

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
