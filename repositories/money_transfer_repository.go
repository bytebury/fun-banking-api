package repositories

import (
	"golfer/database"
	"golfer/models"

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
	return repository.db.Preload("Account").First(&moneyTransfer, moneyTransferID).Error
}

func (repository MoneyTransferRepository) FindByAccount(accountID string, moneyTransfers *[]models.MoneyTransfer) error {
	return repository.db.Preload("User").Find(&moneyTransfers, "account_id = ?", accountID).Error
}

func (repository MoneyTransferRepository) Update(moneyTransfer *models.MoneyTransfer) error {
	return repository.db.Save(&moneyTransfer).Error
}

func (repository MoneyTransferRepository) Delete(moneyTransferID string) error {
	return repository.db.Delete(&models.MoneyTransfer{}, "id = ?", moneyTransferID).Error
}
