package repositories

import (
	"golfer/database"
	"golfer/models"

	"gorm.io/gorm"
)

type EmployeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository() *EmployeeRepository {
	return &EmployeeRepository{
		db: database.DB,
	}
}

func (repository EmployeeRepository) Create(employee *models.Employee) error {
	return repository.db.Create(&employee).Error
}

func (repository EmployeeRepository) FindByBank(bankID string, employees *[]models.Employee) error {
	return repository.db.Preload("User").Preload("Bank").Find(&employees, "bank_id = ?", bankID).Error
}

func (repository EmployeeRepository) FindByUser(userID string, employees *[]models.Employee) error {
	return repository.db.Preload("User").Preload("Bank.User").Find(&employees, "user_id = ?", userID).Error
}

func (repository EmployeeRepository) Delete(employeeID string) error {
	return repository.db.Delete(&models.Employee{}, "id = ?", employeeID).Error
}
