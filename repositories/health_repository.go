package repositories

import (
	"golfer/database"
	"golfer/models"

	"gorm.io/gorm"
)

type HealthRepository struct {
	db *gorm.DB
}

func NewHealthRepository() *HealthRepository {
	return &HealthRepository{
		db: database.DB,
	}
}

func (repository HealthRepository) GetHealthCheck(health *models.Health) error {
	health.Name = "Fun Banking"
	health.Version = "1.0.0"
	health.Message = "Everything is up and running!"

	repository.db.Model(&models.User{}).Count(&health.NumberOfUsers)
	repository.db.Model(&models.Bank{}).Count(&health.NumberOfBanks)
	repository.db.Model(&models.Customer{}).Count(&health.NumberOfCustomers)
	repository.db.Model(&models.MoneyTransfer{}).Count(&health.NumberOfTransfers)

	return nil
}
