package banking

import (
	"funbanking/internal/infrastructure/persistence"

	"gorm.io/gorm"
)

type EmployeeRepository interface {
	FindAllByBankID(bankID string, employees *[]Employee) error
	FindAllByUserID(userID string, employees *[]Employee) error
	Create(employee *Employee) error
	Delete(employeeID string) error
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository() EmployeeRepository {
	return employeeRepository{db: persistence.DB}
}

func (r employeeRepository) FindAllByBankID(bankID string, employees *[]Employee) error {
	return r.db.Preload("User").Preload("Bank").Find(&employees, "bank_id = ?", bankID).Error
}

func (r employeeRepository) FindAllByUserID(userID string, employees *[]Employee) error {
	return r.db.Preload("User").Preload("Bank.User").Find(&employees, "user_id = ?", userID).Error
}

func (r employeeRepository) Create(employee *Employee) error {
	return r.db.Create(&employee).Error
}

func (r employeeRepository) Delete(employeeID string) error {
	return r.db.Delete(&Employee{}, "id = ?", employeeID).Error
}
