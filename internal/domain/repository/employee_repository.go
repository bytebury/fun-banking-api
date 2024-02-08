package repository

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/infrastructure/persistence"

	"gorm.io/gorm"
)

type EmployeeRepository interface {
	FindAllByBankID(bankID string, employees *[]model.Employee) error
	FindAllByUserID(userID string, employees *[]model.Employee) error
	Create(employee *model.Employee) error
	Delete(id string) error
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository() EmployeeRepository {
	return employeeRepository{db: persistence.DB}
}

func (r employeeRepository) FindAllByBankID(bankID string, employees *[]model.Employee) error {
	return r.db.Preload("User").Preload("Bank").Find(&employees, "bank_id = ?", bankID).Error
}

func (r employeeRepository) FindAllByUserID(userID string, employees *[]model.Employee) error {
	return r.db.Preload("User").Preload("Bank.User").Find(&employees, "user_id = ?", userID).Error
}

func (r employeeRepository) Create(employee *model.Employee) error {
	return r.db.Create(&employee).Error
}

func (r employeeRepository) Delete(id string) error {
	return r.db.Delete(&model.Employee{}, "id = ?", id).Error
}
