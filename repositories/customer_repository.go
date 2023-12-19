package repositories

import (
	"golfer/database"
	"golfer/models"

	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository() *CustomerRepository {
	return &CustomerRepository{
		db: database.DB,
	}
}

func (repository CustomerRepository) Create(customer *models.Customer) error {
	return repository.db.Create(&customer).Error
}

func (repository CustomerRepository) FindByID(customerID string, customer *models.Customer) error {
	return repository.db.Preload("Bank").Preload("Bank.User").First(&customer, customerID).Error
}

func (repository CustomerRepository) Update(customer *models.Customer) error {
	return repository.db.Save(&customer).Error
}

func (repository CustomerRepository) Delete(customerID string) error {
	return repository.db.Delete(&models.Customer{}, "id = ?", customerID).Error
}
