package repository

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/infrastructure/persistence"

	"gorm.io/gorm"
)

type HealthRepository interface {
	GetHealthCheck(health *model.Health) error
}

type healthRepository struct {
	db *gorm.DB
}

func NewHealthRepository() HealthRepository {
	return healthRepository{db: persistence.DB}
}

func (r healthRepository) GetHealthCheck(health *model.Health) error {
	health.Name = "Fun Banking"
	health.Version = "1.0.0"
	health.Message = "Everything is up and running!"

	// TODO(marcello): add these
	// r.db.Model(&model.User{}).Count(&health.NumberOfUsers)
	// r.db.Model(&model.Bank{}).Count(&health.NumberOfBanks)
	// r.db.Model(&model.Customer{}).Count(&health.NumberOfCustomers)
	// r.db.Model(&model.Transaction{}).Count(&health.NumberOfTransactions)

	return nil
}
