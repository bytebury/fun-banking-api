package repositories

import (
	"golfer/database"
	"golfer/models"

	"gorm.io/gorm"
)

type BankRepository struct {
	db *gorm.DB
}

func NewBankRepository() *BankRepository {
	return &BankRepository{
		db: database.DB,
	}
}

func (repository BankRepository) Create(bank *models.Bank) error {
	return repository.db.Create(&bank).Error
}

func (repository BankRepository) FindBanksByUserID(userID string, banks *[]models.Bank) error {
	return repository.db.Preload("User").Find(&banks, "user_id = ?", userID).Error
}

func (repository BankRepository) FindByID(bankID string, bank *models.Bank) error {
	return repository.db.Preload("User").First(&bank, bankID).Error
}

func (repository BankRepository) FindCustomers(bankID string, customers *[]models.Customer) error {
	return repository.db.Find(&customers, "bank_id = ?", bankID).Error
}

func (repository BankRepository) Update(bank *models.Bank) error {
	return repository.db.Save(&bank).Error
}

func (repository BankRepository) Delete(bankID string) error {
	return repository.db.Delete(&models.Bank{}, "id = ?", bankID).Error
}
