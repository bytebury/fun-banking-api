package repository

import (
	"funbanking/internal/infrastructure/persistence"

	"gorm.io/gorm"
)

type EmployeeRepository interface {

}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository() EmployeeRepository {
	return employeeRepository{ db: persistence.DB }
}
